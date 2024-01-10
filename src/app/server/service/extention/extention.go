package extention

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
)

func RequestExtention(Url string) ([]common.Extention, error) {
	//TODO:规定插件接收请求的api
	url := Url + "/plugin_manage/api/v1/extentions"

	r, err := httputils.Get(url, &httputils.Params{})
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	if r.StatusCode != http.StatusOK {
		logger.Error("server process error:" + strconv.Itoa(r.StatusCode))
		return nil, err
	}

	resp := &common.CommonResult{}
	if err := json.Unmarshal(r.Body, resp); err != nil {
		logger.Error("解析结果出错%v", err.Error())
		return nil, err
	}
	if resp.Code != http.StatusOK {
		logger.Error(resp.Message)
		return nil, err
	}

	var data []map[string]interface{}
	if err := json.Unmarshal(resp.Data, &data); err != nil {
		panic(err)
	}
	extentions := common.ParseParameters(data)
	return extentions, err
}
