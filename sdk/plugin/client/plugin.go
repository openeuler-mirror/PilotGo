package client

import (
	"encoding/json"
	"errors"

	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
)

// 用于初始化Client
type PluginInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Email       string `json:"email"`
	Url         string `json:"url"`
	PluginType  string `json:"plugin_type"`
	ReverseDest string `json:"reverse_dest"`
}

// 用于插件与PilotGo server通讯
type PluginFullInfo struct {
	PluginInfo
	Extentions  []common.Extention
	Permissions []common.Permission
}

func (c *Client) GetPluginInfo(name string) (*PluginInfo, error) {
	if !c.IsBind() {
		return nil, errors.New("unbind PilotGo-server platform")
	}
	url := c.Server() + "/api/v1/pluginapi/plugins"
	r, err := httputils.Get(url, nil)
	if err != nil {
		return nil, err
	}

	resp := &PluginInfo{}
	if err := json.Unmarshal(r.Body, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
