package tag

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gitee.com/openeuler/PilotGo/app/server/service/plugin"
	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
)

type Tag struct {
	UUID string `json:"machineuuid"`
	Type string `json:"type"`
	Data string `json:"data"`
}

// 向所有插件发送uuidlist
func RequestTag(UUIDList []string) ([]Tag, error) {
	//TODO：获取在线插件列表
	plugins, err := plugin.GetPlugins()
	if err != nil {
		return nil, err
	}
	msg := []Tag{}
	//向url发送请求
	for _, v := range plugins {
		//TODO:规定插件接收请求的api
		url := v.Url + "/gettags"
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
		var tags []Tag
		if err := resp.ParseData(&tags); err != nil {
			logger.Error(err.Error())
		}
		msg = append(msg, tags...)
	}
	return msg, err
}