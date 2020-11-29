package znet

import (
	"Myzinx/zinx/utils"
	"Myzinx/zinx/ziface"
	"fmt"
	"net"
)

//iServer的接口实现，定义一个Server的服务器模块
type Server struct {
	Name        string                         //服务器的名称
	IPVersion   string                         //服务器绑定的ip版本
	IP          string                         //服务器监听的IP
	Port        int                            //服务器监听的端口
	MsgHandler  ziface.IMsgHandle              //当前server的消息管理模块，用来绑定MsgID和对应的处理业务API关系
	ConnMgr     ziface.IConnManager            //该server的连接管理器
	OnConnStart func(conn ziface.IConnection) //该Server创建连接之后自动调用Hook函数-->OnConnStart()
	OnConnStop  func(conn ziface.IConnection) //该Server销毁连接之前自动调用的Hook函数-->OnConnStop()
}

//启动服务器
func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name: %s,listener at IP:%s, Port:%d is starting\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version:%s, MaxConn:%d, MaxPackageSize:%d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

	go func() {
		//0、开启消息队列及Worker工作池
		s.MsgHandler.StartWorkerPool()

		//1.获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if nil != err {
			fmt.Println("resolve tcp addr error:", err)
		}
		//2.监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if nil != err {
			fmt.Println("listen", s.IPVersion, "err:", err)
			return
		}
		fmt.Println("start Zinx server success,", s.Name, "success Listenning..")
		var cid uint32
		cid = 0

		//3.阻塞等待客户端连接，处理客户端连接业务(读写)
		for {
			//如果有客户端连接过来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if nil != err {
				fmt.Println("Accept err", err)
				continue
			}

			//设置最大连接个数的判断，如果超过最大连接，那么则关闭此新的连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				//TODO 给客户端响应一个超出最大连接的错误包
				fmt.Println("====> Too Many Connections MaxConn=", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			//将处理新链接的业务方法和conn进行绑定，得到我们的连接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++
			//启动当前的连接业务处理
			go dealConn.Start()
		}
	}()

}

//停止服务器
func (s *Server) Stop() {
	//TODO 将一些服务器的资源、状态或者一些已经开辟的连接信息 进行停止或者回收
	fmt.Println("[STOP Zinx server name]", s.Name)
	s.ConnMgr.ClearConn()
}

//运行服务器
func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()

	//TODO 做一些启动服务器之后的额外业务

	//阻塞状态
	select {}
}

//路由功能：给当前的服务注册一个路由方法，供客户端的连接处理使用
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router success!")
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr

}

/*
	初始化Server模块的方法
*/
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}
//注册OnConnStart钩子函数方法
func (s *Server)SetOnConnStart(hookFunc func(connection ziface.IConnection)){
	s.OnConnStart = hookFunc
}

//注册OnConnStop钩子方法
func (s *Server)SetOnConnStop(hookFunc func(connection ziface.IConnection)){
	s.OnConnStop = hookFunc
}

//调用OnConnStart钩子函数方法
func (s *Server)CallOnConnStart(conn  ziface.IConnection){
	if s.OnConnStart!=nil{
		fmt.Println("----->Call OnConnStart()...")
		s.OnConnStart(conn)
	}
}

//调用OnConnStop钩子方法
func (s *Server)CallOnConnStop(conn  ziface.IConnection){
	if s.OnConnStop!=nil{
		fmt.Println("----->Call OnConnStop()...")
		s.OnConnStop(conn)
	}
}