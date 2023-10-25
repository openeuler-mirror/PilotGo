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
 * Date: 2022-01-17 02:43:29
 * LastEditTime: 2023-02-21 18:24:37
 * Description: provide agent rpm manager functions.
 ******************************************************************************/
package baseos

import (
	"bufio"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/utils"
	"gitee.com/openeuler/PilotGo/utils/os/common"
	"github.com/duke-git/lancet/fileutil"
)

// 获取全部安装的rpm包列表
func (b *BaseOS) GetAllRpm() ([]string, error) {
	listRpm := make([]string, 0)
	exitc, result, stde, err := utils.RunCommand("rpm -qa")
	if exitc == 0 && result != "" && stde == "" && err == nil {
		reader := strings.NewReader(result)
		scanner := bufio.NewScanner(reader)
		for {

			if !scanner.Scan() {
				break
			}
			line := scanner.Text()
			line = strings.TrimSpace(line)
			listRpm = append(listRpm, line)
		}
		return listRpm, nil
	}
	logger.Error("failed to get list of installed RPM packages: %d, %s, %s, %v", exitc, result, stde, err)
	return nil, fmt.Errorf("failed to get list of installed RPM packages: %d, %s, %s, %v", exitc, result, stde, err)

}

// 获取源软件包名以及源
func (b *BaseOS) GetRpmSource(rpm string) ([]common.RpmSrc, error) {
	Getlist := make([]common.RpmSrc, 0)
	listRpmSource := make([]string, 0)
	listRpmName := make([]string, 0)
	listRpmProvides := make([]string, 0)
	exitc, result, stde, err := utils.RunCommand("yum --nogpgcheck provides " + rpm)
	if exitc == 0 && result != "" && stde == "" && err == nil {
		reader := strings.NewReader(result)
		scanner := bufio.NewScanner(reader)
		for {
			if !scanner.Scan() {
				break
			}
			line := scanner.Text()
			line = strings.TrimSpace(line)
			reg2 := regexp.MustCompile(`^[R].*`)
			x := reg2.FindAllString(line, -1)
			if x == nil {
				continue
			}
			str2 := strings.Fields(x[0])
			listRpmSource = append(listRpmSource, str2[2])
		}
		reader = strings.NewReader(result)
		scanner = bufio.NewScanner(reader)
		for {
			if !scanner.Scan() {
				break
			}
			line := scanner.Text()
			line = strings.TrimSpace(line)
			reg1 := regexp.MustCompile(`^.*[.]+.*:`)
			x := reg1.FindAllString(line, -1)
			if x == nil {
				continue
			}
			str1 := strings.Fields(x[0])
			listRpmName = append(listRpmName, str1[0])
		}
		reader = strings.NewReader(result)
		scanner = bufio.NewScanner(reader)
		for {
			if !scanner.Scan() {
				break
			}
			line := scanner.Text()
			line = strings.TrimSpace(line)
			reg := regexp.MustCompile(`Provide.*:.*`)
			x := reg.FindAllString(line, -1)
			if x == nil {
				continue
			}
			str3 := strings.Fields(x[0])
			str := str3[2] + str3[3] + str3[4]
			listRpmProvides = append(listRpmProvides, str)
		}
		for key, value := range listRpmSource {
			tmp := common.RpmSrc{}
			tmp.Name = listRpmName[key]
			tmp.Provides = listRpmProvides[key]
			tmp.Repo = value
			Getlist = append(Getlist, tmp)
		}
		return Getlist, nil
	}
	logger.Error("failed to get source software package name and source: %d, %s, %s, %v", exitc, result, stde, err)
	return nil, fmt.Errorf("failed to get source software package name and source: %d, %s, %s, %v", exitc, result, stde, err)

}

// 按行使用正则语言查找结构体的属性信息
func readInfo(reader *strings.Reader, reg string) (string, error) {
	scanner := bufio.NewScanner(reader)
	var result string
	for {
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		line = strings.TrimSpace(line)
		reg := regexp.MustCompile(reg)
		x := reg.FindAllString(line, -1)
		if x == nil {
			continue
		}
		str := strings.Fields(x[0])
		length := len(str)
		if length < 3 {
			continue
		} else if length == 3 {
			result = str[2]
			return result, nil
		} else {
			i := 3
			result = str[2]
			for {
				if i == length {
					break
				}
				result = result + " " + str[i]
				i += 1

			}
			return result, nil
		}
	}
	return string(""), fmt.Errorf("failed to match struct properties")
}

