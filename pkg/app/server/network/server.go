package network

import (
	"fmt"
	"net"
)

type SocketServer struct {
	// MessageProcesser *protocol.MessageProcesser
	OnAccept func(net.Conn)
	OnStop   func()
}

func (s *SocketServer) Run(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}

		s.OnAccept(conn)
	}
}

func (s *SocketServer) Stop() {

}
