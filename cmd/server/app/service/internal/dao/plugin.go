/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package dao

import "gitee.com/openeuler/PilotGo/pkg/dbmanager/mysqlmanager"

type PluginModel struct {
	ID          int    `gorm:"type:int"`
	UUID        string `gorm:"type:varchar(50)" json:"uuid"`
	CustomName  string `gorm:"type:varchar(100)" json:"custom_name"`
	Name        string `gorm:"type:varchar(100)" json:"name"`
	Version     string `gorm:"type:varchar(50)" json:"version"`
	Description string `gorm:"type:text" json:"description"`
	Author      string `gorm:"type:varchar(50)" json:"author"`
	Email       string `gorm:"type:varchar(100)" json:"email"`
	Url         string `gorm:"type:varchar(200)" json:"url"`
	// 插件类型，iframe/microapp
	PluginType string `gorm:"type:varchar(50)" json:"plugin_type"`
	Enabled    int    `gorm:"type:int" json:"enabled"`
}

func (m *PluginModel) TableName() string {
	return "plugin"
}

func RecordPlugin(plugin *PluginModel) error {
	err := mysqlmanager.MySQL().Create(&plugin).Error
	return err
}

// 查询所有插件信息
func QueryPlugins() ([]*PluginModel, error) {
	var plugins []*PluginModel
	if err := mysqlmanager.MySQL().Find(&plugins).Error; err != nil {
		return nil, err
	}
	return plugins, nil
}

// 分页查询
func GetPluginPaged(offset, size int) (int64, []*PluginModel, error) {
	var count int64
	var pluginModels []*PluginModel
	err := mysqlmanager.MySQL().Model(PluginModel{}).Order("id desc").Offset(offset).Limit(size).Find(&pluginModels).Offset(-1).Limit(-1).Count(&count).Error
	return count, pluginModels, err
}

// 更新插件使能状态
func UpdatePluginEnabled(plugin *PluginModel) error {
	var p PluginModel
	err := mysqlmanager.MySQL().Model(&p).Where("uuid = ?", plugin.UUID).Update("enabled", plugin.Enabled).Error
	return err
}

// 更新插件信息
func UpdatePluginInfo(plugin *PluginModel) error {
	var p PluginModel
	err := mysqlmanager.MySQL().Model(&p).Where("uuid = ?", plugin.UUID).Updates(plugin).Error
	return err
}

// 删除插件
func DeletePlugin(uuid string) error {
	err := mysqlmanager.MySQL().Where("uuid=?", uuid).Delete(&PluginModel{}).Error
	return err
}

// 查询单个插件信息
func QueryPluginById(uuid string) (*PluginModel, error) {
	var plugins *PluginModel
	if err := mysqlmanager.MySQL().Where("uuid=?", uuid).Find(&plugins).Error; err != nil {
		return nil, err
	}
	return plugins, nil
}

func IsExistCustomName(name string) (bool, error) {
	var pm PluginModel
	err := mysqlmanager.MySQL().Where("custom_name = ?", name).Find(&pm).Error
	return pm.ID != 0, err
}