func (b *BaseOS) GetRpmInfo(rpm string) (*common.RpmInfo, error) {
	rpminfo := common.RpmInfo{}

	exitc, result, stde, err := utils.RunCommand("rpm -qi " + rpm)
	if exitc == 1 && strings.Replace(result, "\n", "", -1) == "package bind is not installed" && stde == "" && err == nil {
		//未安装该软件包情况
		logger.Error(" %s's RPM package not installed: %d, %s, %s, %v", rpm, exitc, result, stde, err)
		return nil, fmt.Errorf("%s's RPM package not installed: %d, %s, %s, %v", rpm, exitc, result, stde, err)
	} else if exitc == 127 && result == "" && strings.Contains(stde, "command not found") && err == nil {
		//未安装rpm工具
		logger.Error("rpm not installed: %d, %s, %s, %v", exitc, result, stde, err)
		return nil, fmt.Errorf("rpm not installed: %d, %s, %s, %v", exitc, result, stde, err)
	} else if exitc == 0 && len(result) > 0 && stde == "" && err == nil {
		reader := strings.NewReader(result)
		str, err := readInfo(reader, `^Name.*`)
		if err != nil && len(str) != 0 {
			logger.Error("failed to read RPM package name properties")
			return nil, fmt.Errorf("failed to read RPM package name properties: %s", err)
		}
		rpminfo.Name = str
		reader = strings.NewReader(result)
		str, err = readInfo(reader, `^Version.*`)
		if err != nil && len(str) != 0 {
			logger.Error("failed to read RPM package Version properties")
			return nil, fmt.Errorf("failed to read RPM package Version properties: %s", err)
		}
		rpminfo.Version = str
		reader = strings.NewReader(result)
		str, err = readInfo(reader, `^Release.*`)
		if err != nil && len(str) != 0 {
			logger.Error("failed to read RPM package Release properties")
			return nil, fmt.Errorf("failed to read RPM package Release properties: %s", err)
		}
		rpminfo.Release = str
		// reader = strings.NewReader(result)
		// str, err = readInfo(reader, `^Architecture.*`)
		// if err != nil && len(str) != 0 {
		// 	logger.Error("读取rpm包Architecture属性失败")
		// 	return RpmInfo{}, fmt.Errorf("读取rpm包Architecture属性失败"), err
		// }
		rpminfo.Architecture = strings.Split(rpm, ".")[len(strings.Split(rpm, "."))-1]
		reader = strings.NewReader(result)
		str, err = readInfo(reader, `^Install Date.*`)
		if err != nil && len(str) != 0 {
			logger.Error("failed to read RPM package InstallDate properties")
			return nil, fmt.Errorf("failed to read RPM package InstallDate properties: %s", err)
		}
		rpminfo.InstallDate = str
		reader = strings.NewReader(result)
		str, err = readInfo(reader, `^Size.*`)
		if err != nil && len(str) != 0 {
			logger.Error("failed to read RPM package Size properties")
			return nil, fmt.Errorf("failed to read RPM package Size properties: %s", err)
		}
		rpminfo.Size = str
		reader = strings.NewReader(result)
		str, err = readInfo(reader, `^License.*`)
		if err != nil && len(str) != 0 {
			logger.Error("failed to read RPM package License properties")
			return nil, fmt.Errorf("failed to read RPM package License properties: %s", err)
		}
		rpminfo.License = str
		reader = strings.NewReader(result)
		str, err = readInfo(reader, `^Signature.*`)
		if err != nil && len(str) != 0 {
			logger.Error("failed to read RPM package Signature properties")
			return nil, fmt.Errorf("failed to read RPM package Signature properties: %s", err)
		}
		rpminfo.Signature = str
		reader = strings.NewReader(result)
		str, err = readInfo(reader, `^Packager.*`)
		if err != nil && len(str) != 0 {
			logger.Error("failed to read RPM package Packager properties")
			return nil, fmt.Errorf("failed to read RPM package Packager properties: %s", err)
		}
		rpminfo.Packager = str
		reader = strings.NewReader(result)
		str, err = readInfo(reader, `^Vendor.*`)
		if err != nil && len(str) != 0 {
			logger.Error("failed to read RPM package Vendor properties")
			return nil, fmt.Errorf("failed to read RPM package Vendor properties: %s", err)
		}
		rpminfo.Vendor = str
		reader = strings.NewReader(result)
		str, err = readInfo(reader, `^URL.*`)
		if err != nil && len(str) != 0 {
			logger.Error("failed to read RPM package URL properties")
			return nil, fmt.Errorf("failed to read RPM package URL properties: %s", err)
		}
		rpminfo.URL = str
		reader = strings.NewReader(result)
		str, err = readInfo(reader, `^Summary.*`)
		if err != nil && len(str) != 0 {
			logger.Error("failed to read RPM package Summary properties")
			return nil, fmt.Errorf("failed to read RPM package Summary properties:%s", err)
		}
		rpminfo.Summary = str
		return &rpminfo, nil
	} else {
		logger.Error("failed to get rpm info: %d, %s, %s, %v", exitc, result, stde, err)
		return nil, fmt.Errorf("failed to get rpm info: %d, %s, %s, %v", exitc, result, stde, err)
	}
}

