/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"gitee.com/openeuler/PilotGo/cmd/server/app/config"
	userservice "gitee.com/openeuler/PilotGo/cmd/server/app/service/user"
)

var Issue = "PilotGo"

type UserClaims struct {
	jwt.StandardClaims

	UserId   uint
	UserName string
}

func GenerateUserToken(user userservice.ReturnUser) (string, error) {
	expirationTime := time.Now().Add(6 * 60 * time.Minute) //到期时间
	claims := &UserClaims{
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
	tokenString, err := token.SignedString([]byte(config.OptionsConfig.JWT.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseUser(c *gin.Context) (*userservice.User, error) {
	claims, err := parseMyClaims(c)
	if err != nil {
		return nil, err
	}

	user, err := userservice.QueryUserByID(int(claims.UserId))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func parseMyClaims(c *gin.Context) (*UserClaims, error) {
	cookie, err := c.Request.Cookie("Admin-Token") //Get authorization header
	if err != nil {
		return nil, err
	}

	claims, err := parseClaims(cookie.Value, &UserClaims{})
	if err != nil {
		return nil, err
	}
	m, ok := claims.(*UserClaims)
	if !ok {
		return nil, errors.New("invalid plugin claims")
	}
	return m, nil
}

type PluginClaims struct {
	jwt.StandardClaims

	Name string
	UUID string
}

func GeneratePluginToken(name, uuid string) (string, error) {
	claims := &PluginClaims{
		Name: name,
		UUID: uuid,

		StandardClaims: jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
			Issuer:   Issue,
			Subject:  "plugin token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.OptionsConfig.JWT.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParsePluginClaims(c *gin.Context) (*PluginClaims, error) {
	cookie, err := c.Request.Cookie("PluginToken") //Get authorization header
	if err != nil {
		return nil, err
	}

	claims, err := parseClaims(cookie.Value, &PluginClaims{})
	if err != nil {
		return nil, err
	}
	m, ok := claims.(*PluginClaims)
	if !ok {
		return nil, errors.New("invalid plugin claims")
	}
	return m, nil
}

func parseToken(tokenString string, clames jwt.Claims) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, clames, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(config.OptionsConfig.JWT.SecretKey), nil
	})
	return token, err
}

func parseClaims(tokenString string, claims jwt.Claims) (jwt.Claims, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("token is empty")
	}

	token, err := parseToken(tokenString, claims)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if token != nil && !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	return token.Claims, nil
}
