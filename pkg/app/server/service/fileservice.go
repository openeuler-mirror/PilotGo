/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: zhanghan
 * Date: 2022-05-26 10:25:52
 * LastEditTime: 2022-06-02 10:16:10
 * Description: agent config file service
 ******************************************************************************/
package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"openeuler.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
)

// 获取时间的日期函数 => 20200426-17:36:04
func NowTime() string {
	time := time.Now()
	year := time.Year()
	month := time.Month()
	day := time.Day()
	hour := time.Hour()
	minute := time.Minute()
	second := time.Second()
	nowtime := fmt.Sprintf("%d%02d%02d-%02d:%02d:%02d", year, month, day, hour, minute, second)
	return nowtime
}

// dhcp方式配置网络
func NetworkDHCP(net []interface{}) (text string) {
	for _, n := range net {
		nn := n.(map[string]interface{})
		for key, value := range nn {
			if key == "BOOTPROTO" {
				text += key + "=" + "dhcp" + "\n"
			} else if key == "IPADDR" {
				break
			} else if key == "NETMASK" {
				break
			} else if key == "GATEWAY" {
				break
			} else if key == "DNS1" {
				break
			} else if key == "DNS2" {
				break
			} else {
				text += key + "=" + value.(string) + "\n"
			}
		}
	}
	return
}

// static方式配置网络
func NetworkStatic(net []interface{}, ip string, netmask string, gateway string, dns1 string, dns2 string) (text string) {
	for _, n := range net {
		nn := n.(map[string]interface{})
		for key, value := range nn {
			if key == "BOOTPROTO" {
				text += key + "=" + "static" + "\n"
			} else if key == "IPADDR" {
				text += key + "=" + ip + "\n"
			} else if key == "NETMASK" {
				text += key + "=" + netmask + "\n"
			} else if key == "GATEWAY" {
				text += key + "=" + gateway + "\n"
			} else if key == "DNS1" {
				text += key + "=" + dns1 + "\n"
			} else if key == "DNS2" && len(dns2) != 0 {
				text += key + "=" + dns2 + "\n"
			} else {
				text += key + "=" + value.(string) + "\n"
			}
		}
	}
	if ok := strings.Contains(text, "IPADDR"); !ok {
		t := "IPADDR" + "=" + ip + "\n"
		text += t
	}
	if ok := strings.Contains(text, "NETMASK"); !ok {
		t := "NETMASK" + "=" + netmask + "\n"
		text += t
	}
	if ok := strings.Contains(text, "GATEWAY"); !ok {
		t := "GATEWAY" + "=" + gateway + "\n"
		text += t
	}
	if ok := strings.Contains(text, "DNS1"); !ok {
		t := "DNS1" + "=" + dns1 + "\n"
		text += t
	}
	if ok := strings.Contains(text, "DNS2"); !ok {
		if len(dns2) != 0 {
			t := "DNS2" + "=" + dns2 + "\n"
			text += t
		}
	}
	return
}

func SaveFileToDatabase(file *model.Files) error {
	filename := file.FileName
	if len(filename) == 0 {
		return errors.New("请输入配置文件名字")
	}

	filepath := file.FilePath
	if len(filepath) == 0 {
		return errors.New("请输入下发文件路径")
	}
	temp, err := dao.IsExistFile(filename)
	if err != nil {
		return err
	}
	if temp {
		return errors.New("文件名字已存在，请重新输入")
	}

	filetype := file.Type
	if len(filetype) == 0 {
		return errors.New("请选择文件类型")
	}

	description := file.Description
	if len(description) == 0 {
		return errors.New("请添加文件描述")
	}

	batchId := file.ControlledBatch

	text := file.File
	if len(text) == 0 {
		return errors.New("请重新检查文件内容")
	}

	fd := model.Files{
		UserUpdate:      file.UserUpdate,
		UserDept:        file.UserDept,
		FileName:        filename,
		FilePath:        filepath,
		Type:            filetype,
		Description:     description,
		ControlledBatch: batchId,
		TakeEffect:      file.TakeEffect,
		File:            text,
	}
	return dao.SaveFile(fd)
}

func DeleteFile(fileIds []int) error {
	for _, fileId := range fileIds {
		err := dao.DeleteFile(fileId)
		if err != nil {
			logger.Error(err.Error())
		}
		err = dao.DeleteHistoryFile(fileId)
		if err != nil {
			logger.Error(err.Error())
		}
	}
	return nil
}
func UpdateFile(file *model.Files) error {
	id := file.ID
	err := dao.SaveHistoryFile(id)
	if err != nil {
		return err
	}
	user := file.UserUpdate
	userDept := file.UserDept
	filename := file.FileName
	description := file.Description
	batchId := file.ControlledBatch
	text := file.File
	ExistIdBool, err := dao.IsExistId(file.ID)
	if err != nil {
		return err
	}
	if !ExistIdBool {
		return errors.New("id有误,请重新确认该文件是否存在")
	}
	if ok, lastfileId, fileName, err := dao.IsExistFileLatest(id); ok {
		if err != nil {
			return err
		}
		fname := strings.Split(fileName, "-")
		f := model.HistoryFiles{
			FileName: fname[0],
		}
		err = dao.UpdateLastFile(lastfileId, f)
		if err != nil {
			return err
		}
	}
	f := model.Files{
		Type:            file.Type,
		FileName:        filename,
		FilePath:        file.FilePath,
		Description:     description,
		UserUpdate:      user,
		UserDept:        userDept,
		ControlledBatch: batchId,
		TakeEffect:      file.TakeEffect,
		File:            text,
	}
	return dao.UpdateFile(id, f)
}

func LastFileRollBack(file *model.RollBackFiles) error {
	lastfileId := file.HistoryFileID
	fileId := file.FileID
	user := file.UserUpdate
	userDept := file.UserDept
	lastfileText, err := dao.LastFileText(lastfileId)
	if err != nil {
		return err
	}
	if ok, _, _, err := dao.IsExistFileLatest(fileId); !ok {
		if err != nil {
			return err
		}
		err := dao.SaveLatestFile(fileId)
		if err != nil {
			return err
		}
	}
	fd := model.Files{
		UserUpdate: user,
		UserDept:   userDept,
		File:       lastfileText,
	}
	return dao.UpdateFile(fileId, fd)
}
