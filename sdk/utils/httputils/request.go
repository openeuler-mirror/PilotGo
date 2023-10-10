package httputils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func request(method, url string, param *Params) (*Response, error) {
	// 处理form参数
	if param != nil && len(param.Form) > 0 {
		s := ""
		for k, v := range param.Form {
			s += fmt.Sprintf("&%s=%s", k, v)
		}
		url = url + "?" + s[1:]
	}

	// 处理body参数
	var bodyReader io.Reader
	if param != nil && param.Body != nil {
		bodyBytes, err := json.Marshal(param.Body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	// 处理header参数
	if param != nil {
		if len(param.Header) > 0 {
			for k, v := range param.Header {
				req.Header.Add(k, v)
			}
		}

		// 如果存在body数据，则自动加入到Content-Type信息当中
		if param.Body != nil {
			typeStr := req.Header.Get("Content-Type")
			if typeStr == "" {
				req.Header.Set("Content-Type", "application/json")
			} else {
				if !strings.Contains(typeStr, "application/json") {
					req.Header.Set("Content-Type", typeStr+"; application/json")
				}
			}
		}
	}

	hc := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}}
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Response{
			StatusCode: resp.StatusCode,
		}, err
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       bs,
	}, nil
}

type Params struct {
	// Header 参数会被添加到请求header当中
	Header map[string]string
	// Form 参数会被格式化到url当中
	Form map[string]string
	// Body 参数会被序列化成json字符串
	Body interface{}
}

type Response struct {
	// 返回状态码
	StatusCode int
	// 返回body数组，[]byte
	Body []byte
}

func Post(url string, params *Params) (*Response, error) {
	return request("POST", url, params)
}

func Get(url string, params *Params) (*Response, error) {
	return request("GET", url, params)
}

func Put(url string, params *Params) (*Response, error) {
	return request("PUT", url, params)
}

func Delete(url string, params *Params) (*Response, error) {
	return request("DELETE", url, params)
}
