/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Wed Sep 27 17:35:12 2023 +0800
 */
package client

import (
	"encoding/json"
	"errors"
	"fmt"

	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/plugin/jwt"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
)

func (c *Client) ApplyConfig(batch *common.Batch, path, content string) error {
	serverInfo, err := c.Registry.Get("pilotgo-server")
	if err != nil {
		return err
	}
	url := fmt.Sprintf("http://%s:%s/api/v1/pluginapi/apply_config", serverInfo.Address, serverInfo.Port)

	r, err := httputils.Put(url, &httputils.Params{
		Cookie: map[string]string{
			jwt.TokenCookie: c.token,
		},
	})
	if err != nil {
		return err
	}

	resp := &struct {
		Status string
		Error  string
	}{}
	if err := json.Unmarshal(r.Body, resp); err != nil {
		return err
	}
	if resp.Status != "ok" {
		return errors.New(resp.Error)
	}

	return nil
}
