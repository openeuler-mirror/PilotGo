package client

import (
	"encoding/json"
	"errors"

	"gitee.com/openeuler/PilotGo-plugins/sdk/utils"
)

func (c *Client) ApplyConfig(batch []string, path, content string) error {
	url := c.Server + "/api/v1/pluginapi/apply_config"
	data, err := utils.Request("PUT", url)
	if err != nil {
		return err
	}

	resp := &struct {
		Status string
		Error  string
	}{}
	if err := json.Unmarshal(data, resp); err != nil {
		return err
	}
	if resp.Status != "ok" {
		return errors.New(resp.Error)
	}

	return nil
}
