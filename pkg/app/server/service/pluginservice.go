package service

import (
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/plugin"
)

func PluginLists() []model.PluginLists {
	m := plugin.GetManager()
	plugins, _ := m.GetAll()
	res := []model.PluginLists{}
	status := 0
	url := ""
	for _, name := range plugins {
		if m.IsPluginLoaded(name) {
			status = 1
		}
		url, _ = m.GetWebExtension(name)
		res = append(res, model.PluginLists{Name: name, Status: status, Url: url})
	}
	return res
}

func LoadPlugin(pluginInfo model.LoadPlugin) error {
	m := plugin.GetManager()
	return m.Load(pluginInfo.Name)
}

func UnLoadPlugin(pluginInfo model.UnLoadPlugin) error {
	m := plugin.GetManager()
	return m.UnLoad(pluginInfo.Name)
}
