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
)

type Message struct {
	UUID   string `json:"message_uuid"`
	Type   int    `json:"message_type"`
	Status int    `json:"status"`
	Data   interface{}
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
