package ziface

/*
	封包、拆包、模块
	直接面向TCP连接中的数据流，用于处理TCP粘包问题
*/

type IDataPack interface {
	GetHeadLen() uint32                //获取包的头的长度方法
	Pack(msg IMessage) ([]byte, error) //封包方法
	Unpack([]byte) (IMessage, error)   //拆包方法
}
