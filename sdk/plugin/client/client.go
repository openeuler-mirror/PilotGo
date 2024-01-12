package client

import (
	"sync"

	"gitee.com/openeuler/PilotGo/sdk/common"
	"github.com/gin-gonic/gin"
)

type GetTagsCallback func([]string) []common.Tag

type Client struct {
	PluginInfo *PluginInfo
	// 并发锁
	l sync.Mutex

	// 远程PilotGo server地址
	server string

	// 用于event消息处理
	eventChan        chan *common.EventMessage
	eventCallbackMap map[int]EventCallback

	// 用于异步command及script执行结果处理机
	asyncCmdResultChan      chan *common.AsyncCmdResult
	cmdProcessorCallbackMap map[string]CallbackHandler

	// 用于处理主机标签
	getTagsCallback GetTagsCallback

	// 用于平台扩展点功能
	extentions []common.Extention

	//用于权限校验
	permissions []common.Permission

	//bind阻塞功能支持
	mu   sync.Mutex
	cond *sync.Cond
}

var global_client *Client

func DefaultClient(desc *PluginInfo) *Client {
	global_client = &Client{
		PluginInfo: desc,

		eventChan:        make(chan *common.EventMessage, 20),
		eventCallbackMap: make(map[int]EventCallback),

		asyncCmdResultChan:      make(chan *common.AsyncCmdResult, 20),
		cmdProcessorCallbackMap: make(map[string]CallbackHandler),
		extentions:              []common.Extention{},
		permissions:             []common.Permission{},
	}
	global_client.cond = sync.NewCond(&global_client.mu)

	return global_client
}

func GetClient() *Client {
	return global_client
}

func (client *Client) Server() string {
	return client.server
}

// RegisterHandlers 注册一些插件标准的API接口，清单如下：
// GET /plugin_manage/info
func (client *Client) RegisterHandlers(router *gin.Engine) {
	// 提供插件基本信息
	mg := router.Group("/plugin_manage/", func(c *gin.Context) {
		c.Set("__internal__client_instance", client)
	})
	{
		mg.GET("/info", infoHandler)
		// 绑定PilotGo server
		mg.PUT("/bind", bindHandler)
	}

	api := router.Group("/plugin_manage/api/v1/")
	{
		api.GET("/gettags", func(c *gin.Context) {
			c.Set("__internal__client_instance", client)
		}, tagsHandler)

		api.POST("/event", func(c *gin.Context) {
			c.Set("__internal__client_instance", client)
		}, eventHandler)

		api.PUT("/command_result", func(c *gin.Context) {
			c.Set("__internal__client_instance", client)
		}, commandResultHandler)
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

func (client *Client) OnGetTags(callback GetTagsCallback) {
	client.getTagsCallback = callback
}

// client是否bind PilotGo server
func (client *Client) IsBind() bool {
	client.l.Lock()
	defer client.l.Unlock()

	return !(client.server == "")
}
