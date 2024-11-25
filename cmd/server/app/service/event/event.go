package event

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"gitee.com/openeuler/PilotGo/cmd/server/app/service/plugin"
	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
)

// 发布event事件
func PublishEvent(msg common.EventMessage) error {
	if eventServer, connected := isPluginEventConnected(); connected {
		err := publishEvent(eventServer, msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func publishEvent(eventServer string, msg common.EventMessage) error {
	url := eventServer + "/plugin/event/publishEvent"
	r, err := httputils.Put(url, &httputils.Params{
		Body: &msg,
	})
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		return errors.New("server process error:" + strconv.Itoa(r.StatusCode))
	}

	resp := &common.CommonResult{}
	if err := json.Unmarshal(r.Body, resp); err != nil {
		return err
	}
	if resp.Code != http.StatusOK {
		return errors.New(resp.Message)
	}
	data := &struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}{}
	if err := resp.ParseData(data); err != nil {
		return err
	}
	return nil
}

func isPluginEventConnected() (string, bool) {
	var eventServer string
	var connected bool

	plugins, err := plugin.GetPlugins()
	if err != nil {
		return "", false
	}

	for _, p := range plugins {
		if p.Name == "event" {
			eventServer = p.Url
			connected = p.ConnectStatus
			break
		}
	}

	if eventServer == "" || connected == false {
		return "", false
	}
	return eventServer, connected
}
