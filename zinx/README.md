11.11
ZinxV0.1-基础的server
方法：
    -启动服务器：基本的服务开发 1创建addr， 2创建listnner， 3处理客户端的基本的业务，回显功能
    -停止服务器：
    -运行服务器：调用Start()方法，调用之后做阻塞处理，在之间可以做今后的一个扩展功能
    -初始化server
属性：
    -name名称
    -监听的IP
    -监听的端口   

####################################
11.12
ZinxV0.2-简单的连接封装和业务绑定
连接的模块：
    -方法：
        1)启动连接Start()
        2)停止连接Stop()
        3)获取当前连接的conn对象(套接字)-->GetTCPConnection() *net.TCPConn
        4)得到链接ID-->GetConnID() uint32
        5)得到客户端连接的地址和端口-->RemoteAddr() net.Addr
        6)发送数据的方法Send()-->Send(data []byte) error
        7)连接所绑定的处理业务的函数类型-->type HandleFunc func(*net.TCPConn,[]byte,int) error
    -属性：
        1)socket TCP套接字-->Conn *net.TCPConn
        2)连接的ID-->ConnID uint32
        3)当前连接的状态(是否已经关闭)-->isClosed bool
        4)与当前连接所绑定的处理业务方法-->handleAPI ziface.HandleFunc
        5)等待连接被动退出的channel-->ExitChan chan bool
       
11.14
基础router模块：
    --Request请求封装：将连接和数据绑定在一起
        -属性：
            1)连接IConnection-->	conn ziface.IConnection
            2)请求数据-->data []byte
        -方法：
            1)得到当前链接-->func (r *Request)GetConnection() ziface.IConnection
            2)得到当前数据-->func (r *Request)GetData() []byte
    --Router模块：
        -抽象的IRouter：
            1)处理业务之前的方法-->Prehandle(request IRequest)
            2)处理业务的主方法-->Handle(request IRequest)
            3)处理业务之后的方法-->PostHandle(request IRequest)
        -具体的BaseRouter：
            1)处理业务之前的方法-->func (br *BaseRouter) Prehandle(request ziface.IRequest) {}
            2)处理业务的主方法-->func (br *BaseRouter) Handle(request ziface.IRequest) {}
            3)处理业务之后的方法 -->func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
            
11.15
zinx集成router模块：
    --IServer增添路由添加功能-->AddRouter(router IRouter)
    --Server类增添router成员-->HandleAPI去掉
    --Connection类绑定一个Router成员
    --在Connection调用已经注册的Router处理业务         
            
 使用ZinxV0.3开发：
    --1、创建一个server句柄，使用Zinx的api
      2、给当前zinx框架添加一个自定义的router
      3、启动server
    --需要继承BaseRouter，实现PreHandle、Handle、PostHandle这3个方法
#############################################################################    
11.15
ZinxV0.4全局配置
    --服务器应用/conf/zinx.json(用户进行填写)
    --创建一个zinx的全局配置模块utils/globalobj.go ：①init方法读取用户配制好的zinx.json文件-->globalobj对象中; ②提供一个全局的GlobalObject对象
    --将zinx框架中全部的硬代码，用globalobj里面的参数进行替换
    --使用ZinxV0.4开发:基于Zinx的服务器应用程序conf/zinx.json

11.16
消息封装：
    --定义一个消息的结构
        -属性：
            1、消息的ID
            2、消息的长度
            3、消息的内容
        -方法：setter,  getter    
    --定义一个解决TCP粘包问题的封包拆包的模块：
        -针对Message进行TLV格式的封装：func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error)           TLV:
            1、写Message的长度
            2、写Message的ID
            3、写Message的内容
        -针对Message进行TLV格式的拆包：func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error)
            1、先读取固定长度的head-->消息内容的长度和消息的类型
            2、再根据消息内容的长度，再进行一次读写，从conn中读取消息的内容
    --将消息封装机制集成到Zinx框架中：
        1、将Message添加到Request属性中
        2、修改链接读取数据的机制，将之前的单纯的读取byte改成拆包形式的读取按照TLV形式读取
    --使用ZinxV0.5开发

11.17
多路由模式
    --消息管理模块(支持多路由业务api调度管理)--MsgHandler
        -属性：
            1、集合-消息ID和对应的router的关系-map -->Apis map[uint32] ziface.IRouter
        -方法：
            1、根据msgID来索引调度路由方法-->func (mh *MsgHandle)DoMsgHandler(request ziface.IRequest)
            2、添加路由方法到map集合中-->func (mh *MsgHandle)AddRouter(msgID uint32, router ziface.IRouter)
    --消息管理模块集成到Zinx框架中:
        1、将server模块中的Router属性替换成MsgHandler属性
        2、将server之前的AddRouter修改，调用MsgHandler的AddRouter-->AddRouter(msgID uint32, router IRouter)
        3、将connection模块Router属性 替换成MsgHandler， 修改初始化Connection方法
        4、Connection的之前调度Router的业务替换成MsgHandler调度，修改StartReader方法
    --使用ZinxV0.6开发

11.20
读写协程分离
    --添加一个Reader和Writer之间通信的channel
    --添加一个Writer goroutine
    --Reader 由之前直接发送给客户端 改成 发送给通信channel
    --启动Reader和Writer一同工作
    --使用ZinxV0.7开发








   看29集
     
 