package agentmanager

import (
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
	pnet "openeluer.org/PilotGo/PilotGo/pkg/net"
	"openeluer.org/PilotGo/PilotGo/pkg/protocol"
)

type AgentMessageHandler func(*Agent, *protocol.Message) error

type Agent struct {
	UUID             string
	Version          string
	conn             net.Conn
	MessageProcesser *protocol.MessageProcesser
}

// 通过给定的conn连接初始化一个agent并启动监听
func NewAgent(conn net.Conn) (*Agent, error) {
	agent := &Agent{
		UUID:             "agent",
		conn:             conn,
		MessageProcesser: protocol.NewMessageProcesser(),
	}

	go func() {
		agent.startListen()
	}()

	if err := agent.Init(); err != nil {
		return nil, err
	}

	return agent, nil
}

func (a *Agent) bindHandler(t int, f AgentMessageHandler) {
	a.MessageProcesser.BindHandler(t, func(c protocol.MessageContext, msg *protocol.Message) error {
		return f(c.(*Agent), msg)
	})
}

func (a *Agent) startListen() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("server processor panic error:%s", err.(error).Error())
			a.conn.Close()
		}
	}()

	readBuff := []byte{}
	for {
		buff := make([]byte, 1024)
		n, err := a.conn.Read(buff)
		if err != nil {
			fmt.Println("read error:", err)
			return
		}
		readBuff = append(readBuff, buff[:n]...)

		// 切割frame
		i, f := protocol.TlvDecode(&readBuff)
		if i != 0 {
			readBuff = readBuff[i:]

			msg := protocol.ParseMessage(*f)
			a.MessageProcesser.ProcessMessage(a, msg)
		}
	}
}

// 远程获取agent端的信息进行初始化
func (a *Agent) Init() error {
	// TODO: 此处绑定所有的消息处理函数
	a.bindHandler(protocol.Heartbeat, func(a *Agent, msg *protocol.Message) error {
		logger.Info("process heartbeat from processor, remote addr:%s, data:%s",
			a.conn.RemoteAddr().String(), msg.String())
		return nil
	})

	a.bindHandler(protocol.AgentInfo, func(a *Agent, msg *protocol.Message) error {
		logger.Info("process heartbeat from processor, remote addr:%s, data:%s",
			a.conn.RemoteAddr().String(), msg.String())
		return nil
	})

	data, err := a.AgentInfo()
	if err != nil {
		logger.Error("fail to get agent info, address:%s", a.conn.RemoteAddr().String())
	}
	d := data.(map[string]interface{})
	logger.Debug("response agent info is %v", d)
	a.UUID = d["agent_uuid"].(string)

	a.Version = d["agent_version"].(string)

	// add agent to agent manager
	AddAgent(a)

	return nil
}

// 远程在agent上运行脚本
func (a *Agent) RunScript(cmd string) (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.RunScript,
		Data: struct {
			Command string
		}{
			Command: cmd,
		},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

func (a *Agent) sendMessage(msg *protocol.Message, wait bool, timeout time.Duration) (*protocol.Message, error) {
	logger.Debug("send message:%s", msg.String())

	if msg.UUID == "" {
		msg.UUID = uuid.New().String()
	}

	if wait {
		waitChan := make(chan *protocol.Message)
		a.MessageProcesser.WaitMap.Store(msg.UUID, waitChan)

		pnet.Send(a.conn, protocol.TlvEncode(msg.Encode()))

		// wail for response
		data := <-waitChan
		return data, nil
	}

	// just send data
	return nil, pnet.Send(a.conn, protocol.TlvEncode(msg.Encode()))
}

// 远程获取agent端的系统信息
func (a *Agent) GetOSInfo() (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.OsInfo,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 远程获取agent端的CPU信息
func (a *Agent) GetCPUInfo() (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.CPUInfo,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 远程获取agent端的内存信息
func (a *Agent) GetMemoryInfo() (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.MemoryInfo,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 远程获取agent端的内核信息
func (a *Agent) GetSysctlInfo() (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.SysctlInfo,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 临时修改agent端系统参数
func (a *Agent) ChangeSysctl(args string) (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.SysctlChange,
		Data: args,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 查看某个内核参数的值
func (a *Agent) SysctlView(args string) (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.SysctlView,
		Data: args,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 远程获取agent端的系统信息
func (a *Agent) AgentInfo() (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.AgentInfo,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}
