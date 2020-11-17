package main

import (
	"fmt"
	"net"
	"time"
)

/*
	模拟客户端
*/

func main()  {
	fmt.Println("client start")
	time.Sleep(1*time.Second)
	//1、直接连接远程服务器，得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if nil!=err{
		fmt.Println("client start err,exit!")
		return
	}
	for{
		//2、连接调用Write写数据
		_, err := conn.Write([]byte("Hello Zinx V0.4..."))
		if nil!=err{
			fmt.Println("write conn err",err)
			return
		}
		buf := make([]byte,512)
		_, err = conn.Read(buf)
		if nil!=err{
			fmt.Println("read buf error")
			return
		}
		fmt.Printf(" server cal back: %s\n",buf)
		//cpu阻塞
		time.Sleep(1*time.Second)
	}

}
