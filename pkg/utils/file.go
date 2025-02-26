/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Mon Apr 18 16:28:19 2022 +0800
 */
package utils

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gitee.com/openeuler/PilotGo/sdk/common"
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

// 根据正则表达式读文件，并返回结构体
func ReadFilePattern(path, pattern string) ([]common.File, error) {
	var matchingFiles []common.File
	path = strings.TrimRight(path, "/") + "/"
	// 编译正则表达式
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	// 打开目录s
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// 读取目录中的文件
	files, err := dir.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	dir.Close()
	// 遍历文件并匹配正则表达式
	for _, file := range files {
		if re.MatchString(file) {
			// 符合正则表达式的文件，读取其内容
			f := common.File{
				Path: path,
				Name: file,
			}
			f.Content, err = FileReadString(path + file)
			if err != nil {
				return nil, err
			}
			matchingFiles = append(matchingFiles, f)
		}
	}
	return matchingFiles, nil
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
		//解码参数
		decodedate, err := base64.StdEncoding.DecodeString(data.(string))
		if err != nil {
			return "", err
		}
		err = FileSaveString(fullname, string(decodedate))
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
		//解码参数
		decodedate, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			return err
		}
		//更新文件
		err = FileSaveString(fullname, string(decodedate))
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
