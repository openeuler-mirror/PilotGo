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
 * LastEditTime: 2022-03-21 17:16:17
 * Description: 用户登录、增删改查
 ******************************************************************************/
package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/tealeg/xlsx"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/common"
	"openeluer.org/PilotGo/PilotGo/pkg/common/dto"
	"openeluer.org/PilotGo/PilotGo/pkg/common/response"
	"openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"
	"openeluer.org/PilotGo/PilotGo/pkg/utils"
)

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
		username = utils.RandomString(5)
	}
	if len(password) == 0 {
		password = "123456"
	}
	if len(email) == 0 {
		response.Response(c, http.StatusOK,
			422,
			nil,
			"邮箱不能为空!")
		return
	}
	if dao.IsEmailExist(email) {
		response.Response(c, http.StatusOK,
			422,
			nil,
			"邮箱已存在!")
		return
	}
	if departPid == 1 {
		user.UserType = 1
	} else {
		user.UserType = 2
	}
	user = model.User{ //Create user
		Username:     username,
		Password:     password,
		Phone:        phone,
		Email:        email,
		DepartName:   depart,
		DepartFirst:  departPid,
		DepartSecond: departId,
		UserType:     user.UserType,
		RoleID:       roleId,
	}
	mysqlmanager.DB.Save(&user)

	response.Success(c, nil, "添加用户成功!") //Return result
}

func Login(c *gin.Context) {

	var user model.User //Data verification
	c.Bind(&user)
	email := user.Email
	password := user.Password
	mysqlmanager.DB.Where("email = ?", email).Find(&user)

	if user.ID == 0 {
		response.Response(c, http.StatusOK,
			400,
			nil,
			"用户不存在!")
		return
	}
	bpassword := []byte(password)
	bemail := []byte(email)
	bbpassword, err := common.JsAesDecrypt(bpassword, bemail)
	if err != nil {
		response.Response(c, http.StatusOK,
			400,
			nil,
			"密码解密失败")
		return
	}
	btspassword := string(bbpassword)
	if user.Password != btspassword {
		response.Response(c, http.StatusOK,
			400,
			nil,
			"密码错误!")
		return
	}
	token, err := common.ReleaseToken(user) // Issue token
	if err != nil {
		response.Response(c, http.StatusInternalServerError,
			500,
			nil,
			"服务器内部错误!")
		log.Printf("token生成错误:%v", err)
		return
	}
	response.Success(c, gin.H{"token": token, "departName": user.DepartName, "departId": user.DepartSecond, "userType": user.UserType, "roleId": user.RoleID}, "登陆成功!")
}

// 退出
func Logout(c *gin.Context) {

	response.Success(c, nil, "退出成功!")

}

func Info(c *gin.Context) {
	user, _ := c.Get("x-user")
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{"user": dto.ToUserDto(user.(model.User))},
	})
}

// 查询所有用户
func UserAll(c *gin.Context) {
	var users []model.User
	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)

	if err != nil {
		response.Response(c, http.StatusOK, 400, gin.H{"status": false}, err.Error())
		return
	}
	mysqlmanager.DB.Find(&users)
	datas := make([]map[string]interface{}, 0)
	for _, user := range users {
		data := make(map[string]interface{})
		data["id"] = user.ID
		data["departPId"] = user.DepartFirst
		data["departid"] = user.DepartSecond
		data["departName"] = user.DepartName
		data["username"] = user.Username
		data["phone"] = user.Phone
		data["email"] = user.Email
		data["userType"] = user.UserType
		roleids := user.RoleID
		roleId := strings.Split(roleids, ",")
		var roles []string
		for _, id := range roleId {
			userRole := model.UserRole{}
			i, err := strconv.Atoi(id)
			if err != nil {
				response.Response(c, http.StatusOK, 400, gin.H{"status": false}, err.Error())
				return
			}
			mysqlmanager.DB.Where("id = ?", i).Find(&userRole)
			role := userRole.Role
			roles = append(roles, role)
		}
		data["role"] = roles
		datas = append(datas, data)
	}
	common.Reverse(&datas)
	total, data, err := model.SearchAll(query, datas)
	if err != nil {
		response.Response(c, http.StatusOK, 400, gin.H{"status": false}, err.Error())
		return
	}
	model.JsonPagination(c, data, total, query)
}

