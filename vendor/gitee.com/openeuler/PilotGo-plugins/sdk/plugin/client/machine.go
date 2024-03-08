package client

import (
	"encoding/json"

	"gitee.com/openeuler/PilotGo-plugins/sdk/common"
	"gitee.com/openeuler/PilotGo-plugins/sdk/utils/httputils"
)

func (c *Client) MachineList() ([]*common.MachineNode, error) {
	url := c.Server + "/api/v1/pluginapi/machine_list"
	r, err := httputils.Get(url, nil)
	if err != nil {
		return nil, err
	}

	result := []*common.MachineNode{}
	if err := json.Unmarshal(r.Body, &result); err != nil {
		return nil, err
	}
	return result, nil
}
