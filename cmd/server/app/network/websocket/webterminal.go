/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package websocket

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

const (
	MsgData   = '1'
	MsgResize = '2'
)

type Resize struct {
	Columns int
	Rows    int
}
type Terminal struct {
	StdinPipe io.WriteCloser
	Session   *ssh.Session
	WsConn    *websocket.Conn
}

func NewTerminal(ws_conn *websocket.Conn, sshClient *ssh.Client) (*Terminal, error) {
	session, err := sshClient.NewSession()
	if err != nil {
		return nil, err
	}
	stdinPipe, err := session.StdinPipe()
	if err != nil {
		return nil, err
	}
	terminal := &Terminal{
		StdinPipe: stdinPipe,
		Session:   session,
		WsConn:    ws_conn,
	}
	session.Stdout = terminal
	session.Stderr = terminal
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err := session.RequestPty("xterm", 150, 30, modes); err != nil {
		return nil, err
	}
	if err := session.Shell(); err != nil {
		return nil, err
	}
	return terminal, nil
}
func (t *Terminal) Write(p []byte) (n int, err error) {
	writer, err := t.WsConn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return 0, err
	}
	defer writer.Close()
	return writer.Write(p)
}
func (t *Terminal) Close() error {
	if t.Session != nil {
		t.Session.Close()
	}
	return t.WsConn.Close()
}
func (t *Terminal) Read(p []byte) (n int, err error) {
	for {
		msgType, reader, err := t.WsConn.NextReader()
		if err != nil {
			return 0, err
		}
		if msgType != websocket.TextMessage {
			continue
		}
		return reader.Read(p)
	}
}
func (t *Terminal) LoopRead(logBuff *bytes.Buffer, context context.Context) error {
	for {
		select {
		case <-context.Done():
			return errors.New("LoopRead exit")
		default:
			_, wsData, err := t.WsConn.ReadMessage()
			if err != nil {
				return fmt.Errorf("reading webSocket message err:%s", err)
			}
			body, err := Decode(wsData[1:])
			if err != nil {
				logger.Error("webSocket message decode err:%s", err)
			}
			switch wsData[0] {
			case MsgResize:
				var args Resize
				if err := json.Unmarshal(body, &args); err != nil {
					return fmt.Errorf("ssh pty resize windows err:%s", err)
				}
				if args.Columns > 0 && args.Rows > 0 {
					if err := t.Session.WindowChange(args.Rows, args.Columns); err != nil {
						return fmt.Errorf("ssh pty resize windows err:%s", err)
					}
				}
			case MsgData:
				if _, err := t.StdinPipe.Write(body); err != nil {
					return fmt.Errorf("StdinPipe write err:%s", err)
				}
				if _, err := logBuff.Write(body); err != nil {
					return fmt.Errorf("logBuff write err:%s", err)
				}
			}
		}
	}
}
func (t *Terminal) SessionWait() error {
	if err := t.Session.Wait(); err != nil {
		return err
	}
	return nil
}
func Decode(p []byte) ([]byte, error) {
	decodeString, err := base64.StdEncoding.DecodeString(string(p))
	if err != nil {
		return decodeString, err
	}
	return decodeString, nil
}
func Encode(p []byte) []byte {
	encodeToString := base64.StdEncoding.EncodeToString(p)
	return []byte(encodeToString)
}
