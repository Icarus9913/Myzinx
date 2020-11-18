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
	fmt.Println("Call Router Handle...")
	//先读取客户端的数据，再回写ping..ping...ping
	fmt.Println("recv from client: msgID=", request.GetMsgID(), ",data=", string(request.GetData()))
	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if nil != err {
		fmt.Println(err)
	}
}

func main() {
	//1.创建一个server句柄，使用Zinx的api
	s := znet.NewServer() //s返回的是IServer的接口
	//2.给当前zinx框架添加一个自定义的router
	s.AddRouter(&PingRouter{})

	//3.启动server
	s.Serve()
}
