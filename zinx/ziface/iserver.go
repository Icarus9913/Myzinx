package ziface

//定义一个服务器接口
type IServer interface {
	Start()                   //启动服务器
	Stop()                    //停止服务器
	Serve()                   //运行服务器
	AddRouter(router IRouter) //路由功能：给当前的服务注册一个路由方法，供客户端的连接处理使用
}
