/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Tue Feb 21 00:17:56 2023 +0800
 */
package baseos

import (
	"bufio"
	"fmt"
	"os/user"
	"strings"

	"gitee.com/openeuler/PilotGo/pkg/utils"
	"gitee.com/openeuler/PilotGo/pkg/utils/os/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
)

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func (b *BaseOS) GetCurrentUserInfo() common.CurrentUser {
	u, err := user.Current()
	handleErr(err)
	userinfo := common.CurrentUser{
		Username:  u.Username,
		Userid:    u.Gid,
		GroupName: u.Name,
		Groupid:   u.Uid,
		HomeDir:   u.HomeDir,
	}
	return userinfo
}

func (b *BaseOS) GetAllUserInfo() ([]common.AllUserInfo, error) {
	exitc, tmp, stde, err := utils.RunCommand("cat /etc/passwd")
	if exitc == 0 && tmp != "" && stde == "" && err == nil {
		reader := strings.NewReader(tmp)
		scanner := bufio.NewScanner(reader)
		var allUsers []common.AllUserInfo
		for {

			if !scanner.Scan() {
				break
			}
			line := scanner.Text()
			line = strings.TrimSpace(line)
			strSlice := strings.Split(line, ":")
			users := common.AllUserInfo{
				Username:    strSlice[0],
				UserId:      strSlice[2],
				GroupId:     strSlice[3],
				Description: strSlice[4],
				HomeDir:     strSlice[5],
				ShellType:   strSlice[6],
			}
			allUsers = append(allUsers, users)
		}

		return allUsers, nil
	}
	logger.Error("failed to get passwd: %d, %s, %s, %v", exitc, tmp, stde, err)
	return nil, fmt.Errorf("failed to get passwd: %d, %s, %s, %v", exitc, tmp, stde, err)
}

// 创建新的用户，并修改新用户密码
func (b *BaseOS) AddLinuxUser(username, password string) error {
	exitc1, output1, stde1, err1 := utils.RunCommand("useradd -m " + username)
	if exitc1 == 0 && output1 == "" && stde1 == "" && err1 == nil {
		logger.Info("successfully created user: %s", output1)

		//下面两个是管道的两端
		//linux可以使用  echo "password" | passwd --stdin username
		//直接更改密码
		exitc2, output2, stde2, err2 := utils.RunCommand("echo \"" + password + "\" | passwd --stdin " + username)
		if exitc2 == 0 && output2 != "" && stde2 == "" && err2 == nil {
			logger.Info("successfully changed user passwd: %s", output2)
			return nil
		}
		logger.Error("failed to change passwd: %d, %s, %s, %v", exitc2, output2, stde2, err2)
		return fmt.Errorf("failed to change passwd: %d, %s, %s, %v", exitc2, output2, stde2, err2)
	}
	logger.Error("failed to add user: %d, %s, %s, %v", exitc1, output1, stde1, err1)
	return fmt.Errorf("failed to add user: %d, %s, %s, %v", exitc1, output1, stde1, err1)

}

// 删除用户
func (b *BaseOS) DelUser(username string) (string, error) {
	exitc, tmp, stde, err := utils.RunCommand(fmt.Sprintf("userdel -r %s", username))
	if exitc == 0 && tmp == "" && stde == "" && err == nil {
		logger.Info("successfully deleted user: %s", tmp)
		return tmp, nil
	}
	logger.Error("failed to delete user: %d, %s, %s, %v", exitc, tmp, stde, err)
	return "", fmt.Errorf("failed to delete user: %d, %s, %s, %v", exitc, tmp, stde, err)
}

// chmod [-R] 权限值 文件名
func (b *BaseOS) ChangePermission(permission, file string) (string, error) {
	exitc, tmp, stde, err := utils.RunCommand(fmt.Sprintf("chmod %s %s", permission, file))
	if exitc == 0 && tmp == "" && stde == "" && err == nil {
		logger.Info("successfully changed file permissions: %s", tmp)
		return tmp, nil
	}
	logger.Error("failed to change file permissions: %d, %s, %s, %v", exitc, tmp, stde, err)
	return "", fmt.Errorf("failed to change file permissions: %d, %s, %s, %v", exitc, tmp, stde, err)
}

// chown [-R] 所有者 文件或目录
func (b *BaseOS) ChangeFileOwner(user, file string) (string, error) {
	exitc, tmp, stde, err := utils.RunCommand(fmt.Sprintf("chown -R %s %s", user, file))
	if exitc == 0 && tmp == "" && stde == "" && err == nil {
		logger.Info("successfully changed file owner: %s", tmp)
		return tmp, nil
	}
	logger.Error("failed to change file owner: %d, %s, %s, %v", exitc, tmp, stde, err)
	return "", fmt.Errorf("failed to change file owner: %d, %s, %s, %v", exitc, tmp, stde, err)
}
