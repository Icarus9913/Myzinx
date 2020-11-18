package znet

import (
	"Myzinx/zinx/utils"
	"Myzinx/zinx/ziface"
	"fmt"
	"net"
)

//iServer的接口实现，定义一个Server的服务器模块
type Server struct {
	Name      string         //服务器的名称
	IPVersion string         //服务器绑定的ip版本
	IP        string         //服务器监听的IP
	Port      int            //服务器监听的端口
	Router    ziface.IRouter //当前的Server添加一个router，server注册的连接对应的处理业务
}



//启动服务器
func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name: %s,listener at IP:%s, Port:%d is starting\n",
		utils.GlobalObject.Name,utils.GlobalObject.Host,utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version:%s, MaxConn:%d, MaxPackageSize:%d\n",
		utils.GlobalObject.Version,utils.GlobalObject.MaxConn,utils.GlobalObject.MaxPackageSize)

	go func() {
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
			//将处理新链接的业务方法和conn进行绑定，得到我们的连接模块
			dealConn := NewConnection(conn, cid, s.Router)
			cid++
			//启动当前的连接业务处理
			go dealConn.Start()
		}
	}()

}

//停止服务器
func (s *Server) Stop() {
	//TODO 将一些服务器的资源、状态或者一些已经开辟的连接信息 进行停止或者回收
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
func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router=router
	fmt.Println("Add Router success!")
}

/*
	初始化Server模块的方法
*/
func NewServer() ziface.IServer {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		Router: nil,
	}
	return s
}
