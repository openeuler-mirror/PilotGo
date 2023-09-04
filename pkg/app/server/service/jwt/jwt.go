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
 * Date: 2021-11-1 15:08:08
 * LastEditTime: 2023-09-04 11:09:10
 * Description: jwt是一个基于token的轻量级认证方式
 ******************************************************************************/
package jwt

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/dao"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/service/common"
)

var jwtKey = []byte("a_secret_crect")
var Issue = "PilotGo"

type MyClaims struct {
	jwt.StandardClaims

	UserId   uint
	UserName string
}

func ReleaseToken(user dao.User) (string, error) {
	expirationTime := time.Now().Add(6 * 60 * time.Minute) //到期时间
	claims := &MyClaims{
		UserId:   user.ID,
		UserName: user.Username,

		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    Issue,
			Subject:   "user token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	return token, err
}

func ParseUser(c *gin.Context) (*common.User, error) {
	claims, err := ParseMyClaims(c)
	if err != nil {
		return nil, err
	}

	user, err := dao.QueryUserByID(int(claims.UserId))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func ParseMyClaims(c *gin.Context) (*MyClaims, error) {
	var tokenString string
	var token *jwt.Token
	var cookie *http.Cookie
	var claims *MyClaims

	var err error
	var ok bool

	cookie, err = c.Request.Cookie("Admin-Token") //Get authorization header
	if err != nil {
		goto OnError
	}
	tokenString = cookie.Value
	if tokenString == "" {
		err = fmt.Errorf("token is empty")
		goto OnError
	}

	token, err = ParseToken(tokenString)
	if err != nil {
		goto OnError
	}

	if token != nil && !token.Valid {
		err = fmt.Errorf("token is invalid")
		goto OnError
	}

	claims, ok = token.Claims.(*MyClaims)
	if !ok {
		err = fmt.Errorf("token claims is invalid")
		goto OnError
	}
	return claims, nil

OnError:
	return nil, err
}
