package main

import (
	"Myzinx/zinx/ziface"
	"Myzinx/zinx/znet"
)

//当前客户端建立连接之后的hook函数
func OnConnectionAdd(conn ziface.IConnection)  {
	//创建一个Player对象


	//给客户端发送MsgID

	//


}

func main()  {
	//创建zinx server句柄
	s := znet.NewServer("MMO Game Zinx")

	//连接创建和销毁的HOOK钩子函数


	//注册一些路由业务

	//启动服务
	s.Serve()

}
