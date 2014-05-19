<?php
	//此处使用thinkphp中的socket客户端包
    include './socket.class.php';
    
    $socket = new Socket(array(
        "host"=>"127.0.0.1",
        "port"=>8888
    ));
    
    $socket->connect();
    
    $req = array();
    $req["a"] = "set";
    $req["k"] = "php1";
    $req["v"] = "sdfdsfds";
    
    $socket->write(json_encode($req));
    
    print_r($socket->read(1024));
    
    $req["a"] = "set";
    $req["k"] = "php2";
    $req["v"] = "sdfdsfdsdffd";
    $req["t"] = "30";
    
    
    $socket->write(json_encode($req));
    
    print_r($socket->read(1024));
    
    $socket->disconnect();