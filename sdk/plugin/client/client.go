package client

import (
	"gitee.com/openeuler/PilotGo/sdk/common"
	"github.com/gin-gonic/gin"
)

type Client struct {
	Server     string
	PluginInfo *PluginInfo

	// 用于event消息处理
	eventChan        chan *common.EventMessage
	eventCallbackMap map[int]EventCallback

	// 用于异步command及script执行结果处理机
	asyncCmdResultChan      chan *common.AsyncCmdResult
	cmdProcessorCallbackMap map[string]CallbackHandler
}

var global_client *Client
var BaseInfo *PluginInfo

func DefaultClient(desc *PluginInfo) *Client {
	BaseInfo = desc

	global_client = &Client{
		PluginInfo: desc,

		eventChan:        make(chan *common.EventMessage, 20),
		eventCallbackMap: make(map[int]EventCallback),

		asyncCmdResultChan:      make(chan *common.AsyncCmdResult, 20),
		cmdProcessorCallbackMap: make(map[string]CallbackHandler),
	}

	return global_client
}

func GetClient() *Client {
	return global_client
}

// RegisterHandlers 注册一些插件标准的API接口，清单如下：
// GET /plugin_manage/info
func (client *Client) RegisterHandlers(router *gin.Engine) {
	// 提供插件基本信息
	mg := router.Group("/plugin_manage/")
	{
		mg.GET("/info", InfoHandler)
	}

	api := router.Group("/plugin_manage/api/v1/")
	{
		api.POST("/event", func(c *gin.Context) {
			c.Set("__internal__client_instance", client)
		}, EventHandler)

		api.PUT("/command_result", func(c *gin.Context) {
			c.Set("__internal__client_instance", client)
		}, CommandResultHandler)
	}

	// pg := router.Group("/plugin/" + desc.Name)
	// {
	// 	pg.Any("/*any", func(c *gin.Context) {
	// 		c.Set("__internal__reverse_dest", dest)
	// 		ReverseProxyHandler(c)
	// 	})
	// }

	// TODO: start command result process service
	client.startEventProcessor()
	client.startCommandResultProcessor()
}
