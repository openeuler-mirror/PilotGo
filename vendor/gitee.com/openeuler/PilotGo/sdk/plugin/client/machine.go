package client

import (
	"encoding/json"
	"errors"

	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
)

func (c *Client) MachineList() ([]*common.MachineNode, error) {
	if !c.IsBind() {
		return nil, errors.New("unbind PilotGo-server platform")
	}
	url := "http://" + c.Server() + "/api/v1/pluginapi/machine_list"
	r, err := httputils.Get(url, &httputils.Params{
		Cookie: map[string]string{
			TokenCookie: c.token,
		},
	})
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

func (c *Client) MachineInfoByUUID(machine_uuid string) (*common.MachineNode, error) {
	if !c.IsBind() {
		return nil, errors.New("unbind PilotGo-server platform")
	}
	url := "http://" + c.Server() + "/api/v1/pluginapi/machine_info?machine_uuid=" + machine_uuid
	r, err := httputils.Get(url, &httputils.Params{
		Cookie: map[string]string{
			TokenCookie: c.token,
		},
	})
	if err != nil {
		return nil, err
	}

	result := struct {
		Code int `json:"code"`
		Data *struct {
			ID          int    `json:"id"`
			Departid    int    `json:"departid"`
			Departname  string `json:"departname"`
			IP          string `json:"ip"`
			UUID        string `json:"uuid"`
			CPU         string `json:"cpu"`
			Runstatus   string `json:"runstatus"`
			Maintstatus string `json:"maintstatus"`
			Systeminfo  string `json:"systeminfo"`
		} `json:"data"`
	}{}
	if err := json.Unmarshal(r.Body, &result); err != nil {
		return nil, err
	}

	res := &common.MachineNode{
		UUID:        machine_uuid,
		Department:  result.Data.Departname,
		IP:          result.Data.IP,
		CPUArch:     result.Data.CPU,
		OS:          result.Data.Systeminfo,
		RunStatus:   result.Data.Runstatus,
		MaintStatus: result.Data.Maintstatus,
	}
	return res, nil
}
