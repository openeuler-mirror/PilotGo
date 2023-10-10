package client

import (
	"gitee.com/openeuler/PilotGo/sdk/common"
	"github.com/gin-gonic/gin"
)

type Client struct {
	Server     string
	PluginInfo *PluginInfo

	// 用于event消息处理
	eventChan chan *common.EventMessage
}

var global_client *Client
var BaseInfo *PluginInfo

func DefaultClient(desc *PluginInfo) *Client {
	BaseInfo = desc

	global_client = &Client{
		PluginInfo: desc,
	}

	return global_client
}

func GetClient() *Client {
	return global_client
}

// RegisterHandlers 注册一些插件标准的API接口，清单如下：
// GET /plugin_manage/info
func (c *Client) RegisterHandlers(router *gin.Engine) {
	// 提供插件基本信息
	mg := router.Group("/plugin_manage/")
	{
		mg.GET("/info", InfoHandler)
	}

	api := router.Group("/plugin_manage/api/v1/")
	{
		api.PUT("/event", func(c *gin.Context) {
			c.Set("__internal__client_instance", c)
		}, EventHandler)
	}

	// pg := router.Group("/plugin/" + desc.Name)
	// {
	// 	pg.Any("/*any", func(c *gin.Context) {
	// 		c.Set("__internal__reverse_dest", dest)
	// 		ReverseProxyHandler(c)
	// 	})
	// }
}
