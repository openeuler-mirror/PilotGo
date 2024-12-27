/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package user

import (
	"errors"
	"strings"

	"gitee.com/openeuler/PilotGo/cmd/server/app/service/internal/dao"
	"gitee.com/openeuler/PilotGo/pkg/utils"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/tealeg/xlsx"
)

type User = dao.User

type UserInfo struct {
	ID         uint   `json:"id"`
	DepartId   int    `json:"departId,omitempty"`
	DepartName string `json:"departName,omitempty"`
	Username   string `json:"username,omitempty" `
	Password   string `json:"password,omitempty"`
	Phone      string `json:"phone,omitempty"`
	Email      string `json:"email,omitempty" binding:"required" msg:"邮箱不能为空"`
	RoleID     string `json:"roleId"`
}

type ReturnUser struct {
	ID         uint     `json:"id"`
	DepartId   int      `json:"departId"`
	DepartName string   `json:"departName"`
	Username   string   `json:"username"`
	Phone      string   `json:"phone"`
	Email      string   `json:"email"`
	Roles      []string `json:"role"`
}

type UserDto struct {
	Name     string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

// 登录
func Login(user *UserInfo) (string, int, string, error) {
	email := user.Email
	pwd := user.Password
	EmailBool, err := dao.IsEmailExist(email)
	if err != nil {
		return "", 0, "", err
	}
	if !EmailBool {
		return "", 0, "", errors.New("用户不存在")
	}

	u, err := dao.GetUserByEmail(email)
	if err != nil {
		return "", 0, "", errors.New("查询邮箱密码错误")
	}

	err = utils.ComparePassword(u.Password, pwd)
	if err != nil {
		return "", 0, "", errors.New("密码错误")
	}
	depart, err := dao.GetDepartById(u.DepartId)
	if err != nil {
		return "", 0, "", errors.New("不存在此部门")
	}
	roleids, err := dao.GetRolesByUid(u.ID)
	return depart.Depart, u.DepartId, utils.Int2String(roleids), err
}

// 添加用户
func Register(user *UserInfo) error {
	username := user.Username
	password := user.Password
	email := user.Email
	phone := user.Phone
	//depart := user.DepartName
	departId := user.DepartId
	roleId := user.RoleID
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
		Username: username,
		Password: string(bs),
		Phone:    phone,
		Email:    email,
		DepartId: departId,
	}
	err = dao.AddUser(u)
	if err != nil {
		return err
	}
	roleIds := utils.String2Int(roleId)
	return dao.UpdateU2R(u.ID, roleIds)
}

func DeleteUser(Email string) error {
	//获取用户信息
	u, err := dao.GetUserByEmail(Email)
	if err != nil {
		return err
	}
	if Email == "admin" {
		return errors.New("admin用户不可删除")
	}
	//删除用户权限表的数据
	err = dao.DeleteByUid(u.ID)
	if err != nil {
		return err
	}
	//删除用户
	return dao.DeleteUser(Email)
}

// 修改用户信息
func UpdateUser(user UserInfo) error {
	new := dao.User{
		DepartId: user.DepartId,
		Username: user.Username,
		Phone:    user.Phone,
	}
	//获取用户信息
	_, err := dao.GetUserByEmail(user.Email)
	if err != nil {
		return err
	}
	//修改user表
	err = dao.UpdateUser(user.Email, new)
	if err != nil {
		return err
	}

	// TODO:
	/*
		//修改修改权限userrole表
		roleIds := utils.String2Int(user.RoleID)
		err = dao.UpdateU2R(u.ID, roleIds)*/
	return err
}

func UpdatePassword(email, newPWD string) error {
	return dao.UpdatePassword(email, newPWD)
}

func ResetPassword(email string) error {
	return dao.ResetPassword(email)
}

// 分页搜索所有用户
func UserSearchPaged(email string, offset, size int) (int64, []ReturnUser, error) {
	count, users, err := dao.UserSearchPaged(email, offset, size)
	if err != nil {
		return 0, nil, err
	}
	var returnUsers []ReturnUser
	for _, user := range users {
		depart, err := dao.GetDepartById(user.DepartId)
		if err != nil {
			logger.Error("%s", errors.New("不存在此部门").Error())
		}
		userinfo := ReturnUser{
			ID:         user.ID,
			DepartId:   user.DepartId,
			DepartName: depart.Depart,
			Username:   user.Username,
			Phone:      user.Phone,
			Email:      user.Email,
		}
		roleids, err := dao.GetRolesByUid(user.ID)
		if err != nil {
			logger.Error("%s", err.Error())
		}
		userinfo.Roles, err = dao.GetNamesByRoleIds(roleids)
		if err != nil {
			logger.Error("%s", err.Error())
		}
		returnUsers = append(returnUsers, userinfo)
	}
	return count, returnUsers, err
}

