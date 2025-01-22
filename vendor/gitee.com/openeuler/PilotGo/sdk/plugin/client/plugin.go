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
	"errors"

	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
)

// 用于初始化Client
type PluginInfo struct {
	MenuName    string `json:"menuName"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Email       string `json:"email"`
	Url         string `json:"url"`
	Icon        string `json:"icon"`
	PluginType  string `json:"plugin_type"`
	ReverseDest string `json:"reverse_dest"`
}

// 用于插件与PilotGo server通讯
type PluginFullInfo struct {
	PluginInfo
	Extentions  []common.Extention
	Permissions []common.Permission
}

func (c *Client) GetPlugins() ([]*PluginInfo, error) {
	if !c.IsBind() {
		return nil, errors.New("unbind PilotGo-server platform")
	}
	url := "http://" + c.Server() + "/api/v1/pluginapi/plugins"
	r, err := httputils.Get(url, &httputils.Params{
		Cookie: map[string]string{
			TokenCookie: c.token,
		},
	})
	if err != nil {
		return nil, err
	}

	resp := struct {
		Code int           `json:"code"`
		Data []*PluginInfo `json:"data"`
	}{}
	if err := json.Unmarshal(r.Body, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}
