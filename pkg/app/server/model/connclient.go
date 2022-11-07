package model

import "github.com/gorilla/websocket"

type ConnClient struct {
	Conn *websocket.Conn
}
