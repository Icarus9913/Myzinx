package znet

import (
	"Myzinx/zinx/utils"
	"Myzinx/zinx/ziface"
	"bytes"
	"encoding/binary"
	"errors"
)

//封包，拆包的具体模块
type DataPack struct {
}

//拆包封包实例的一个初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

//获取包的头的长度方法
func (dp *DataPack) GetHeadLen() uint32 {
	//DataLen uint32(4字节) + ID uint32(4字节)
	return 8
}

//封包方法
//dataLen| msgID|data
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuffer := bytes.NewBuffer([]byte{})

	//将dataLen写进dataBuff中
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgLen()); nil != err {
		return nil, err
	}

	//将MsgId写进dataBuff中
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgId()); nil != err {
		return nil, err
	}

	//将data写进dataBuff中
	if err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetData()); nil != err {
		return nil, err
	}

	return dataBuffer.Bytes(), nil
}

//拆包方法  (将包的Head信息读出来) 之后再根据head信息里的data的长度，再进行一次读
func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuffer := bytes.NewBuffer(binaryData)

	//只解压head信息，得到dataLen和MsgID
	msg := &Message{}

	//读dataLen
	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.DataLen); nil != err {
		return nil, err
	}

	//读MsgID
	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.Id); nil != err {
		return nil, err
	}

	//判断dataLen是否已经超出了我们允许的最大包长度
	if utils.GlobalObject.MaxPackageSize>0 && msg.DataLen>utils.GlobalObject.MaxPackageSize{
		return nil, errors.New("too Large msg data recv!")
	}
	return msg, nil
}
