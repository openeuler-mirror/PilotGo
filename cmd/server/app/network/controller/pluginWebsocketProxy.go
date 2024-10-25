package controller

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gitee.com/openeuler/PilotGo/cmd/server/app/network/jwt"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/auditlog"
	"gitee.com/openeuler/PilotGo/cmd/server/app/service/plugin"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
	"github.com/gin-gonic/gin"
	uuidservice "github.com/google/uuid"
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
	errChan := make(chan error, 1)
	defer close(errChan)

	name := c.Param("plugin_name")
	p, err := plugin.GetPlugin(name)
	if err != nil {
		c.String(http.StatusNotFound, "plugin not found: "+err.Error())
		return
	}

	if err := clientAuthentication(c); err != nil {
		c.String(http.StatusUnauthorized, err.Error())
		return
	}

	target_addr := strings.Replace(strings.Split(p.Url, "//")[1], "/plugin/"+name, "", 1)
	targetURL_str := fmt.Sprintf("ws://%s/ws/proxy", target_addr)
	ishttp, err := httputils.ServerIsHttp("http://" + target_addr)
	if err != nil {
		c.String(http.StatusNotFound, "parse plugin url error: "+err.Error())
		return
	}
	if ishttp && strings.Split(targetURL_str, "://")[0] == "wss" {
		targetURL_str = "ws://" + strings.Split(targetURL_str, "://")[1]
	}
	if !ishttp && strings.Split(targetURL_str, "://")[0] == "ws" {
		targetURL_str = "wss://" + strings.Split(targetURL_str, "://")[1]
	}

	logger.Debug("proxy plugin request to: %s", targetURL_str)

	client_wsconn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.String(http.StatusBadGateway, fmt.Sprintf("upgrade client WebSocket failed: %s", err.Error()))
		return
	}
	defer client_wsconn.Close()

	header, err := targetDirector(c.Request)
	if err != nil {
		if err := client_wsconn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, err.Error())); err != nil {
			logger.Error("websocket writemessage close: %s", err.Error())
		}
	}

	target_wsconn, _, err := dialer.Dial(targetURL_str, header)
	if err != nil {
		if err := client_wsconn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, fmt.Sprintf("dial to target WebSocket failed: %s", err.Error()))); err != nil {
			logger.Error("websocket writemessage close: %s", err.Error())
		}
		return
	}
	defer target_wsconn.Close()

	go transferMessages(client_wsconn, target_wsconn, errChan)
	go transferMessages(target_wsconn, client_wsconn, errChan)

	err = <-errChan
	logger.Error(err.Error())
	wserr := err.(*WebsocketError)
	switch wserr.Code {
	case WebsocketProxyReadError:
		if err := wserr.DstConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, wserr.Text)); err != nil {
			logger.Error("websocket writemessage close: %s", err.Error())
		}
	case WebsocketProxyWriteError:
		if err := wserr.SrcConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, wserr.Text)); err != nil {
			logger.Error("websocket writemessage close: %s", err.Error())
		}
	case WebsocketProxySingleError:
		if err := wserr.SingleConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, wserr.Text)); err != nil {
			logger.Error("websocket writemessage close: %s", err.Error())
		}
	}
}

func targetDirector(_r *http.Request) (http.Header, error) {
	cookie, err := _r.Cookie("Admin-Token")
	if err != nil {
		return nil, fmt.Errorf("fail to get client cookie: %s", err.Error())

	}
	header := http.Header{}
	header.Set("token", cookie.Value)
	return header, nil
}

func transferMessages(_srcConn, _dstConn *websocket.Conn, _err_ch chan error) {
	for {
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

func clientAuthentication(_ctx *gin.Context) error {
	u, err := jwt.ParseUser(_ctx)
	if err != nil {
		return fmt.Errorf("user token error: %s", err.Error())
	}
	log := &auditlog.AuditLog{
		LogUUID:    uuidservice.New().String(),
		ParentUUID: "",
		Module:     auditlog.ModulePlugin,
		Status:     auditlog.StatusOK,
		UserID:     u.ID,
		Action:     "parse Plugin",
	}
	auditlog.Add(log)
	return nil
}
