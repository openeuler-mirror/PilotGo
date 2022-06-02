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
 * Date: 2022-05-31 10:25:52
 * LastEditTime: 2022-06-02 10:16:10
 * Description: get config file
 ******************************************************************************/
package os

import (
	"io/ioutil"

	"openeluer.org/PilotGo/PilotGo/pkg/utils"
)

const (
	// repo路径
	RepoPath = "/etc/yum.repos.d"
)

func GetFiles(filePath string) (fs []string, err error) {
	files, err := ioutil.ReadDir(filePath)
	if err != nil {
		return fs, err
	}
	for _, file := range files {
		if file.IsDir() {
			tmp, err := GetFiles(filePath + "/" + file.Name())
			if err != nil {
				return fs, err
			}
			fs = append(fs, tmp...)
		} else {
			fs = append(fs, file.Name())
		}
	}
	return fs, nil
}

func UpdateFile(path, filename, data interface{}) (lastversion interface{}, err error) {
	fullname := path.(string) + "/" + filename.(string)
	if !utils.IsFileExist(fullname) {
		err := utils.FileSaveString(fullname, data.(string))
		if err != nil {
			return nil, err
		}
	}
	lastversion, err = utils.FileReadString(fullname)
	if err != nil {
		return nil, err
	}
	utils.RunCommand("rm -rf " + fullname)
	err = utils.FileSaveString(fullname, data.(string))
	if err != nil {
		return lastversion, err
	}
	return lastversion, nil
}
