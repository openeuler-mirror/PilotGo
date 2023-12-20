/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: guozhengxin
 * Date: 2022-05-26 10:25:52
 * LastEditTime: 2023-06-28 15:59:45
 * Description: plugin info record
 ******************************************************************************/
package dao

import "gitee.com/openeuler/PilotGo/dbmanager/mysqlmanager"

type PluginModel struct {
	ID          int    `gorm:"type:int"`
	UUID        string `gorm:"type:varchar(50)" json:"uuid"`
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

// 查询单个插件信息
func QueryPlugin(name string) (*PluginModel, error) {
	var plugins *PluginModel
	if err := mysqlmanager.MySQL().Where("name=?", name).Find(&plugins).Error; err != nil {
		return nil, err
	}
	return plugins, nil
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

// 查询插件
func GetURLAndName(uuid string) (string, string, error) {
	var plugins *PluginModel
	if err := mysqlmanager.MySQL().Where("uuid = ?", uuid).Find(&plugins).Error; err != nil {
		return "", "", err
	}
	return plugins.Url, plugins.Name, nil
}

// 查询单个插件信息
func QueryPluginById(uuid string) (*PluginModel, error) {
	var plugins *PluginModel
	if err := mysqlmanager.MySQL().Where("uuid=?", uuid).Find(&plugins).Error; err != nil {
		return nil, err
	}
	return plugins, nil
}
