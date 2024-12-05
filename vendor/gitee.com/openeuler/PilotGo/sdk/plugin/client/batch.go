/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Fri Jan 19 11:08:30 2024 +0800
 */
package client

import (
	"encoding/json"
	"errors"

	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
)

func (c *Client) BatchList() ([]*common.BatchList, error) {
	if !c.IsBind() {
		return nil, errors.New("unbind PilotGo-server platform")
	}
	url := "http://" + c.Server() + "/api/v1/pluginapi/batch_list"
	r, err := httputils.Get(url, &httputils.Params{
		Cookie: map[string]string{
			TokenCookie: c.token,
		},
	})
	if err != nil {
		return nil, err
	}

	result := struct {
		Code int                 `json:"code"`
		Data []*common.BatchList `json:"data"`
	}{}
	if err := json.Unmarshal(r.Body, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

func (c *Client) BatchUUIDList(batchId string) ([]string, error) {
	if !c.IsBind() {
		return nil, errors.New("unbind PilotGo-server platform")
	}
	url := "http://" + c.Server() + "/api/v1/pluginapi/batch_uuid?batchId=" + batchId
	r, err := httputils.Get(url, &httputils.Params{
		Cookie: map[string]string{
			TokenCookie: c.token,
		},
	})
	if err != nil {
		return nil, err
	}

	result := struct {
		Code int      `json:"code"`
		Data []string `json:"data"`
	}{}
	if err := json.Unmarshal(r.Body, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}
