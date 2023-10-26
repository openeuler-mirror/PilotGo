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
 * LastEditTime: 2023-09-04 13:54:45
 * Description: 用户模块逻辑代码
 ******************************************************************************/
package user

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"gitee.com/openeuler/PilotGo/app/server/dao"
	"gitee.com/openeuler/PilotGo/app/server/service/common"
	"gitee.com/openeuler/PilotGo/utils"
	"github.com/tealeg/xlsx"
)

type User = dao.User

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
			userRole := row.Cells[4].Value         // 5：角色
			roleId, err := dao.GetRoleId(userRole) //角色对应id和用户类型
			if err != nil {
				return UserExit, err
			}
			password := strings.Split(email, "@")[0]
			u := &dao.User{
				Username:     userName,
				Phone:        phone,
				Password:     password,
				Email:        email,
				DepartName:   departName, //4：部门
				DepartFirst:  pid,
				DepartSecond: id,
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

func DeleteUser(Email string) error {
	if err := dao.DeleteUser(Email); err != nil {
		return err
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

func UpdatePassword(email, newPWD string) (dao.User, error) {

	u, err := dao.UpdatePassword(email, newPWD)
	if err != nil {
		return u, err
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
	if err != nil {
		return nil, 0, err
	}
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

func Login(user *dao.User) (string, int, string, error) {
	email := user.Email
	pwd := user.Password
	EmailBool, err := dao.IsEmailExist(email)
	if err != nil {
		return "", 0, "", err
	}
	if !EmailBool {
		return "", 0, "", errors.New("用户不存在")
	}

	u, err := dao.UserInfo(email)
	if err != nil {
		return "", 0, "", errors.New("查询邮箱密码错误")
	}

	err = utils.ComparePassword(u.Password, pwd)
	if err != nil {
		return "", 0, "", errors.New("密码错误")
	}

	return u.DepartName, u.DepartSecond, u.RoleID, nil
}

func Register(user *dao.User) error {
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

	u := &dao.User{ //Create user
		Username:     username,
		Password:     string(bs),
		Phone:        phone,
		Email:        email,
		DepartName:   depart,
		DepartFirst:  departPid,
		DepartSecond: departId,
		RoleID:       roleId,
	}
	err = dao.AddUser(u)
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

func GetUserRoles(username string) ([]string, error) {
	result := []string{}

	user, err := dao.QueryUserByName(username)
	if err != nil {
		return nil, err
	}

	// 查找角色
	roleids := user.RoleID
	roleId := strings.Split(roleids, ",")
	for _, id_str := range roleId {
		id, err := strconv.Atoi(id_str)
		if err != nil {
			return nil, err
		}
		role, err := dao.RoleIdToGetAllInfo(id)
		if err != nil {
			return nil, err
		}

		result = append(result, role.Role)
	}

	return result, nil
}
