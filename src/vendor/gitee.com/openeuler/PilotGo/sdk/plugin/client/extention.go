package client

import (
	"gitee.com/openeuler/PilotGo/sdk/common"
)

func (c *Client) RegisterExtention(exts []*common.Extention) {
	c.extentions = exts
}
