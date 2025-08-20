/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan <zhanghan@kylinos.cn>
 * Date: Mon Nov 25 16:52:07 2024 +0800
 */
package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"gitee.com/openeuler/PilotGo/pkg/global"
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
	ok := global.GW.GetServiceStatus("event-service")
	if !ok {
		return "", false
	}

	service := global.GW.GetService("event-service")
	return fmt.Sprintf("http://%s:%s", service.Address, service.Port), true
}
