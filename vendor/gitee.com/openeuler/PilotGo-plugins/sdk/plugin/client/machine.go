package client

import (
	"encoding/json"
	"io"
	"net/http"
)

type MachineNode struct {
	UUID       string
	Department string
	IP         string
	CPUArch    string
	OSInfo     string
	State      int
}

func (c *Client) MachineList() ([]*MachineNode, error) {
	url := c.Server + "/api/v1/pluginapi/machine_list"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	hc := &http.Client{}
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := []*MachineNode{}
	if err := json.Unmarshal(bs, &result); err != nil {
		return nil, err
	}
	return result, nil
}