// 高级搜索
func UserSearch(c *gin.Context) {
	var user model.User
	var users []model.User
	c.Bind(&user)
	var email = user.Email

	mysqlmanager.DB.Where("email LIKE ?", "%"+email+"%").Find(&users)
	response.Response(c, http.StatusOK,
		200,
		gin.H{"data": users},
		"查询成功!")
}

// 重置密码
func ResetPassword(c *gin.Context) {
	var user model.User
	email := c.Query("email")

	if dao.IsEmailExist(email) {
		mysqlmanager.DB.Model(&user).Where("email=?", email).Update("password", "123456")
		response.Response(c, http.StatusOK,
			200,
			gin.H{"data": user},
			"密码重置成功!")
		return
	} else {
		response.Fail(c, nil, "无此用户!")
	}
}

// 删除用户
type Userdel struct {
	Emails []string `gorm:"type:varchar(30);not null" json:"email,omitempty" form:"email"`
}

func DeleteUser(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			err.Error())
		return
	}
	var userdel Userdel
	bodys := string(body)
	err = json.Unmarshal([]byte(bodys), &userdel)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			err.Error())
		return
	}
	var user model.User
	for _, userEmail := range userdel.Emails {
		mysqlmanager.DB.Where("email=?", userEmail).Unscoped().Delete(user)
	}
	response.Response(c, http.StatusOK,
		200,
		nil,
		"用户删除成功!")
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
	mysqlmanager.DB.Where("email = ?", email).Find(&user)
	if user.DepartName != departName && user.Phone != phone {
		u := model.User{
			DepartFirst:  Pid,
			DepartSecond: id,
			DepartName:   departName,
		}
		fmt.Println("1", Pid, id, departName, phone)
		mysqlmanager.DB.Model(&user).Where("email=?", email).Updates(&u)
		mysqlmanager.DB.Model(&user).Where("email=?", email).Update("phone", phone)
		response.Response(c, http.StatusOK,
			200,
			gin.H{"data": user},
			"用户信息修改成功")
		return
	}
	if user.DepartName != departName && user.Phone != phone {
		fmt.Println("2", Pid, id, departName, phone)
		mysqlmanager.DB.Model(&user).Where("email=?", email).Update("phone", phone)
		response.Response(c, http.StatusOK,
			200,
			gin.H{"data": user},
			"用户信息修改成功")
		return
	}
	if user.DepartName != departName && user.Phone == phone {
		fmt.Println("3", Pid, id, departName, phone)
		u := model.User{
			DepartFirst:  Pid,
			DepartSecond: id,
			DepartName:   departName,
		}
		mysqlmanager.DB.Model(&user).Where("email=?", email).Updates(&u)
		response.Response(c, http.StatusOK,
			200,
			gin.H{"data": user},
			"用户信息修改成功")
	}
}

//一键导入用户数据
func ImportUser(c *gin.Context) {
	form, _ := c.MultipartForm()

	files := form.File["upload"]
	if files == nil {
		response.Fail(c, nil, "Please select a file first!")
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
		for _, sheet := range xlFile.Sheets {

			for rowIndex, row := range sheet.Rows {
				user := model.User{}
				depart := model.DepartNode{}

				//跳过第一行表头信息
				if rowIndex == 0 {
					continue
				}
				user.Username = row.Cells[0].Value
				user.Phone = row.Cells[1].Value
				user.Email = row.Cells[2].Value
				if dao.IsEmailExist(user.Email) {
					UserExit = append(UserExit, user.Email)
					continue
				}
				// 设置默认密码为123456
				user.Password = "123456"
				user.DepartName = row.Cells[3].Value
				mysqlmanager.DB.Where("depart=?", user.DepartName).Find(&depart)
				user.DepartSecond = depart.ID
				user.DepartFirst = depart.PID
				mysqlmanager.DB.Save(&user)
			}
		}
	}
	if len(UserExit) == 0 {
		response.Response(c, http.StatusOK,
			200,
			nil,
			"import success")
	} else {
		response.Response(c, http.StatusOK,
			200, gin.H{
				"UserExit": UserExit,
			}, "以上用户已经存在")
	}

}
