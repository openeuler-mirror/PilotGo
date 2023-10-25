/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: yangzhao1
 * Date: 2022-04-13 15:55:23
 * LastEditTime: 2022-04-18 16:00:11
 * Description: provide agent log manager of pilotgo
 ******************************************************************************/

package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// 将string写入到指定文件
func FileSaveString(filePath string, data string) error {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	data_length := len(data)
	send_count := 0
	for {
		n, err := f.WriteString(data[send_count:])
		if err != nil {
			return err
		}
		if n+send_count >= data_length {
			send_count += n
			break
		}
	}
	return nil
}

// 读取文件所有数据，返回字符串
func FileReadString(filePath string) (string, error) {
	f, err := os.Open(filePath)
	defer func(file *os.File) {
		// ignore file close error
		file.Close()
	}(f)
	if err != nil {
		return "", err
	}

	var result []byte
	readBuff := make([]byte, 1024*4)
	for {
		n, err := f.Read(readBuff)
		if err != nil {
			if err == io.EOF {
				if n != 0 {
					result = append(result, readBuff[:n]...)
				}
				break
			}
			return "", err
		}
		result = append(result, readBuff[:n]...)
	}
	return string(result), nil
}

// 获取文件状态，return (true, false):文件存在； (false, true):文件不存在但存在同名目录
func IsFileExist(filePath string) (bool, bool) {
	info, err := os.Stat(filePath)
	if err != nil {
		return false, false
	}

	if info.IsDir() {
		return false, true
	}

	return true, false
}

func GetFiles(filePath string, scanSub bool) (fs []string, err error) {
	files, err := os.ReadDir(filePath)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() && scanSub {
			tmp, err := GetFiles(filepath.Join(filePath, file.Name()), scanSub)
			if err != nil {
				return nil, err
			}
			fs = append(fs, tmp...)
		} else {
			fs = append(fs, file.Name())
		}
	}
	return fs, nil
}

func UpdateFile(path, filename, data interface{}) (lastversion string, err error) {
	fullname := path.(string) + "/" + filename.(string)
	if fok, dok := IsFileExist(fullname); !fok {
		if dok {
			return "", fmt.Errorf("存在同名目录,文件下发失败")
		}
		err := FileSaveString(fullname, data.(string))
		if err != nil {
			return "", err
		}
	}
	lastversion, err = FileReadString(fullname)
	if err != nil {
		return "", err
	}
	err = FileSaveString(fullname, data.(string))
	return lastversion, err
}
