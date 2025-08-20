/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Wangjunqi123 <wangjunqi@kylinos.cn>
 * Date: Fri Oct 25 16:02:57 2024 +0800
 */
package controller

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"gitee.com/openeuler/PilotGo/pkg/global"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: 10 * time.Second,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	dialer = websocket.Dialer{
		Proxy: nil,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		HandshakeTimeout: 10 * time.Second,
	}
)

type WebsocketError struct {
	Code       int
	SrcConn    *websocket.Conn
	DstConn    *websocket.Conn
	SingleConn *websocket.Conn
	Text       string
}

const (
	WebsocketProxyReadError int = iota
	WebsocketProxyWriteError
	WebsocketProxySingleError
)

func (we *WebsocketError) Error() string {
	str := ""
	switch we.Code {
	case WebsocketProxyReadError:
		str = fmt.Sprintf("websocket proxy read error: %s", we.Text)
	case WebsocketProxyWriteError:
		str = fmt.Sprintf("websocket proxy write error: %s", we.Text)
	case WebsocketProxySingleError:
		str = fmt.Sprintf("websocket proxy error: %s", we.Text)
	}
	return str
}

func PluginWebsocketGatewayHandler(c *gin.Context) {
	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	doneChan := make(chan struct{})
	defer func() {
		close(doneChan)
		wg.Wait()
		close(errChan)
	}()

	name := c.Param("serviceName")
	p := global.GW.GetService(name)

	targetURL_str := fmt.Sprintf("ws://%s/ws/proxy", p.Address+":"+p.Port)

	logger.Debug("websocket proxy plugin request: %s->%s", c.Request.RemoteAddr, targetURL_str)

	client_wsconn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.String(http.StatusBadGateway, fmt.Sprintf("upgrade client WebSocket failed: %s", err.Error()))
		return
	}
	defer client_wsconn.Close()

	header, err := targetDirector(c)
	if err != nil {
		if err := client_wsconn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error())); err != nil {
			logger.Error("websocket writemessage close: %s", err.Error())
		}
	}

	target_wsconn, _, err := dialer.Dial(targetURL_str, header)
	if err != nil {
		if err := client_wsconn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, fmt.Sprintf("dial to plugin server WebSocket failed: %s", err.Error()))); err != nil {
			logger.Error("websocket writemessage close: %s", err.Error())
		}
		return
	}
	defer target_wsconn.Close()

	wg.Add(1)
	go transferMessages(client_wsconn, target_wsconn, errChan, &wg, doneChan)
	wg.Add(1)
	go transferMessages(target_wsconn, client_wsconn, errChan, &wg, doneChan)

	err = <-errChan
	logger.Error("%s", err.Error())
	wserr := err.(*WebsocketError)
	switch wserr.Code {
	case WebsocketProxyReadError:
		if err := wserr.DstConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, wserr.Text)); err != nil {
			logger.Error("write close message error: %s(wserr.text: %s)", err.Error(), wserr.Text)
		}
	case WebsocketProxyWriteError:
		if err := wserr.SrcConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, wserr.Text)); err != nil {
			logger.Error("write close message error: %s(wserr.text: %s)", err.Error(), wserr.Text)
		}
	case WebsocketProxySingleError:
		if err := wserr.SingleConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, wserr.Text)); err != nil {
			logger.Error("write close message error: %s(wserr.text: %s)", err.Error(), wserr.Text)
		}
	}
}

func targetDirector(_ctx *gin.Context) (http.Header, error) {
	header := http.Header{}

	header.Set("clientId", _ctx.Query("clientId"))

	if clientIP, clientPort, err := net.SplitHostPort(_ctx.Request.RemoteAddr); err == nil {
		if prior, ok := _ctx.Request.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP + ":" + clientPort
		}
		header.Set("X-Forwarded-For", clientIP+":"+clientPort)
	}

	header.Set("X-Forwarded-Proto", "http")
	if _ctx.Request.TLS != nil {
		header.Set("X-Forwarded-Proto", "https")
	}
	return header, nil
}

func transferMessages(_srcConn, _dstConn *websocket.Conn, _err_ch chan error, _wg *sync.WaitGroup, _donechan chan struct{}) {
	defer _wg.Done()
	for {
		select {
		case <-_donechan:
			return
		default:
			messageType, message, err := _srcConn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
					_err_ch <- &WebsocketError{
						Code:    WebsocketProxyReadError,
						SrcConn: _srcConn,
						DstConn: _dstConn,
						Text:    fmt.Sprintf("websocket src conn %s closed(%v->%v): %s", _srcConn.RemoteAddr().String(), _srcConn.RemoteAddr().String(), _dstConn.RemoteAddr().String(), err.Error()),
					}
					return
				}
				_err_ch <- &WebsocketError{
					Code:    WebsocketProxyReadError,
					SrcConn: _srcConn,
					DstConn: _dstConn,
					Text:    fmt.Sprintf("error while reading message(%v->%v, msgType: %d): %s, %s", _srcConn.RemoteAddr().String(), _dstConn.RemoteAddr().String(), messageType, err.Error(), message),
				}
				return
			}

			if err := _dstConn.WriteMessage(messageType, message); err != nil {
				_err_ch <- &WebsocketError{
					Code:    WebsocketProxyWriteError,
					SrcConn: _srcConn,
					DstConn: _dstConn,
					Text:    fmt.Sprintf("error while writing message(%v->%v): %s", _srcConn.RemoteAddr().String(), _dstConn.RemoteAddr().String(), err.Error()),
				}
				return
			}
		}
	}
}
