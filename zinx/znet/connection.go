package znet

import (
	"Myzinx/zinx/ziface"
	"errors"
	"fmt"
	"io"
	"net"
)

/*
	连接模块
*/
type Connection struct {
	Conn       *net.TCPConn      //当前连接的socketTCP套接字
	ConnID     uint32            //连接的ID
	isClosed   bool              //当前的连接状态
	ExitChan   chan bool         //告知当前连接已经退出的/停止 channel(油Reader告知Writer退出)
	msgChan    chan []byte       //无缓冲管道，用于读、写goroutine之间的消息通信
	MsgHandler ziface.IMsgHandle //消息的管理MsgID和对应的处理业务API关系

}

//初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, msghandle ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
		MsgHandler: msghandle,
	}
	return c
}

//连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running..]")
	defer fmt.Println("[Reader is exist!], ConnID=", c.ConnID, " remote addr is ", c.RemoteAddr().String())
	defer c.Stop()
	for {
		//读取客户端的数据到buf中，最大512字节
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if nil != err {
		//	fmt.Println("recv buf err", err)
		//	continue
		//}

		//创建一个拆包解包对象
		dp := NewDataPack()

		//读取客户端Msg Head 二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); nil != err {
			fmt.Println("read msg head error", err)
			break
		}

		//拆包，得到msgID和msgDataLen放在msg消息中
		msg, err := dp.Unpack(headData)
		if nil != err {
			fmt.Println("unpack error", err)
			break
		}
		//根据dataLen 再次读取Data，放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); nil != err {
				fmt.Println("read msg data error:", err)
				break
			}
		}
		msg.SetData(data)

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		//从路由中，找到注册绑定的conn对应的router调用
		//根据绑定好的MsgID找到对应处理api业务 执行
		go c.MsgHandler.DoMsgHandler(&req)
	}
}

//写消息goroutine，专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer goroutine is running]")
	defer fmt.Println("[conn Writer exit!]",c.RemoteAddr().String())

	//不断的阻塞的等待channel的消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); nil != err {
				fmt.Println("Send data error:", err)
				return
			}
		case <-c.ExitChan:
			//代表Reader已经推出，此时Writer也要退出
			return

		}
	}

}

//启动连接，让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start().. ConnID=", c.ConnID)
	//启动当前连接的读数据的业务
	go c.StartReader()
	//启动从当前连接写数据的业务
	go c.StartWriter()
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
	//告知Writer关闭
	c.ExitChan <- true
	//回收资源
	close(c.ExitChan)
	close(c.msgChan)
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

//提供一个SendMsg方法，将我们要发送给客户端的数据，先进行封包，再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}
	//将data进行封包 MsgDataLen| MsgID| Data
	dp := NewDataPack()
	//MsgDataLen|MsgID|Data
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if nil != err {
		fmt.Println("Pack error msg id=", msgId)
		return errors.New("Pack error msg")
	}
	//将数据发送给客户端
	c.msgChan <- binaryMsg

	return nil
}
