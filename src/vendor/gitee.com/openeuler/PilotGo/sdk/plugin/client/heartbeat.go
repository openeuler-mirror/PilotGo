package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
)

const (
	HeartbeatInterval = 30 * time.Second
	HeartbeatKey      = "heartbeat:"
)

// 插件连接状态
type PluginStatus struct {
	Connected   bool
	LastConnect time.Time
}

func (client *Client) sendHeartBeat() {
	clientID := client.PluginInfo.Url
	go func() {
		for {
			err := client.sendHeartbeat(clientID)
			if err != nil {
				logger.Error("Heartbeat failed:%v", err)
			}
			time.Sleep(HeartbeatInterval)
		}
	}()
}

func (client *Client) sendHeartbeat(clientID string) error {
	p := &struct {
		PluginUrl string `json:"clientID"`
	}{
		PluginUrl: clientID,
	}

	ServerUrl := "http://" + client.Server() + "/api/v1/pluginapi/heartbeat"
	resp, err := httputils.Post(ServerUrl, &httputils.Params{
		Body: p,
		Cookie: map[string]string{
			TokenCookie: client.token,
		},
	})
	if err != nil {
		return err
	}
	res := &struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
	}{}
	if err := json.Unmarshal(resp.Body, res); err != nil {
		return err
	}
	if res.Code != http.StatusOK {
		return fmt.Errorf("heartbeat failed with status: %v", res.Code)
	}
	return nil
}
