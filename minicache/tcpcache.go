package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

/**
 * tcpsocket 键值对缓存
 * 使用方式 cache -p 9999 -ct 10
 */
var p *string = flag.String("p", "8888", "is post")

var ct *string = flag.String("ct", "30", "is cleartime")

var Data = make(map[string]string)

var Expire = make(map[string]int)



/***************************************** func *************************************************/
func set(d map[string]string, con net.Conn) {
	k := d["k"]
	v := d["v"]
	t := d["t"]

	if k != "" && v != "" {
		Data[k] = v
		if t != "" {
			timex := time.Now().Unix()
			timer := int(timex)
			exp, _ := strconv.Atoi(t)
			Expire[k] = timer + exp
		}
		result("1", "", con)
	} else {
		result("2", "", con)
	}
}

func get(d map[string]string, con net.Conn) {
	k := d["k"]
	if k != "" {
		result("1", Data[k], con)
	} else {
		result("2", "", con)
	}
}

func del(d map[string]string, con net.Conn) {
	k := d["k"]
	if k != "" {
		delete(Data, k)
		result("1", "", con)
	} else {
		result("2", "", con)
	}
}

//更新一个key的有效时间
func update(d map[string]string, con net.Conn) {
	k := d["k"]
	t := d["t"]

	if k != "" && t != "" {
		timex := time.Now().Unix()
		timer := int(timex)
		exp, _ := strconv.Atoi(t)
		Expire[k] = timer + exp
		result("1", "", con)
	} else {
		result("2", "", con)
	}
}

func all(d map[string]string, con net.Conn) {
	str := mapToJson(Data)
	result("1", str, con)
}

func count(d map[string]string, con net.Conn) {
	str := strconv.Itoa(len(Data))
	result("1", str, con)
}
/******************************************** tool *************************************************************/
func Deleteexpired() {
	t := time.Now().Unix()
	timer := int(t)
	if len(Expire) > 0 {
		for k, v := range Expire {
			if v < timer {
				delete(Data, k)
				delete(Expire, k)
			}
		}
	}
}

func gc(cleartime string) {
	fmt.Println(cleartime)
	step, _ := strconv.Atoi(cleartime)
	timer1 := time.NewTicker(time.Duration(step) * time.Second)

	for {
		select {
		case <-timer1.C:
			go Deleteexpired()
		}
	}
	/*for {
		time.Sleep(1 * time.Second)
		go Deleteexpired()
	}*/
}

func mapToJson(d map[string]string) string {
	t := ""
	if len(d) > 0 {
		for k, v := range d {
			t += ",\"" + k + "\":\"" + v + "\""
		}
		if t != "" {
			t = strings.Trim(t, ",")
			t = "{" + t + "}"
		}
	}
	return t
}

func result(code string, data string, con net.Conn) {
	var res = make(map[string]string)
	res["code"] = code
	res["data"] = data
	b, err2 := json.Marshal(res)
	if err2 != nil {

	}
	con.Write(b)
}

func main() {
	flag.Parse()
	port := *p
	cleartime := *ct
	go gc(cleartime)
	var (
		host   = "127.0.0.1"
		remote = host + ":" + port
		data   = make([]byte, 1024)
	)
	fmt.Println("Initiating server... (Ctrl-C to stop)", "Listen：", port)
	lis, err := net.Listen("tcp", remote)
	defer lis.Close()
	if err != nil {
		fmt.Println("Error when listen: ", remote)
		os.Exit(-1)
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			fmt.Println("Error accepting client: ", err.Error())
			os.Exit(0)
		}

		go func(con net.Conn) {

			//fmt.Println("New connection: ", con.RemoteAddr())
			for {

				length, err := con.Read(data)
				if err != nil {
					fmt.Printf("Client %v quit.\n", con.RemoteAddr())
					con.Close()
					return
				}

				b := data[0:length]
				var d = make(map[string]string)

				err1 := json.Unmarshal(b, &d)

				if err1 != nil {
					result("5", "", con)
				}
				switch d["a"] {
				case "set":
					set(d, con)
					break
				case "get":
					get(d, con)
					break
				case "del":
					del(d, con)
					break
				case "all":
					all(d, con)
					break
				case "count":
					count(d, con)
					break
				case "update":
					update(d, con)
					break
				default:

				}
			}
		}(conn)
	}
}
