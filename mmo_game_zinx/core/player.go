package core

import (
	"Myzinx/zinx/ziface"
	"fmt"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"sync"
)

//玩家对象
type Player struct {
	Pid  int32              //玩家ID
	Conn ziface.IConnection //当前玩家的连接(用于和客户端的连接)
	X    float32            //平面的x坐标
	Y    float32            //高度
	Z    float32            //平面y坐标(注意不是Y)
	V    float32            //旋转的0-360角度
}

/*
	Player ID生成器
*/

var PidGen int32 = 1  //用来生成玩家ID的计数器
var IdLock sync.Mutex //保护PidGen的Mutex

//创建一个玩家的方法
func NewPlayer(conn ziface.IConnection) *Player {
	//生成一个玩家ID
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()

	//创建一个玩家对象
	p := &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)), //随机在160坐标点 基于X轴若干偏移
		Y:    0,
		Z:    float32(140 + rand.Intn(20)), //随机在140坐标点，基于Y轴若干偏移
		V:    0,                            //角度为0
	}
	return p
}

/*
	提供一个发送给客户端消息的方法
	主要是将pb的protobuf数据序列化之后，再调用zinx的SendMsg方法
*/
func (p *Player)SendMsg(msgId uint32,data proto.Message)  {
	//将proto Message结构体序列化 转换成二进制
	msg ,err := proto.Marshal(data)
	if nil!=err{
		fmt.Println("marshal msg err:",err)
		return
	}

	//将二进制文件 通过zinx框架的sendmsg将数据发送给客户端
	if nil==p.Conn{
		fmt.Println("connection in player is nil")
		return
	}
	if err := p.Conn.SendMsg(msgId, msg);nil!=err{
		fmt.Println("Player SendMsg error!")
		return
	}
	return
}

//告知客户端玩家Pid，同步已经生成的玩家ID给客户端

