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
 * Date: 2021-12-18 02:33:55
 * LastEditTime: 2022-04-11 16:15:51
 * Description: 用户登录、增删改查
 ******************************************************************************/
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/service"
	"openeluer.org/PilotGo/PilotGo/pkg/utils"
	"openeluer.org/PilotGo/PilotGo/pkg/utils/response"
)

func GetUserRole(c *gin.Context) {
	roles := dao.AllUserRole()
	response.Success(c, gin.H{"role": roles}, "获取用户角色")
}

func Register(c *gin.Context) {
	var user model.User
	c.Bind(&user)
	username := user.Username
	password := user.Password
	email := user.Email
	phone := user.Phone
	depart := user.DepartName
	departId := user.DepartSecond
	departPid := user.DepartFirst
	roleId := user.RoleID

	if len(username) == 0 { //Data verification
		username = service.RandomString(5)
	}
	if len(password) == 0 {
		password = "123456"
	}
	if len(email) == 0 {
		response.Response(c, http.StatusOK, 422, nil, "邮箱不能为空!")
		return
	}
	if dao.IsEmailExist(email) {
		response.Response(c, http.StatusOK, 422, nil, "邮箱已存在!")
		return
	}

	user_type := service.UserType(roleId)

	user = model.User{ //Create user
		Username:     username,
		Password:     password,
		Phone:        phone,
		Email:        email,
		DepartName:   depart,
		DepartFirst:  departPid,
		DepartSecond: departId,
		UserType:     user_type,
		RoleID:       roleId,
	}
	dao.AddUser(user)

	response.Success(c, nil, "添加用户成功!") //Return result
}

func Login(c *gin.Context) {
	var user model.User //Data verification
	c.Bind(&user)
	email := user.Email
	password := user.Password

	if !dao.IsEmailExist(email) {
		response.Response(c, http.StatusOK, 400, nil, "用户不存在!")
		return
	}

	DecryptedPassword, err := utils.JsAesDecrypt(password, email)
	if err != nil {
		response.Response(c, http.StatusOK, 400, nil, "密码解密失败")
		return
	}

	DBpassword, departName, roleId, departId, userType := dao.UserPassword(email)
	if DBpassword != DecryptedPassword {
		response.Response(c, http.StatusOK, 400, nil, "密码错误!")
		return
	}

	// Issue token
	token, err := service.ReleaseToken(user)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}
	response.Success(c, gin.H{"token": token, "departName": departName, "departId": departId, "userType": userType, "roleId": roleId}, "登陆成功!")
}

// 退出
func Logout(c *gin.Context) {

	response.Success(c, nil, "退出成功!")

}

func Info(c *gin.Context) {
	user, _ := c.Get("x-user")
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{"user": model.ToUserDto(user.(model.User))},
	})
}

// 查询所有用户
func UserAll(c *gin.Context) {
	query := &utils.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Response(c, http.StatusOK, 400, gin.H{"status": false}, err.Error())
		return
	}

	users := dao.UserAll()
	utils.Reverse(&users)

	total, data, err := utils.SearchAll(query, users)
	if err != nil {
		response.Response(c, http.StatusOK, 400, gin.H{"status": false}, err.Error())
		return
	}
	utils.JsonPagination(c, data, total, query)
}

// 高级搜索
func UserSearch(c *gin.Context) {
	var user model.User
	c.Bind(&user)
	var email = user.Email

	query := &utils.PaginationQ{}
	err := c.ShouldBindQuery(query)
	if err != nil {
		response.Response(c, http.StatusOK, 400, gin.H{"status": false}, err.Error())
		return
	}

	users := dao.UserSearch(email)
	utils.Reverse(&users)

	total, data, err := utils.SearchAll(query, users)
	if err != nil {
		response.Response(c, http.StatusOK, 400, gin.H{"status": false}, err.Error())
		return
	}
	utils.JsonPagination(c, data, total, query)
}

// 重置密码
func ResetPassword(c *gin.Context) {
	var user model.User
	c.Bind(&user)
	var email = user.Email
	u, err := dao.ResetPassword(email)
	if err != nil {
		response.Response(c, http.StatusOK, 400, nil, err.Error())
	} else {
		response.Response(c, http.StatusOK, 200, gin.H{"data": u}, "密码重置成功!")
	}
}

// 删除用户
func DeleteUser(c *gin.Context) {
	var userdel model.Userdel
	c.ShouldBind(&userdel)

	for _, userEmail := range userdel.Emails {
		dao.DeleteUser(userEmail)
	}
	response.Response(c, http.StatusOK, 200, nil, "用户删除成功!")
}

//修改用户信息
func UpdateUser(c *gin.Context) {
	var user model.User
	c.Bind(&user)
	email := user.Email
	phone := user.Phone
	Pid := user.DepartFirst
	id := user.DepartSecond
	departName := user.DepartName

	u := dao.UserInfo(email)

	if u.DepartName != departName && u.Phone != phone {
		dao.UpdateUserDepart(email, departName, Pid, id)
		dao.UpdateUserPhone(email, phone)
		response.Response(c, http.StatusOK, 200, gin.H{"data": user}, "用户信息修改成功")
		return
	}
	if u.DepartName == departName && u.Phone != phone {
		dao.UpdateUserPhone(email, phone)
		response.Response(c, http.StatusOK, 200, gin.H{"data": user}, "用户信息修改成功")
		return
	}
	if u.DepartName != departName && u.Phone == phone {
		dao.UpdateUserDepart(email, departName, Pid, id)
		response.Response(c, http.StatusOK, 200, gin.H{"data": user}, "用户信息修改成功")
	}
}

//一键导入用户数据
func ImportUser(c *gin.Context) {
	form, _ := c.MultipartForm()

	files := form.File["upload"]
	if files == nil {
		response.Response(c, http.StatusOK, 400, nil, "请先选择要上传的文件")
		return
	}
	UserExit := make([]string, 0)
	for _, file := range files {
		name := file.Filename
		c.SaveUploadedFile(file, name)

		xlFile, error := xlsx.OpenFile(name)
		if error != nil {
			return
		}

		UserExit = service.ReadFile(xlFile, UserExit)
	}

	if len(UserExit) == 0 {
		response.Response(c, http.StatusOK, 200, nil, "导入用户信息成功")
	} else {
		response.Response(c, http.StatusOK, 200, gin.H{"UserExit": UserExit}, "以上用户已经存在")
	}
}
