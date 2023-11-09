package client

import (
	"encoding/json"

	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
)

func (c *Client) MachineList() ([]*common.MachineNode, error) {
	url := "http://" + c.Server + "/api/v1/pluginapi/machine_list"
	r, err := httputils.Get(url, nil)
	if err != nil {
		return nil, err
	}

	result := struct {
		Code int                   `json:"code"`
		Data []*common.MachineNode `json:"data"`
	}{}
	if err := json.Unmarshal(r.Body, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}
