package client

import (
	"errors"
	"strings"

	"gitee.com/openeuler/PilotGo/sdk/common"
)

func (c *Client) RegisterPermission(pers []common.Permission) error {
	for _, v := range pers {
		if strings.Contains(v.Resource, "/") {
			return errors.New("permission-resource string contains /")
		} else {
			c.permissions = append(c.permissions, v)
		}
	}
	return nil
}
