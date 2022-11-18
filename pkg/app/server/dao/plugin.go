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
 * LastEditTime: 2022-06-02 10:16:10
 * Description: plugin info record
 ******************************************************************************/
package dao

import (
	"openeuler.org/PilotGo/PilotGo/pkg/global"
)

type PluginModel struct {
	ID          int    `gorm:"type:int"`
	UUID        string `gorm:"type:varchar(50)"`
	Name        string `gorm:"type:varchar(100)"`
	Version     string `gorm:"type:varchar(50)"`
	Description string `gorm:"type:text"`
	Author      string `gorm:"type:varchar(50)"`
	Email       string `gorm:"type:varchar(100)"`
	Url         string `gorm:"type:varchar(200)"`
	Enabled     int    `gorm:"type:int"`
}

func (m *PluginModel) TableName() string {
	return "plugin"
}

func RecordPlugin(plugin *PluginModel) error {
	err := global.PILOTGO_DB.Create(&plugin).Error
	return err
}

// 查询所有插件信息
func QueryPlugins() ([]*PluginModel, error) {
	var plugins []*PluginModel
	if err := global.PILOTGO_DB.Find(&plugins).Error; err != nil {
		return nil, err
	}
	return plugins, nil
}

// 更新插件使能状态
func UpdatePluginEnabled(plugin *PluginModel) error {
	err := global.PILOTGO_DB.Model(&plugin).Update("Enabled", plugin.Enabled).Error
	return err
}

// 删除插件
func DeletePlugin(uuid string) error {
	err := global.PILOTGO_DB.Where("uuid=?", uuid).Delete(&PluginModel{}).Error
	return err
}
