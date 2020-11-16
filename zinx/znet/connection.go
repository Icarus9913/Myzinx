package znet

import (
	"Myzinx/zinx/utils"
	"Myzinx/zinx/ziface"
	"fmt"
	"net"
)

/*
	连接模块
*/
type Connection struct {
	Conn     *net.TCPConn   //当前连接的socketTCP套接字
	ConnID   uint32         //连接的ID
	isClosed bool           //当前的连接状态
	ExitChan chan bool      //告知当前连接已经退出的/停止 channel
	Router   ziface.IRouter //该链接处理的方法Router
}

//初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		ExitChan: make(chan bool, 1),
		Router:   router,
	}
	return c
}

//连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running..")
	defer fmt.Println("ConnID=", c.ConnID, " Reader is exist, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()
	for {
		//读取客户端的数据到buf中，最大512字节
		buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		_, err := c.Conn.Read(buf)
		if nil != err {
			fmt.Println("recv buf err", err)
			continue
		}

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			data: buf,
		}
		//执行注册的路由方法
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
		//从路由中，找到注册绑定的conn对应的router调用
	}
}

//启动连接，让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start().. ConnID=", c.ConnID)
	//启动当前连接的读数据的业务
	go c.StartReader()
	//TODO 启动从当前连接写数据的业务
}

//停止连接，结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop().. ConnID=", c.ConnID)
	//如果当前连接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	//关闭socket连接
	c.Conn.Close()
	//回收资源
	close(c.ExitChan)
}

//获取当前连接的绑定socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前连接模块的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端的TCP状态IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//发送数据，将数据发送给远程的客户端
func (c *Connection) Send(data []byte) error {
	return nil
}
