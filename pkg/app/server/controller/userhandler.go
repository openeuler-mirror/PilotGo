package controller

/**
 * @Author: zhang han
 * @Date: 2021/10/28 14:58
 * @Description: 用户登录和注册
 */

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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
	enable := user.Enable

	if len(username) == 0 { //Data verification
		username = utils.RandomString(5)
	}
	if len(password) == 0 {
		password = "123456"
	}
	if len(email) == 0 {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"邮箱不能为空!")
		return
	}
	if dao.IsEmailExist(email) {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			"邮箱已存在!")
		return
	}

	user = model.User{ //Create user
		Username: username,
		Password: password,
		Phone:    phone,
		Email:    email,
		Enable:   enable,
	}
	mysqlmanager.DB.Create(&user)

	response.Success(c, nil, "注册成功!") //Return result
}

func Login(c *gin.Context) {

	var user model.User //Data verification
	c.Bind(&user)
	email := user.Email
	password := user.Password

	mysqlmanager.DB.Where("email = ?", email).Find(&user)

	if user.ID == 0 {
		response.Response(c, http.StatusBadRequest,
			400,
			nil,
			"用户不存在!")
		return
	}

	bpassword := []byte(password)
	bemail := []byte(email)
	bpassword, _ = common.JsAesDecrypt(bpassword, bemail)
	btspassword := string(bpassword)
	if user.Password != btspassword {
		response.Response(c, http.StatusBadRequest,
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
	response.Success(c, gin.H{"token": token}, "登陆成功!")
}

func Info(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{"user": dto.ToUserDto(user.(model.User))},
	})
}

// 查询所有用户
func UserAll(c *gin.Context) {
	users := model.User{}
	query := &model.PaginationQ{}
	err := c.ShouldBindQuery(query)

	if model.HandleError(c, err) {
		return
	}
	list, total, err := users.All(query)
	if model.HandleError(c, err) {
		return
	}
	// 返回数据开始拼装分页的json
	model.JsonPagination(c, list, total, query)
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
type PasswordReset struct {
	Emails string `gorm:"type:varchar(30);not null" json:"email,omitempty" form:"email"`
}

func ResetPassword(c *gin.Context) {
	var user model.User
	var reset PasswordReset

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			err.Error())
		return
	}
	bodys := string(body)
	err = json.Unmarshal([]byte(bodys), &reset)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity,
			422,
			nil,
			err.Error())
		return
	}

	if dao.IsEmailExist(reset.Emails) {
		mysqlmanager.DB.Model(&user).Where("email=?", reset.Emails).Update("password", "123456")
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
	if dao.IsEmailExist(email) {
		// 修改手机号
		mysqlmanager.DB.Model(&user).Where("email=?", email).Update("phone", phone)

		response.Response(c, http.StatusOK,
			200,
			gin.H{"data": user},
			"User update successfully!")
		return
	} else {
		response.Fail(c, nil, "No user found!")
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
