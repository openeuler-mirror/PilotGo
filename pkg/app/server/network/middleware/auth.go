package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service/auth"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
)

func AuthMiddleware(c *gin.Context) {
	claims, err := auth.ParseMyClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "claims解析错误:" + err.Error()})
		c.Abort()
		return
	}
	logger.Debug("request from %d, %s", claims.UserId, claims.UserName)

	ok, err := auth.CheckAuth(claims.UserName, c.Request.RequestURI, "get")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "check auth error:" + err.Error()})
		c.Abort()
		return
	}
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "no permission"})
		c.Abort()
		return
	}

	c.Next()
}

func NeedPermission(resource, action string) func(c *gin.Context) {
	return func(c *gin.Context) {
		claims, err := auth.ParseMyClaims(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "claims解析错误:" + err.Error()})
			c.Abort()
			return
		}
		logger.Debug("request from %d, %s", claims.UserId, claims.UserName)

		ok, err := auth.CheckAuth(claims.UserName, resource, action)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "check auth error:" + err.Error()})
			c.Abort()
			return
		}
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "no permission"})
			c.Abort()
			return
		}

		c.Next()
	}
}