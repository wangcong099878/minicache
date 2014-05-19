package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	_ "strings"
)

var str string
var msg = make([]byte, 1024)


func send(str map[string]string, con net.Conn) (s string, e error) {
	b, err2 := json.Marshal(str)
	if err2 != nil {
		return "", err2
	}

	_, err := con.Write(b)

	if err != nil {
		return "", err
	}

	length, err := con.Read(msg)
	if err != nil {
		return "", err
	}
	result := string(msg[0:length])
	return result, nil
}

func main() {

	var (
		host   = "127.0.0.1"
		port   = "8888"
		remote = host + ":" + port
	)

	con, err := net.Dial("tcp", remote)

	defer con.Close()
	if err != nil {
		fmt.Println("Server not found.")
		os.Exit(-1)
	}

	fmt.Println("Connection OK.")

	//写入数据
	/*req := make(map[string]string)
	req["a"] = "set"
	req["k"] = "haha"
	req["v"] = "sdhkfljdslksd"
	s, err := send(req, con)

	//写入多个
	req["k"] = "haha1"
	req["v"] = "sfdhjkhsjkh"
	s, err = send(req, con)
	fmt.Println(s)
	//写入多个
	req["k"] = "haha2"
	req["v"] = "sfdhjkhsjkh"
	s, err = send(req, con)
	fmt.Println(s)
	//写入多个
	req["k"] = "haha3"
	req["v"] = "sfdhjkhsjkh"
	req["t"] = "30"
	s, err = send(req, con)
	fmt.Println(s)
	//写入多个
	req["k"] = "haha4"
	req["v"] = "sfdhjkhsjkh"
	s, err = send(req, con)
	fmt.Println(s)*/

	//获取数据
	/*req1 := make(map[string]string)
	req1["a"] = "get"
	req1["k"] = "haha"
	s1, err := send(req1, con)

	fmt.Println(s1)*/

	//获取全部
	req2 := make(map[string]string)
	req2["a"] = "all"
	s2, err := send(req2, con)
	fmt.Println(s2)

	//获取总量
	req3 := make(map[string]string)
	req3["a"] = "count"
	s3, err := send(req3, con)

	fmt.Println(s3)

	var d = make(map[string]string)

	err8 := json.Unmarshal([]byte(s2), &d)
	if err8 != nil {
		fmt.Println("有错误")
	}

	fmt.Println(d["data"])

	//删除一条数据
	/*req = make(map[string]string)
	req["a"] = "del"
	req["k"] = "haha"
	s, err = send(req, con)

	fmt.Println(s)*/

}
