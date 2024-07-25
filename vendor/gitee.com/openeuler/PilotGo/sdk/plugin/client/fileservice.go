package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/utils/httputils"
)

func (client *Client) FileUpload(filePath string, filename string) error {
	// 以二进制方式上传文件
	file := filepath.Join(filePath, filename)
	bodyBuf, contentType, err := getUploadBody(file)
	if err != nil {
		return err
	}

	upload_addr := "http://" + client.Server() + "/api/v1/pluginapi/upload?filename=" + filename
	// 判断服务端是否是http协议
	ishttp, err := httputils.ServerIsHttp(upload_addr)
	if err != nil {
		return err
	}
	if !ishttp {
		upload_addr = fmt.Sprintf("https://%s", strings.Split(upload_addr, "://")[1])
	}

	req, err := http.NewRequest("POST", upload_addr, bodyBuf)
	if err != nil {
		return err
	}
	defer req.Body.Close()

	req.Header.Set("Content-Type", contentType)
	req.AddCookie(&http.Cookie{
		Name:  TokenCookie,
		Value: client.token,
	})

	hc := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}}
	defer hc.CloseIdleConnections()

	resp, err := hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取返回结果
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("没获取到：%s", err.Error())
		return err
	}
	res := &struct {
		StatusCode int         `json:"code"`
		Data       interface{} `json:"data"`
		Message    string      `json:"msg"`
	}{}
	err = json.Unmarshal(bs, &res)
	if err != nil {
		logger.Error("解析出错:%s", err.Error())
		return err
	}
	if resp.StatusCode == http.StatusOK && res.StatusCode == http.StatusOK {
		return nil
	}
	return errors.New(res.Message)
}

// 以二进制格式上传文件
func getUploadBody(filename string) (*bytes.Reader, string, error) {
	bodyBytes, err := os.ReadFile(filename)
	if err != nil {
		return bytes.NewReader(bodyBytes), "", err
	}
	return bytes.NewReader(bodyBytes), "multipart/form-data", nil
}
