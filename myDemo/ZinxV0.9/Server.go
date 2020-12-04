package main

import (
	"Myzinx/zinx/ziface"
	"Myzinx/zinx/znet"
	"fmt"
)

/*
	基于Zinx框架来开发的服务器端应用程序
*/

//ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

//Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle...")
	//先读取客户端的数据，再回写ping..ping...ping
	fmt.Println("recv from client: msgID=", request.GetMsgID(), ",data=", string(request.GetData()))
	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping"))
	if nil != err {
		fmt.Println(err)
	}
}

//hello Zinx test 自定义路由
type HelloZinxRouter struct {
	znet.BaseRouter
}

//Test Handle
func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle...")
	//先读取客户端的数据，再回写ping..ping...ping
	fmt.Println("recv from client: msgID=", request.GetMsgID(), ",data=", string(request.GetData()))
	err := request.GetConnection().SendMsg(201, []byte("Hello Welcome to Zinx!!"))
	if nil != err {
		fmt.Println(err)
	}
}

//创建连接之后执行钩子函数
func DoConnectionBegin(conn ziface.IConnection)  {
	fmt.Println("===>DoConnectionBegin is Called...")
	if err := conn.SendMsg(202,[]byte("DoConnection BEGIN"));nil!=err{
		fmt.Println(err)
	}
}

//连接断开之前的需要执行的函数
func DoConnectionLost(conn ziface.IConnection)  {
	fmt.Println("===> DoConnectionLost is Called...")
	fmt.Println("conn ID = ",conn.GetConnID(),"is Lost...")
}

func main() {
	//1.创建一个server句柄，使用Zinx的api
	s := znet.NewServer("[zinx V0.7]") //s返回的是IServer的接口

	//2.注册Hook钩子函数(只是注册，并未执行，执行在start里)
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//3.给当前zinx框架添加一个自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})

	//4.启动server
	s.Serve()
}
