package client

import (
	"encoding/json"
	"errors"

	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
)

func (c *Client) ApplyConfig(batch *common.Batch, path, content string) error {
	if !c.IsBind() {
		return errors.New("unbind PilotGo-server platform")
	}
	url := c.Server() + "/api/v1/pluginapi/apply_config"
	r, err := httputils.Put(url, nil)
	if err != nil {
		return err
	}

	resp := &struct {
		Status string
		Error  string
	}{}
	if err := json.Unmarshal(r.Body, resp); err != nil {
		return err
	}
	if resp.Status != "ok" {
		return errors.New(resp.Error)
	}

	return nil
}
