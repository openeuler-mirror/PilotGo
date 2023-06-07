package client

import (
	"encoding/json"

	"gitee.com/openeuler/PilotGo-plugins/sdk/utils"
)

type PluginInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Email       string `json:"email"`
	Url         string `json:"url"`
	ReverseDest string `json:"reverse_dest"`
}

func (c *Client) GetPluginInfo(name string) (*PluginInfo, error) {
	url := c.Server + "/api/v1/pluginapi/plugins"
	data, err := utils.Request("GET", url)
	if err != nil {
		return nil, err
	}

	resp := &PluginInfo{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
