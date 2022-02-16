package protocol

import (
	"encoding/json"
	"errors"
	"strconv"
	"sync"

	"openeluer.org/PilotGo/PilotGo/pkg/logger"
)

const (
	// agent发送的心跳信息
	Heartbeat = 1
	// 获取agent返回kernel配置
	GetKernelConfig = 2
	// 修改agent端kernel配置
	WriteKernelConfig = 3
	// agent端软件包信息
	PackageInfo = 4
	// agent端软件包安装
	PackageInstall = 5
	// agent端软件包升级
	PackageUpdate = 6
	// agent端软件包回滚
	PackageRollback = 7
	// agent端执行shell脚本
	RunScript = 8
	// agent升级
	AgentUpdate = 9
	// agent卸载
	AgentUninstall = 10
	// agent信息获取
	AgentInfo = 11
	// os信息获取
	OsInfo = 12
	//CPU信息数据获取
	CPUInfo = 13
	//内存信息数据获取
	MemoryInfo = 14
	//获取磁盘的IO信息
	DiskInfo = 15
	//内核配置信息数据获取
	SysctlInfo = 16
	//临时修改系统参数
	SysctlChange = 17
	//查看某个内核参数的值
	SysctlView = 18
	//查看服务列表
	ServiceList = 19
	//查看某个服务状态
	ServiceStatus = 20
	//重启某个服务
	ServiceRestart = 21
	//启动某个服务
	ServiceStart = 22
	//关闭某个服务
	ServiceStop = 23
	//获取全部安装的rpm包列表
	AllRpm = 24
	//获取源软件包名以及源
	RpmSource = 25
	//获取软件包信息
	RpmInfo = 26
	//安装rpm软件包
	InstallRpm = 27
	//卸载rpm软件包
	RemoveRpm = 28
	//获取磁盘的使用情况
	DiskUsage = 29
	//创建挂载磁盘的目录
	CreateDiskPath = 30
	//挂载磁盘
	DiskMount = 31
	//卸载磁盘
	DiskUMount = 32
	//磁盘格式化
	DiskFormat = 33
	// 获取当前TCP网络连接信息
	NetTCP = 34
	//获取当前UDP网络连接信息
	NetUDP = 35
	//获取网络读写字节／包的个数
	NetIOCounter = 36
	// 获取网卡配置
	NetNICConfig = 37
	// 获取当前用户信息
	CurrentUser = 38
	// 获取所有用户的信息
	AllUser = 39
	// 创建新的用户，并新建家目录
	AddLinuxUser = 40
	// 删除用户
	DelUser = 41
	// chmod [-R] 权限值 文件名
	ChangePermission = 42
	// chown [-R] 所有者 文件或目录
	ChangeFileOwner = 43
)

type Message struct {
	UUID   string `json:"message_uuid"`
	Type   int    `json:"message_type"`
	Status int    `json:"status"`
	Data   interface{}
	Error  string
}

type MessageContext interface{}

type MessageHandler func(MessageContext, *Message) error

type MessageProcesser struct {
	// in         <-chan *Message

	// 用于绑定默认消息处理函数
	handlerMap map[int]MessageHandler
	// 用于阻塞型的消息发送
	WaitMap sync.Map
}

func (m *Message) Encode() []byte {

	// type32 := int32(m.Type)

	// buffer := bytes.NewBuffer([]byte{})
	// if err := binary.Write(buffer, binary.BigEndian, &type32); err != nil {
	// 	fmt.Println("parse lenth error", err)
	// 	return []byte{}
	// }
	// resultBytes := []byte{}
	// resultBytes = append(resultBytes, buffer.Bytes()...)
	// bytes, err := json.Marshal(m)
	// if err != nil {
	// 	logger.Error("marshal message to json failed, message is:%v", m)
	// }
	// resultBytes = append(resultBytes, []byte(bytes)...)

	bytes, err := json.Marshal(m)
	if err != nil {
		logger.Error("marshal message to json failed, message is:%v", m)
	}

	return bytes
}

func (m *Message) String() string {
	bytes, err := json.Marshal(m)
	if err != nil {
		logger.Error("marshal message to json failed, message is:%v", m)
	}

	return string(bytes)
}

func NewMessageProcesser() *MessageProcesser {
	m := &MessageProcesser{
		handlerMap: map[int]MessageHandler{},
		WaitMap:    sync.Map{},
	}
	return m
}

func (m *MessageProcesser) BindHandler(t int, f MessageHandler) {
	m.handlerMap[t] = f
}

func (m *MessageProcesser) ProcessMessage(ctx MessageContext, msg *Message) error {
	// 如果message uuid在等待队列当中
	value, ok := m.WaitMap.Load(msg.UUID)
	if ok {
		waitChan := value.(chan *Message)
		waitChan <- msg
		return nil
	}

	if f, ok := m.handlerMap[msg.Type]; ok {
		return f(ctx, msg)
	} else {
		return errors.New("unregistered message:" + strconv.Itoa(msg.Type))
	}
}

// 从字节数组中解析构建一个message
func ParseMessage(data []byte) *Message {
	// typeBytes := data[:4]
	// t32 := int32(0)
	// buf := bytes.NewBuffer(typeBytes)
	// binary.Read(buf, binary.BigEndian, &t32)
	// t := int(t32)

	msg := &Message{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		logger.Error("unmarshal message error, error:%s, data:%s", err.Error(), string(data[4:]))
	}

	return msg
}
