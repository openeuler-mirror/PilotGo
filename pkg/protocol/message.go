package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
)

const (
	Heartbeat       = 1
	GetConfig       = 2
	WriteConfig     = 3
	PackageInstall  = 4
	PackageUpdate   = 5
	PackageRollback = 6
	RunScript       = 7
	AgentUpdate     = 8
	AgentUninstall  = 9
)

type Message struct {
	Type int
	Body []byte
}

type MessageContext interface{}

type MessageHandler func(MessageContext, *Message) error

type MessageProcesser struct {
	in         <-chan *Message
	handlerMap map[int]MessageHandler
}

func (m *Message) Encode() []byte {

	type32 := int32(m.Type)

	buffer := bytes.NewBuffer([]byte{})
	if err := binary.Write(buffer, binary.BigEndian, &type32); err != nil {
		fmt.Println("parse lenth error", err)
		return []byte{}
	}
	resultBytes := []byte{}
	resultBytes = append(resultBytes, buffer.Bytes()...)
	resultBytes = append(resultBytes, []byte(m.Body)...)

	return resultBytes
}

func NewMessageProcesser() *MessageProcesser {
	m := &MessageProcesser{}
	m.Init()
	return m
}

func (m *MessageProcesser) Init() {
	m.in = make(<-chan *Message)
	m.handlerMap = map[int]MessageHandler{}
}

func (m *MessageProcesser) BindHandler(t int, f MessageHandler) {
	m.handlerMap[t] = f
}

func (m *MessageProcesser) ProcessMessage(ctx MessageContext, msg *Message) error {
	if f, ok := m.handlerMap[msg.Type]; ok {
		return f(ctx, msg)
	} else {
		return errors.New("unregistered message:" + strconv.Itoa(msg.Type) + string(msg.Body))
	}
}

// 从字节数组中解析构建一个message
func ParseMessage(data []byte) *Message {
	typeBytes := data[:4]
	t32 := int32(0)
	buf := bytes.NewBuffer(typeBytes)
	binary.Read(buf, binary.BigEndian, &t32)
	t := int(t32)

	return &Message{
		Type: t,
		Body: data[4:],
	}
}
