package client

import (
	"encoding/json"
	"errors"

	"gitee.com/openeuler/PilotGo-plugins/sdk/common"
	"gitee.com/openeuler/PilotGo-plugins/sdk/utils/httputils"
)

func (c *Client) ApplyConfig(batch *common.Batch, path, content string) error {
	url := c.Server + "/api/v1/pluginapi/apply_config"
	data, err := httputils.Put(url, nil)
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
