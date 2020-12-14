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

11.20
消息队列及多任务
    --消息队列及Worker工作池实现
        --1、创建一个消息队列:
            --MsgHandler消息管理模块
                -属性：
                    1)消息队列--TaskQueue []chan ziface.IRequest
                    2)worker工作池的数量(WorkerPoolSize uint32)--->在全局配置的参数中获取，也可以在配置文件中让用户设置
        --2、创建多任务worker的工作池并且启动-->创建一个worker的工作池func (mh *MsgHandle) StartWorkerPool():
                                        1)根据workerPoolSize的数量去创建Worker--func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest)
                                        2)每个worker都应该用一个Go去承载
                                            -1、阻塞的等待与当前worker对应的channel的消息
                                            -2、一旦有消息到来，worker应该处理当前消息对应的业务，调用DoMsgHandler()
        --3、将之前的发送消息，全部改成 把消息发送给 消息队列和worker工作池来处理
            --定义一个方法，将消息发送给消息队列工作池的方法-->func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest)：
                1、保证每个worker所收到的request任务是均衡(平均分配)，让哪个worker去处理，只需要将这个request请求发送给对应的taskQueue即可
                2、将消息直接发送给对应的channel
    --将消息队列机制集成到Zinx框架中
        --1、开启并调用消息队列及Worker工作池-->保证WorkerPool只有一个，应该在创建Server模块的时候开启(在server listen之前添加)
        --2、将从客户端处理的消息，发送给当前的Worker工作池来处理-->在已经处理完拆包，得到了request请求，交给工作池来处理
    --使用ZinxV0.8开发

11.22
连接管理
    --创建一个连接管理模块(定义、属性、方法)
        --ConnManager
            -属性：
                1)已经创建的Connection集合map-->connections map[uint32] ziface.IConnection
                2)针对map的互斥锁-->connLock    sync.RWMutex
            -方法：
                1)添加链接-->func (connMgr *ConnManager) Add(conn ziface.IConnection)
                2)删除链接-->func (connMgr *ConnManager) Remove(conn ziface.IConnection)
                3)根据连接ID查找对应的连接-->func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error)
                4)总连接个数-->func (connMgr *ConnManager) Len() int
                5)清理全部的连接-->func (connMgr *ConnManager) ClearConn()
    --将连接管理模块集成到Zinx框架中
        --将ConnManager加入Server模块中
            -给server添加一个ConnMgr属性
            -修改NewServer方法，加入ConnMgr初始化
            -判断当前的链接数量是否已经超出最大值MaxConn
        --每次成功与客户端建立连接后-添加链接到ConnManager中
            -在NewConnection的时候将新的conn加入到ConnMgr中，需要给Connection加入隶属server属性-->给Server提供一个GetConnMgr方法
        --每次与客户端连接断开后，将连接从ConnManager删除
            -在Conn.Stop()方法中，将当前的连接从ConnMgr删除即可
            -当server停止的时候应该清除所有的连接 Stop()方法中加入ConnMgr.ClearConn()
    --给Zinx框架提供 创建连接之后/销毁连接之前 所要处理的一些业务 提供给用户能够注册Hook函数
        -属性
            --该Server创建连接之后自动调用Hook函数-->OnConnStart func(conn ziface.IConnection)
            --该Server销毁连接之前自动调用的Hook函数-->OnConnStop  func(conn ziface.IConnection)
        -方法
            --注册OnConnStart钩子函数方法-->func (s *Server)SetOnConnStart(hookFunc func(connection ziface.IConnection))
            --注册OnConnStop钩子方法-->func (s *Server)SetOnConnStop(hookFunc func(connection ziface.IConnection))
            --调用OnConnStart钩子函数方法-->func (s *Server)CallOnConnStart(conn  ziface.IConnection)
            --调用OnConnStop钩子方法-->func (s *Server)CallOnConnStop(conn  ziface.IConnection)
        -在Conn创建之后调用OnConnStart-->在conn.Start()中调用
        -在Conn销毁之前调用OnConnStop-->在conn.Stop()中调用    
    --使用ZinxV0.9开发
        -注册Hook钩子函数:
        	s.SetOnConnStart(DoConnectionBegin)
        	s.SetOnConnStop(DoConnectionLost)
        	
11.29        	
连接属性配置
    --给Connection模块添加可以配置属性的功能
        -属性
            --连接属性集合map-->property  map[string]interface{}
            --保护连接属性的锁-->propertyLock sync.RWMutex
        -方法
            --设置连接属性-->func (c *Connection)SetProperty(key string, value interface{})
            --获取连接属性-->func (c *Connection)GetProperty(key string) (interface{}, error)
            --移除连接属性-->func (c *Connection)RemoveProperty(key string)

12.13
基于Zinx的服务器应用--MMO多人在线网游
    --协议的定义
    --AOI兴趣点的算法
        --AOI格子的数据类型-Gid
            -属性
                --格子ID-->GID int
                --格子的左边边界坐标-->MinX int
                --格子的右边边界坐标-->MaxX int
                --格子的上边边界坐标-->MinY int
                --格子的下边边界坐标-->MaxY int
                --当前格子内玩家或者物体成员的ID集合-->playerIDs map[int]bool
                --保护当前集合的锁-->pIDLock sync.RWMutex
            -方法
                --初始化当前的格子的方法-->func NewGrid(gID, minX, maxX, minY, maxY int) *Grid
                --给格子添加一个玩家-->func (g *Grid) Add(playerID int)
                --从格子中删除一个玩家-->func (g *Grid) Remove(playerID int)
                --得到当前格子中所有的玩家-->func (g *Grid) GetPlayerIDs() (playerIDs []int)
                --调试使用-打印出格子的基本信息-->func (g *Grid) String() string
        --AOI管理格子(地图)数据类型-AOIManager
            -属性
                --区域的左边界坐标-->MinX int
                --区域的右边界坐标-->MaxX int
                --X方向格子的数量-->CntsX int
                --区域的上边界坐标-->MinY int
                --区域的下边界坐标-->MaxY int
                --Y方向格子的数量-->CntsY int
                --当前区域中有哪些格子-->map-key=格子的ID，value=格子对象-->grids map[int]*Grid
            -方法
                --初始化一个AOI管理区域模块-->func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager
                --调试使用--打印当前AOI模块-->func (m *AOIManager) String() string
                --得到每个格子在X轴方向的宽度-->func (m *AOIManager) gridWidth() int
                --得到每个格子在Y轴方向的长度-->func (m *AOIManager) gridLength() int
                --添加一个PlayerID到一个格子中
                --移除一个格子中的PlayerID
                --通过GID获取全部的PlayerID
                --通过坐标将Player添加到一个格子中
                --通过坐标把一个Player从一个格子中删除
                --通过Player坐标得到当前Player周边九宫格内全部PlayerIDs
                --通过坐标获取对应的玩家所在的GID
    --数据传输协议Protocol buffer
    --玩家的业务
        -玩家上线
        -世界聊天
        -上线位置的信息同步(玩家上线广播)
        -移动位置与广播
        -玩家下线


   看55集
     
 