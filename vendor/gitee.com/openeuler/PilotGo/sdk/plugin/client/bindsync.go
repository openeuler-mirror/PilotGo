package client

import (
	"gitee.com/openeuler/PilotGo/sdk/logger"
)

func (c *Client) Wait4Bind() {
	c.cond.L.Lock()
	logger.Debug("等待bind中...")
	for c.server == "" {
		c.cond.Wait() // 等待条件变量被通知
	}
	c.cond.L.Unlock()
	logger.Debug("bind 成功！")
}
