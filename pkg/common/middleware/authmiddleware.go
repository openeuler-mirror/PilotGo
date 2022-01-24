package middleware

/**
 * @Author: zhang han
 * @Date: 2021/11/1 11:14
 * @Description:
 */

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeluer.org/PilotGo/PilotGo/pkg/common"
	"openeluer.org/PilotGo/PilotGo/pkg/mysqlmanager"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization") //Get authorization header
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不够1"})
			c.Abort()
			return
		}

		tokenString = tokenString[7:]
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不够2"})
			c.Abort()
			return
		}

		userId := claims.UserId // Get the userID in claim after the verification is passed
		var user model.User
		mysqlmanager.DB.First(&user, userId)
		if user.ID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不够3"})
			c.Abort()
			return
		}
		c.Set("user", user) // user exists, write the user's information into the context
		c.Next()
	}
}
