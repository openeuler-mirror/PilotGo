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
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 将string写入到指定文件
func FileSaveString(filePath string, data string) error {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	decodedate, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}
	data_length := len(decodedate)
	send_count := 0
	for {
		n, err := f.Write(decodedate[send_count:])
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
	//判断path是否带有斜杠
	fullname := strings.TrimRight(path.(string), "/") + "/" + filename.(string)
	fok, lok := IsFileExist(fullname)
	if fok {
		lastversion, err = FileReadString(fullname)
		if err != nil {
			return "", err
		}
		err := FileSaveString(fullname, data.(string))
		if err != nil {
			return "", err
		}
		return lastversion, err
	}
	if lok {
		return "", fmt.Errorf(fullname + " is a directory")
	}
	return "", fmt.Errorf(fullname + " does not exist")
}

// 存储文件，文件存在时直接更新内容，文件不存在路径存在时创建新文件，路径不存在时报错
func SaveFile(path, filename, data string) error {
	path = strings.TrimRight(path, "/")
	fullname := path + "/" + filename
	fok, _ := IsFileExist(fullname)
	if fok {
		//更新文件
		err := FileSaveString(fullname, data)
		return err
	} else {
		//判断是否是文件夹
		_, lok := IsFileExist(path)
		if lok {
			//创建文件
			file, err := os.Create(fullname)
			if err != nil {
				return err
			}
			defer file.Close()

			// 写入内容到文件
			decodedate, err := base64.StdEncoding.DecodeString(data)
			data_length := len(decodedate)
			send_count := 0
			for {
				n, err := file.Write(decodedate[send_count:])
				if err != nil {
					return err
				}
				if n+send_count >= data_length {
					send_count += n
					break
				}
			}
			return err
		}
		return fmt.Errorf(path + " does not exist")
	}
}
