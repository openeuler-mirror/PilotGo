/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package dao

import (
	"fmt"
	"math/rand"
	"time"

	"gitee.com/openeuler/PilotGo/pkg/dbmanager/mysqlmanager"
	"gitee.com/openeuler/PilotGo/pkg/global"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/pkg/errors"
)

type Script struct {
	ID             uint   `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Name           string `json:"name"`
	Content        string `json:"content"`
	Description    string `json:"description"`
	UpdatedAt      time.Time
	HistoryVersion []HistoryVersion `gorm:"foreignKey:ScriptID" json:"history_version"`
	Deleted        int              `json:"deleted"` //deleted为1的时候表示删除，一般表示为0
}

type HistoryVersion struct {
	ID          uint   `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	ScriptID    uint   `json:"scriptid"`
	Version     string `gorm:"unique" json:"version"`
	Content     string `json:"content"`
	Description string `json:"description"`
	UpdatedAt   time.Time
	Script      Script //`gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// 添加脚本文件
func AddScript(_script Script) error {
	_script.UpdatedAt = time.Now()
	if err := mysqlmanager.MySQL().Save(&_script).Error; err != nil {
		return err
	}
	return nil
}

func UpdateScript(_script Script) error {
	now := time.Now()

	old_script := Script{}
	if err := mysqlmanager.MySQL().Where("id=?", _script.ID).Find(&old_script).Error; err != nil {
		return nil
	}

	err := mysqlmanager.MySQL().Model(&Script{}).Where("id=?", _script.ID).Updates(Script{
		Name:        _script.Name,
		Content:     _script.Content,
		Description: _script.Description,
		UpdatedAt:   now,
	}).Error
	if err != nil {
		return err
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	version := now.Format("2006-01-02 15:04:05") + "-" + vcode
	history := HistoryVersion{
		ScriptID:    _script.ID,
		Version:     version,
		Content:     old_script.Content,
		Description: old_script.Description,
		UpdatedAt:   now,
	}
	if err := mysqlmanager.MySQL().Save(&history).Error; err != nil {
		return err
	}
	return nil
}

// 根据脚本版本号查询文件是否存在
func IsVersionExist(scriptversion string) (bool, error) {
	var historyscript HistoryVersion
	err := mysqlmanager.MySQL().Where("version=?", scriptversion).Find(&historyscript).Error
	return len(historyscript.Version) > 0, err
}

func DeleteScript(id uint, version string) error {
	if id >= 1 && version == "" {
		if err := mysqlmanager.MySQL().Model(&HistoryVersion{}).Where("script_id = ?", id).Unscoped().Delete(HistoryVersion{}).Error; err != nil {
			return err
		}
		if err := mysqlmanager.MySQL().Model(&Script{}).Where("id=?", id).Unscoped().Delete(Script{}).Error; err != nil {
			return err
		}
		return nil
	}

	VersionExistBool, err := IsVersionExist(version)
	if err != nil {
		return err
	}
	if !VersionExistBool {
		return errors.Errorf("version %s not exist", version)
	}
	if err := mysqlmanager.MySQL().Model(&HistoryVersion{}).Where("version=?", version).Delete(HistoryVersion{}).Error; err != nil {
		return err
	}
	return nil
}

func ShowScriptContent(id uint) (string, error) {
	var script Script
	if err := mysqlmanager.MySQL().Where("id=?", id).Find(&script).Error; err != nil {
		return "", err
	}
	return script.Content, nil
}

func ShowScript(id uint) (*Script, error) {
	var script *Script
	if err := mysqlmanager.MySQL().Where("id=?", id).Find(&script).Error; err != nil {
		return nil, err
	}
	return script, nil
}

func ShowScriptWithVersion(id uint, version string) (string, error) {
	var historyscript HistoryVersion
	if err := mysqlmanager.MySQL().Where("script_id=? and version=?", id, version).Find(&historyscript).Error; err != nil {
		return "", err
	}
	return historyscript.Content, nil
}

func ScriptList(query *response.PaginationQ) ([]*Script, int, error) {
	scripts := make([]*Script, 0)
	if err := mysqlmanager.MySQL().Order("id desc").Limit(query.PageSize).Offset((query.Page - 1) * query.PageSize).Find(&scripts).Error; err != nil {
		return nil, 0, err
	}

	for _, script := range scripts {
		script.UpdatedAt = script.UpdatedAt.Local()
	}

	var total int64
	if err := mysqlmanager.MySQL().Model(&Script{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return scripts, int(total), nil
}

func GetScriptHistoryVersion(scriptid uint) ([]*HistoryVersion, error) {
	history_scripts := make([]*HistoryVersion, 0)
	if err := mysqlmanager.MySQL().Where("script_id=?", scriptid).Find(&history_scripts).Error; err != nil {
		return nil, err
	}

	for _, script := range history_scripts {
		script.UpdatedAt = script.UpdatedAt.Local()
	}
	return history_scripts, nil
}

type DangerousCommands struct {
	ID      uint   `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Command string `gorm:"unique" json:"command"`
	Active  bool   `json:"active"`
}

func CreateDangerousCommands() error {
	var total int64
	if err := mysqlmanager.MySQL().Model(&DangerousCommands{}).Count(&total).Error; err != nil {
		return err
	}

	if total == 0 {
		for _, command := range global.DangerousCommandsList {
			command_db := DangerousCommands{
				Command: command,
				Active:  true,
			}
			if err := mysqlmanager.MySQL().Create(&command_db).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func UpdateCommandsBlackList(_whitelist []uint) error {
	for _, id := range _whitelist {
		if err := mysqlmanager.MySQL().Model(&DangerousCommands{}).Where("id=?", id).Update("active", false).Error; err != nil {
			return err
		}
	}
	return nil
}

func GetDangerousCommandsList() ([]*DangerousCommands, error) {
	commands := make([]*DangerousCommands, 0)
	if err := mysqlmanager.MySQL().Order("id").Find(&commands).Error; err != nil {
		return nil, err
	}
	return commands, nil
}
