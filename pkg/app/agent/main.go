package main

import (
  "fmt"
  "openeluer.org/PilotGo/PilotGo/pkg/app/agent/network"
  "openeluer.org/PilotGo/PilotGo/pkg/protocol"
  "openeluer.org/PilotGo/PilotGo/pkg/utils"
  "os"
  "time"
)

func main() {
	fmt.Println("Start PilotGo agent.")

	// 加载系统配置

	// 初始化日志

	// 与server握手
	client := &network.SocketClient{
		MessageProcesser: protocol.NewMessageProcesser(),
	}
	if err := client.Connect("localhost:8879"); err != nil {
		fmt.Println("connect server failed, error:", err)
		os.Exit(-1)
	}
	regitsterHandler(client)

	for {
		msg := &protocol.Message{
			Type: protocol.Heartbeat,
			Body: []byte(`{"type":1}`),
		}

		if err := client.Send(msg); err != nil {
			fmt.Println("send message failed, error:", err)
		}
		fmt.Println("send heartbeat message")

		time.Sleep(time.Second)

		// 接受远程指令并执行

		if false {
			break
		}
	}

	out, err := utils.RunCommand("date")
	if err == nil {
		fmt.Println(string(out))
	}

	fmt.Println("Stop PilotGo agent.")
}

func regitsterHandler(c *network.SocketClient) {
	c.BindHandler(protocol.Heartbeat, func(c *network.SocketClient, msg *protocol.Message) error {
		fmt.Println(string(msg.Body))
		return nil
	})

	c.BindHandler(protocol.RunScript, func(c *network.SocketClient, msg *protocol.Message) error {
		fmt.Println("process run script command:", string(msg.Body))
		return nil
	})
}