// 判断rpm软件包是否安装/卸载成功
func verifyRpmInstalled(reader *strings.Reader, reg string) bool {
	scanner := bufio.NewScanner(reader)
	for {
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		line = strings.TrimSpace(line)
		reg2 := regexp.MustCompile(reg)
		x := reg2.FindAllString(line, -1)
		if x != nil {
			return true
		}
	}
	return false
}

// 安装rpm软件包
func (b *BaseOS) InstallRpm(rpm string) error {
	exitc, result, stde, err := utils.RunCommand("yum -y --nogpgcheck install " + rpm)
	if exitc == 0 && result != "" && err == nil {
		if verifyRpmInstalled(strings.NewReader(result), `Nothing to do.`) || verifyRpmInstalled(strings.NewReader(result), `无需任何处理。`) {
			logger.Error("failed to run RPM package installation command due to package already being installed")
			return fmt.Errorf("this RPM package is already installed")
		} else if verifyRpmInstalled(strings.NewReader(result), `^Error: Unable to find a match:.*`) {
			logger.Error("failed to run RPM package installation command due to the package not being found in the source")
			return fmt.Errorf("the package cannot be found in the source")
		} else {
			logger.Info("successfully installed %s", rpm)
			return nil
		}
	}
	logger.Error("failed to run RPM package installation command: %d, %s, %s, %v", exitc, result, stde, err)
	return fmt.Errorf("failed to run RPM package installation command: %d, %s, %s, %v", exitc, result, stde, err)

}

// 卸载rpm软件包
func (b *BaseOS) RemoveRpm(rpm string) error {
	exitc, result, stde, err := utils.RunCommand("yum -y --nogpgcheck remove " + rpm)
	if exitc == 0 && result != "" && err == nil {
		if verifyRpmInstalled(strings.NewReader(result), `Nothing to do.`) || verifyRpmInstalled(strings.NewReader(result), `无需任何处理。`) {
			logger.Error("failed to run RPM package uninstallation command due to the package not being found")
			return fmt.Errorf("this RPM package does not exist")
		} else {
			logger.Info("successfully uninstalled %s", rpm)
			return nil
		}
	}
	logger.Error("failed to run RPM package uninstallation command: %d, %s, %s, %v", exitc, result, stde, err)
	return fmt.Errorf("failed to execute RPM package uninstallation command: %d, %s, %s, %v", exitc, result, stde, err)

}

const RepoPath = "/etc/yum.repos.d"

// TODO: yum源文件在agent端打开的情况下调用该接口匹配内容出错
func (b *BaseOS) GetRepoSource() ([]*common.RepoSource, error) {
	files, err := utils.GetFiles(RepoPath, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get repo source file: %s", err)
	}

	var result []*common.RepoSource
	for _, path := range files {
		if strings.HasSuffix(path, ".repo") {
			content, err := fileutil.ReadFileToString(filepath.Join(RepoPath, path))
			if err != nil {
				return nil, err
			}
			repos, err := parseRepoContent(content)
			if err == nil {
				for index := range repos {
					repos[index].File = filepath.Base(path)
				}
				result = append(result, repos...)
			} else {
				return nil, err
			}
		}
	}

	return result, nil
}

func parseRepoContent(content string) ([]*common.RepoSource, error) {
	lines := strings.Split(content, "\n")

	var result []*common.RepoSource
	var currentRepo *common.RepoSource
	for index := range lines {
		line := strings.TrimSpace(lines[index])
		if line == "" || strings.HasPrefix(line, "#") {
			// 空行或者注释
			continue
		}

		if strings.HasPrefix(line, "[") {
			if strings.HasSuffix(line, "]") {
				// 找到一个新的repo
				if currentRepo != nil {
					result = append(result, currentRepo)
				}

				currentRepo = &common.RepoSource{
					ID: strings.Trim(line, "[]"),
				}
				continue
			} else {
				return nil, errors.New("invalid repo content: " + line)
			}
		}

		words := strings.SplitN(line, "=", 2)
		if len(words) == 2 {
			switch words[0] {
			case "name":
				currentRepo.Name = words[1]
			case "baseurl":
				currentRepo.BaseURL = words[1]
			case "mirrorlist":
				currentRepo.MirrorList = words[1]
			case "metalink":
				currentRepo.MetaLink = words[1]
			case "metadata_expire":
				currentRepo.MetadataExpire = words[1]
			case "enabled":
				if strings.TrimSpace(words[1]) == "1" {
					currentRepo.Enabled = 1
				} else {
					currentRepo.Enabled = 0
				}
			case "gpgcheck":
				if strings.TrimSpace(words[1]) == "1" {
					currentRepo.GPGCheck = 1
				} else {
					currentRepo.GPGCheck = 0
				}
			case "gpgkey":
				currentRepo.GPGKey = words[1]
			}
		}
	}
	if currentRepo.ID != "" {
		result = append(result, currentRepo)
	}

	return result, nil
}
