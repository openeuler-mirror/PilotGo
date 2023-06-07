package client

import (
	"encoding/json"
	"io"
	"net/http"
)

type CmdResult struct {
	MachineUUID string
	MachineIP   string
	Code        int
	Stdout      string
	Stderr      string
}

func (c *Client) RunScript(batch []string, cmd string) ([]*CmdResult, error) {
	url := c.Server + "/api/v1/pluginapi/run_script"
	req, err := http.NewRequest("POST", url, nil)
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

	res := []*CmdResult{}
	if err := json.Unmarshal(bs, &res); err != nil {
		return nil, err
	}

	return res, nil
}
