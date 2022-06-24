/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: zhanghan
 * Date: 2022-04-20 16:48:55
 * LastEditTime: 2022-04-20 17:48:55
 * Description: web socket连接结构体
 ******************************************************************************/
package model

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
	"unicode/utf8"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

type PtyRequestMsg struct {
	Term     string
	Columns  uint32
	Rows     uint32
	Width    uint32
	Height   uint32
	Modelist string
}

type Terminal struct {
	Columns uint32 `json:"cols"`
	Rows    uint32 `json:"rows"`
}

type SSHClient struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	IpAddress string `json:"ipaddress"`
	Port      int    `json:"port"`
	Session   *ssh.Session
	Client    *ssh.Client
	Channel   ssh.Channel
}

type ApiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *ApiError) Error() string {
	return e.Message
}

func (sshc *SSHClient) GenerateClient() error {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		config       ssh.Config
		err          error
	)
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(sshc.Password))
	config = ssh.Config{
		Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
	}
	// ssh客户端
	clientConfig = &ssh.ClientConfig{
		User:    sshc.Username,
		Auth:    auth,
		Timeout: 5 * time.Second,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	addr = fmt.Sprintf("%s:%d", sshc.IpAddress, sshc.Port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return err
	}
	sshc.Client = client
	return nil
}

func (sshc *SSHClient) RequestTerminal(terminal Terminal) *SSHClient {
	// 通过ssh客户端创建ssh channel,并请求一个pty伪终端,请求用户的默认会话
	session, err := sshc.Client.NewSession()
	if err != nil {
		log.Println(err)
		return nil
	}
	sshc.Session = session
	// 如果主机信息验证通过, 可以通过ssh client创建一个通道:
	channel, inRequests, err := sshc.Client.OpenChannel("session", nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	sshc.Channel = channel
	go func() {
		for req := range inRequests {
			if req.WantReply {
				req.Reply(false, nil)
			}
		}
	}()
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	var modeList []byte
	for k, v := range modes {
		kv := struct {
			Key byte
			Val uint32
		}{k, v}
		modeList = append(modeList, ssh.Marshal(&kv)...)
	}
	modeList = append(modeList, 0)
	req := PtyRequestMsg{
		Term:     "xterm",
		Columns:  terminal.Columns,
		Rows:     terminal.Rows,
		Width:    uint32(terminal.Columns * 8),
		Height:   uint32(terminal.Columns * 8),
		Modelist: string(modeList),
	}
	// 如果主机信息验证通过, 可以通过ssh client创建一个通道:
	ok, err := channel.SendRequest("pty-req", true, ssh.Marshal(&req))
	if !ok || err != nil {
		log.Println(err)
		return nil
	}
	ok, err = channel.SendRequest("shell", true, nil)
	if !ok || err != nil {
		log.Println(err)
		return nil
	}
	return sshc
}

func (sshc *SSHClient) Connect(ws *websocket.Conn) {

	//这里第一个协程获取用户的输入
	go func() {
		for {
			// p为用户输入
			_, p, err := ws.ReadMessage()
			if err != nil {
				return
			}
			_, err = sshc.Channel.Write(p)
			if err != nil {
				return
			}
		}
	}()

	//第二个协程将远程主机的返回结果返回给用户
	go func() {
		br := bufio.NewReader(sshc.Channel)
		buf := []byte{}
		t := time.NewTimer(time.Microsecond * 100)
		defer t.Stop()
		// 构建一个信道, 一端将数据远程主机的数据写入, 一段读取数据写入ws
		r := make(chan rune)

		// 另起一个协程, 一个死循环不断的读取ssh channel的数据, 并传给r信道直到连接断开
		go func() {
			defer sshc.Client.Close()
			defer sshc.Session.Close()

			for {
				x, size, err := br.ReadRune()
				if err != nil {
					log.Println(err)
					ws.WriteMessage(1, []byte("\033[31m已经关闭连接!\033[0m"))
					ws.Close()
					return
				}
				if size > 0 {
					r <- x
				}
			}
		}()

		// 主循环
		for {
			select {
			// 每隔100微秒, 只要buf的长度不为0就将数据写入ws, 并重置时间和buf
			case <-t.C:
				if len(buf) != 0 {
					err := ws.WriteMessage(websocket.TextMessage, buf)
					buf = []byte{}
					if err != nil {
						log.Println(err)
						return
					}
				}
				t.Reset(time.Microsecond * 100)
			// 前面已经将ssh channel里读取的数据写入创建的通道r, 这里读取数据, 不断增加buf的长度, 在设定的 100 microsecond后由上面判定长度是否返送数据
			case d := <-r:
				if d != utf8.RuneError {
					p := make([]byte, utf8.RuneLen(d))
					utf8.EncodeRune(p, d)
					buf = append(buf, p...)
				} else {
					buf = append(buf, []byte("@")...)
				}
			}
		}
	}()

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
}
