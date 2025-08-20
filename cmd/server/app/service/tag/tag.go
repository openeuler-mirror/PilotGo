/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package tag

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"gitee.com/openeuler/PilotGo/pkg/global"
	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
)

// 向所有插件发送uuidlist
func RequestTag(UUIDList []string) ([]common.Tag, error) {
	//TODO：获取在线插件列表
	p := global.GW.GetAllServices()

	msg := []common.Tag{}
	//向url发送请求
	for _, v := range p {
		//TODO:规定插件接收请求的api
		url := fmt.Sprintf("http://%s:%s/plugin_manage/api/v1/gettags", v["address"], v["port"])
		uuidTags := &struct {
			UUIDS []string `json:"uuids"`
		}{
			UUIDS: UUIDList,
		}
		r, err := httputils.Get(url, &httputils.Params{
			Body: uuidTags,
		})
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		if r.StatusCode != http.StatusOK {
			logger.Error("server process error:" + strconv.Itoa(r.StatusCode))
			continue
		}

		resp := &common.CommonResult{}
		if err := json.Unmarshal(r.Body, resp); err != nil {
			logger.Error("解析结果出错%v", err.Error())
			continue
		}
		if resp.Code != http.StatusOK {
			logger.Error(resp.Message)
			continue
		}
		var tags []common.Tag
		if err := resp.ParseData(&tags); err != nil {
			logger.Error(err.Error())
		}

		servicename := v["serviceName"].(string)
		for _, vt := range tags {
			vt.PluginName = servicename
			msg = append(msg, vt)
		}
	}
	return msg, nil
}
