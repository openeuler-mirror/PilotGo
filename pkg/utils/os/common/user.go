/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Tue Feb 21 19:05:07 2023 +0800
 */
package common

// 获取当前用户信息
type CurrentUser struct {
	Username  string
	Userid    string
	GroupName string
	Groupid   string
	HomeDir   string
}

// 获取所有用户的信息
type AllUserInfo struct {
	Username    string
	UserId      string
	GroupId     string
	Description string
	HomeDir     string
	ShellType   string
}
