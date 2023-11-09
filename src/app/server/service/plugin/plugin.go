package plugin

import (
	"encoding/json"
	"fmt"
	"strings"

	"gitee.com/openeuler/PilotGo/app/server/config"
	"gitee.com/openeuler/PilotGo/app/server/service/extention"
	"gitee.com/openeuler/PilotGo/app/server/service/internal/dao"
	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/plugin/client"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
	"github.com/google/uuid"
)

const (
	PluginEnabled  = 1
	PluginDisabled = 0
)

type Plugin = dao.PluginModel

// 与plugin进行握手，交换必要信息
func Handshake(url string) (*Plugin, error) {
	info, err := requestPluginInfo(url)
	if err != nil {
		logger.Debug("")
		return nil, err
	}

	plugin := &Plugin{
		Name:        info.Name,
		Version:     info.Version,
		Description: info.Description,
		Author:      info.Author,
		Email:       info.Email,
		Url:         info.Url,
		PluginType:  info.PluginType,
		// Status:      common.StatusLoaded,
	}

	return plugin, nil
}

// 发起http请求，提供server地址，同时获取到插件的基本信息
func requestPluginInfo(url string) (*client.PluginInfo, error) {
	conf := config.Config().HttpServer
	url = url + fmt.Sprintf("?server=%s", conf.Addr)
	resp, err := httputils.Get(url, nil)
	if err != nil {
		logger.Debug("request plugin info error:%s", err.Error())
		return nil, err
	}

	info := &client.PluginInfo{}
	err = json.Unmarshal(resp.Body, info)
	if err != nil {
		logger.Debug("unmarshal request plugin info error:%s", err.Error())
	}
	// TODO: check info valid
	return info, nil
}

// 获取plugin清单
func GetPlugins() ([]*Plugin, error) {
	plugins, err := dao.QueryPlugins()
	if err != nil {
		logger.Error("failed to read plugin info from db:%s", err.Error())
		return nil, err
	}

	return plugins, nil
}

// 分页查询
func GetPluginPaged(offset, size int) (int64, []Plugin, error) {
	return dao.GetPluginPaged(offset, size)
}

func GetPlugin(name string) (*Plugin, error) {
	plugin, err := dao.QueryPlugin(name)
	if err != nil {
		logger.Error("failed to read plugin info from db:%s", err.Error())
		return nil, err
	}
	return plugin, nil
}

type AddPluginParam struct {
	Name string `json:"name"`
	Type string `json:"plugin_type"`
	Url  string `json:"url"`
}

func AddPlugin(param *AddPluginParam) error {
	url := param.Url
	logger.Debug("add plugin from %s", url)
	url = strings.TrimRight(url, "/")

	plugin, err := Handshake(url + "/plugin_manage/info")
	if err != nil {
		return err
	}
	plugin.UUID = uuid.New().String()
	plugin.PluginType = param.Type

	if err := dao.RecordPlugin(plugin); err != nil {
		return err
	}
	return nil
}

func DeletePlugin(uuid string) error {
	logger.Debug("delete plugin: %s", uuid)

	if err := dao.DeletePlugin(uuid); err != nil {
		logger.Error("failed to delete plugin info:%s", err.Error())
	}
	return nil
}

func TogglePlugin(uuid string, enable int) ([]common.Extention, error) {
	logger.Debug("toggle plugin: %s to enable %d ", uuid, enable)
	if err := dao.UpdatePluginEnabled(&dao.PluginModel{
		UUID:    uuid,
		Enabled: enable,
	}); err != nil {
		return nil, err
	}
	plugin, err := dao.QueryPluginById(uuid)
	if err != nil {
		return nil, err
	}
	return extention.RequestExtention(*plugin)
}
