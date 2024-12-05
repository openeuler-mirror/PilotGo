/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: wubijie <wubijie@kylinos.cn>
 * Date: Thu Jan 11 20:24:26 2024 +0800
 */
package client

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
)

func (c *Client) RegisterPermission(pers []common.Permission) error {
	for _, v := range pers {
		if strings.Contains(v.Resource, "/") {
			return errors.New("permission-resource string contains /")
		} else {
			c.permissions = append(c.permissions, v)
		}
	}
	return nil
}

func (c *Client) HasPermission(resource, operate string) (bool, error) {
	if !c.IsBind() {
		return false, errors.New("unbind PilotGo-server platform")
	}

	url := "http://" + c.Server() + "/api/v1/pluginapi/permission"
	p := &common.Permission{
		Resource: resource,
		Operate:  operate,
	}

	r, err := httputils.Post(url, &httputils.Params{
		Body: p,
		Cookie: map[string]string{
			TokenCookie: c.token,
		},
	})
	if err != nil {
		return false, err
	}

	res := &struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
		Data    bool   `json:"data"`
	}{}
	if err := json.Unmarshal(r.Body, res); err != nil {
		return false, err
	}
	if res.Code != http.StatusOK {
		return false, errors.New(res.Message)
	}
	return res.Data, nil
}
