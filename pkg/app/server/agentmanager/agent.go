package agentmanager

import (
	"fmt"
	"net"

	"openeluer.org/PilotGo/PilotGo/pkg/protocol"
)

type AgentMessageHandler func(*Agent, *protocol.Message) error

type Agent struct {
	UUID             string
	conn             net.Conn
	MessageProcesser *protocol.MessageProcesser
}

// 通过给定的conn连接初始化一个agent
func NewAgent(conn net.Conn) (*Agent, error) {
	agent := &Agent{
		conn:             conn,
		MessageProcesser: protocol.NewMessageProcesser(),
	}
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

func (a *Agent) StartListen() {
	go func() {
		a.startListen()
	}()
}

func (a *Agent) startListen() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("server processor panic error:", err)
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
	a.bindHandler(protocol.Heartbeat, func(a *Agent, msg *protocol.Message) error {
		fmt.Println("process heartbeat from processor:", a.conn.RemoteAddr(), string(msg.Body))
		return nil
	})

	return nil
}

// 远程在agent上运行脚本
func (a *Agent) RunScript() {

}

// 远程获取agent端的系统信息
func (a *Agent) GetOSInfo() {

}

// 远程获取agent端的系统信息
func (a *Agent) GetInfo() {

}

func Send(conn net.Conn, msg *protocol.Message) (error, error) {
	data := msg.Encode()
	sendData := protocol.TlvEncode(data)

	data_length := len(sendData)
	send_count := 0
	for {
		n, err := conn.Write(sendData[send_count:])
		if err != nil {
			return err, nil
		}
		if n+send_count >= data_length {
			send_count = send_count + n
			break
		}
	}
	return nil, nil
}
