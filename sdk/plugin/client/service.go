/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Sep 27 17:35:12 2023 +0800
 */
package client

import (
	"encoding/json"
	"fmt"

	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/plugin/jwt"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
)

func (c *Client) ServiceStatus(batch *common.Batch, servicename string) ([]*common.ServiceResult, error) {
	serverInfo, err := c.Registry.Get("pilotgo-server")
	if err != nil {
		return []*common.ServiceResult{}, err
	}
	url := fmt.Sprintf("http://%s:%s/api/v1/pluginapi/service/:name", serverInfo.Address, serverInfo.Port)

	p := &common.ServiceStruct{
		Batch:       batch,
		ServiceName: servicename,
	}

	r, err := httputils.Put(url, &httputils.Params{
		Body: p,
		Cookie: map[string]string{
			jwt.TokenCookie: c.token,
		},
	})
	if err != nil {
		return nil, err
	}

	res := &common.Result{}
	if err := json.Unmarshal(r.Body, res); err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (c *Client) StartService(batch *common.Batch, serviceName string) ([]*common.ServiceResult, error) {
	serverInfo, err := c.Registry.Get("pilotgo-server")
	if err != nil {
		return []*common.ServiceResult{}, err
	}
	url := fmt.Sprintf("http://%s:%s/api/v1/pluginapi/start_service", serverInfo.Address, serverInfo.Port)

	p := &common.ServiceStruct{
		Batch:       batch,
		ServiceName: serviceName,
	}

	r, err := httputils.Put(url, &httputils.Params{
		Body: p,
		Cookie: map[string]string{
			jwt.TokenCookie: c.token,
		},
	})
	if err != nil {
		return nil, err
	}

	res := &common.Result{}
	if err := json.Unmarshal(r.Body, res); err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (c *Client) StopService(batch *common.Batch, serviceName string) ([]*common.ServiceResult, error) {
	serverInfo, err := c.Registry.Get("pilotgo-server")
	if err != nil {
		return []*common.ServiceResult{}, err
	}
	url := fmt.Sprintf("http://%s:%s/api/v1/pluginapi/stop_service", serverInfo.Address, serverInfo.Port)

	p := &common.ServiceStruct{
		Batch:       batch,
		ServiceName: serviceName,
	}

	r, err := httputils.Put(url, &httputils.Params{
		Body: p,
		Cookie: map[string]string{
			jwt.TokenCookie: c.token,
		},
	})
	if err != nil {
		return nil, err
	}

	res := &common.Result{}
	if err := json.Unmarshal(r.Body, res); err != nil {
		return nil, err
	}

	return res.Data, nil
}
