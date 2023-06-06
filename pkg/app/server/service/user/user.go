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
 * Date: 2022-03-21 15:32:50
 * LastEditTime: 2022-04-21 15:37:48
 * Description: 用户模块逻辑代码
 ******************************************************************************/
package user

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service/auth"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service/common"
	"openeuler.org/PilotGo/PilotGo/pkg/global"
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
)

type User = dao.User
type Userdel = dao.Userdel

// 随机产生用户名字
func RandomString(n int) string {
	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := make([]byte, n)
	rand.Seed(time.Now().Unix())
	for k := range result {
		result[k] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

// 判断用户类型
func UserType(s string) int {
	// 找出字符串中包含的最小数字，例：“1,2,3,4”，最小的是1
	roleIds := strings.Split(s, ",")
	res := make([]int, len(s))

	for k, v := range roleIds {
		res[k], _ = strconv.Atoi(v)
	}

	min := res[0]
	if len(res) > 1 {
		for _, v := range res {
			if v < min {
				min = v
			}
		}
	}

	var user_type int
	if min > global.OrdinaryUserRoleId {
		user_type = global.OtherUserType
	} else {
		user_type = min - 1
	}
	return user_type
}

// 读取xlsx文件
func ReadFile(xlFile *xlsx.File, UserExit []string) ([]string, error) {
	for _, sheet := range xlFile.Sheets {
		for rowIndex, row := range sheet.Rows {
			//跳过第一行表头信息
			if rowIndex == 0 {
				continue
			}
			userName := row.Cells[0].Value //1:用户名
			phone := row.Cells[1].Value    //2：手机号
			email := row.Cells[2].Value    //3：邮箱
			EmailBool, err := dao.IsEmailExist(email)
			if err != nil {
				return nil, err
			}
			if EmailBool {
				UserExit = append(UserExit, email)
				continue
			}
			departName := row.Cells[3].Value            //4：部门
			pid, id, err := dao.GetPidAndId(departName) // 部门对应的PId和Id
			if err != nil {

				return UserExit, err
			}
			userRole := row.Cells[4].Value                               // 5：角色
			roleId, user_type, err := dao.GetRoleIdAndUserType(userRole) //角色对应id和用户类型
			if err != nil {
				return UserExit, err
			}
			password := global.DefaultUserPassword // 设置默认密码为123456
			u := dao.User{
				Username:     userName,
				Phone:        phone,
				Password:     password,
				Email:        email,
				DepartName:   departName, //4：部门
				DepartFirst:  pid,
				DepartSecond: id,
				UserType:     user_type,
				RoleID:       roleId,
			}
			err = dao.AddUser(u)
			if err != nil {
				return UserExit, err
			}
		}
	}
	return UserExit, nil
}

func DeleteUser(Emails []string) error {
	for _, userEmail := range Emails {
		if err := dao.DeleteUser(userEmail); err != nil {
			return err
		}
	}
	return nil
}

// 修改用户信息
func UpdateUser(user dao.User) (dao.User, error) {
	email := user.Email
	phone := user.Phone
	Pid := user.DepartFirst
	id := user.DepartSecond
	departName := user.DepartName
	u, err := dao.UserInfo(email)
	if err != nil {
		return u, err
	}

	if u.DepartName != departName && u.Phone != phone {
		err := dao.UpdateUserDepart(email, departName, Pid, id)
		if err != nil {
			return u, err
		}
		err = dao.UpdateUserPhone(email, phone)
		if err != nil {
			return u, err
		}
		return u, nil
	}
	if u.DepartName == departName && u.Phone != phone {
		err = dao.UpdateUserPhone(email, phone)
		if err != nil {
			return u, err
		}
		return u, nil
	}
	if u.DepartName != departName && u.Phone == phone {
		err := dao.UpdateUserDepart(email, departName, Pid, id)
		if err != nil {
			return u, err
		}
	}
	return u, nil
}

func ResetPassword(email string) (dao.User, error) {
	u, err := dao.ResetPassword(email)
	if err != nil {
		return u, err
	}
	return u, nil
}

func UserSearch(email string, query *common.PaginationQ) (interface{}, int, error) {
	users, total, err := dao.UserSearch(email)
	data, err := common.DataPaging(query, users, total)
	if err != nil {
		return nil, 0, err
	}
	return data, total, nil
}

func UserAll() ([]dao.ReturnUser, int, error) {
	users, total, err := dao.UserAll()
	if err != nil {
		return users, total, err
	}
	return users, total, nil
}

func Login(user dao.User) (string, string, int, int, string, error) {
	email := user.Email
	pwd := user.Password
	EmailBool, err := dao.IsEmailExist(email)
	if err != nil {
		return "", "", 0, 0, "", err
	}
	if !EmailBool {
		return "", "", 0, 0, "", errors.New("用户不存在")
	}

	u, err := dao.UserInfo(email)
	if err != nil {
		return "", "", 0, 0, "", errors.New("查询邮箱密码错误")
	}

	err = utils.ComparePassword(u.Password, pwd)
	if err != nil {
		return "", "", 0, 0, "", errors.New("密码错误")
	}

	// Issue token
	token, err := auth.ReleaseToken(user)
	if err != nil {
		return "", "", 0, 0, "", err
	}
	return token, u.DepartName, u.DepartSecond, u.UserType, u.RoleID, nil
}

func Register(user dao.User) error {
	username := user.Username
	password := user.Password
	email := user.Email
	phone := user.Phone
	depart := user.DepartName
	departId := user.DepartSecond
	departPid := user.DepartFirst
	roleId := user.RoleID

	if len(username) == 0 { //Data verification
		username = RandomString(5)
	}
	if len(password) == 0 {
		return errors.New("密码不能为空")
	}
	if len(email) == 0 {
		return errors.New("邮箱不能为空")
	}
	EmailBool, err := dao.IsEmailExist(email)
	if err != nil {
		return err
	}
	if EmailBool {
		return errors.New("邮箱已存在")
	}

	bs, err := utils.CryptoPassword(password)
	if err != nil {
		return errors.New("数据加密错误")
	}

	user_type := UserType(roleId)
	user = dao.User{ //Create user
		Username:     username,
		Password:     string(bs),
		Phone:        phone,
		Email:        email,
		DepartName:   depart,
		DepartFirst:  departPid,
		DepartSecond: departId,
		UserType:     user_type,
		RoleID:       roleId,
	}
	err = dao.AddUser(user)
	if err != nil {
		return err
	}
	return nil
}

func GetUserRole() ([]dao.UserRole, error) {
	roles, err := dao.AllUserRole()
	if err != nil {
		return roles, err
	}
	return roles, nil
}

func AddUserRole(userRole *dao.UserRole) error {
	err := dao.AddRole(*userRole)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUserRole(ID int) error {
	ok, err := dao.IsUserBindingRole(ID)
	if err != nil {
		return err
	}
	if !ok {
		err := dao.DeleteRole(ID)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("有用户绑定此角色，不可删除")
	}
}

func UpdateUserRole(UserRole *dao.UserRole) error {
	id := UserRole.ID
	role := UserRole.Role
	description := UserRole.Description
	userRole, err := dao.RoleIdToGetAllInfo(id)
	if err != nil {
		return err
	}
	if userRole.Role != role && userRole.Description != description {
		err = dao.UpdateRoleName(id, role)
		if err != nil {
			return err
		}
		err = dao.UpdateRoleDescription(id, description)
		if err != nil {
			return err
		}
		return nil
	}
	if userRole.Role == role && userRole.Description != description {
		err = dao.UpdateRoleDescription(id, description)
		if err != nil {
			return err
		}
		return nil
	}
	if userRole.Role != role && userRole.Description == description {
		err = dao.UpdateRoleName(id, role)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("没有修改信息")
}
