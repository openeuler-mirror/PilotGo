package net

import (
	"net"
)

func Send(conn net.Conn, data []byte) error {
	data_length := len(data)
	send_count := 0
	for {
		n, err := conn.Write((data)[send_count:])
		if err != nil {
			return err
		}
		if n+send_count >= data_length {
			send_count = send_count + n
			break
		}
	}
	return nil
}