// 根据用户名字查询角色名字
func GetUserRoles(username string) ([]string, error) {
	user, err := dao.GetUserByName(username)
	if err != nil {
		return nil, err
	}
	// 查找角色
	roleids, err := dao.GetRolesByUid(user.ID)
	if err != nil {
		return nil, err
	}
	return dao.GetNamesByRoleIds(roleids)
}

// 根据userid查询user信息
func QueryUserByID(userID int) (*User, error) {
	user, err := dao.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}
	return user, err
}

// 查询某用户信息
func GetUserByEmail(email string) (*ReturnUser, error) {
	user, err := dao.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	depart, err := dao.GetDepartById(user.DepartId)
	if err != nil {
		logger.Error("%s", errors.New("不存在此部门").Error())
	}
	userinfo := &ReturnUser{
		ID:         user.ID,
		DepartId:   user.DepartId,
		DepartName: depart.Depart,
		Username:   user.Username,
		Phone:      user.Phone,
		Email:      user.Email,
	}
	roleids, err := dao.GetRolesByUid(user.ID)
	if err != nil {
		return nil, err
	}
	userinfo.Roles, err = dao.GetNamesByRoleIds(roleids)

	return userinfo, err
}

// 创建管理员账户
func CreateAdministratorUser() error {
	return dao.CreateAdministratorUser()
}

// 分页查询所有用户
func GetUserPaged(offset, size int) (int64, []ReturnUser, error) {
	var returnUsers []ReturnUser
	count, users, err := dao.GetUserPaged(offset, size)
	if err != nil {
		return 0, nil, err
	}
	for _, user := range users {
		depart, err := dao.GetDepartById(user.DepartId)
		if err != nil {
			logger.Error("%s", errors.New("不存在此部门").Error())
		}
		userinfo := ReturnUser{
			ID:         user.ID,
			DepartId:   user.DepartId,
			DepartName: depart.Depart,
			Username:   user.Username,
			Phone:      user.Phone,
			Email:      user.Email,
		}
		roleids, err := dao.GetRolesByUid(user.ID)
		if err != nil {
			logger.Error("%s", err.Error())
		}
		userinfo.Roles, err = dao.GetNamesByRoleIds(roleids)
		if err != nil {
			logger.Error("%s", err.Error())
		}
		returnUsers = append(returnUsers, userinfo)
	}
	return count, returnUsers, err
}

// 读取xlsx文件
func ReadFile(xlFile *xlsx.File, UserExit []string) ([]string, error) {
	for _, sheet := range xlFile.Sheets {
		for rowIndex, row := range sheet.Rows {
			//跳过第一行表头信息
			if rowIndex == 0 {
				continue
			}
			if !isEmptyRow(row) {
				userName := row.Cells[0].Value //1:用户名
				phone := row.Cells[1].Value    //2：手机号
				email := row.Cells[2].Value    //3：邮箱
				EmailBool, err := dao.IsEmailExist(email)
				if err != nil {
					UserExit = append(UserExit, email+"邮箱错误"+err.Error())
					logger.Error("%s", "邮箱错误"+err.Error())
					continue
				}
				if EmailBool {
					UserExit = append(UserExit, email+"已存在")
					continue
				}
				departName := row.Cells[3].Value          //4：部门
				_, id, err := dao.GetPidAndId(departName) // 部门对应的PId和Id
				if err != nil {
					UserExit = append(UserExit, email+"部门错误"+err.Error())
					logger.Error("%s", "部门错误"+err.Error())
					continue
				}
				userRole := row.Cells[4].Value         // 5：角色
				roleId, err := dao.GetRoleId(userRole) //角色对应id和用户类型
				if err != nil {
					UserExit = append(UserExit, email+"角色错误"+err.Error())
					logger.Error("%s", "角色错误"+err.Error())
					continue
				}
				password := strings.Split(email, "@")[0]

				u := &dao.User{
					Username: userName,
					Phone:    phone,
					Password: password,
					Email:    email,
					DepartId: id,
				}
				err = dao.AddUser(u)
				if err != nil {
					UserExit = append(UserExit, email+"添加用户错误"+err.Error())
					logger.Error("%s", "添加用户错误"+err.Error())
					continue
				}
				ur := &dao.UserRole{
					UserID: u.ID,
					RoleID: roleId,
				}
				err = ur.Add()
				if err != nil {
					UserExit = append(UserExit, email+"添加角色关系错误"+err.Error())
					logger.Error("%s", "添加角色关系错误"+err.Error())
					continue
				}
			}
		}
	}
	return UserExit, nil
}

func GetRolesByUid(uid uint) ([]int, error) {
	return dao.GetRolesByUid(uid)
}

// 判断是否为空行
func isEmptyRow(row *xlsx.Row) bool {
	for _, cell := range row.Cells {
		if cell.String() != "" {
			return false // 只要有一个单元格不为空，即认为当前行不为空行
		}
	}
	return true
}
