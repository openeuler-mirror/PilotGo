package httputils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Deprecated: use other api. TODO: will remove this
func Request(method, url string) ([]byte, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	hc := &http.Client{}
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func request(method, url string, param *Params) ([]byte, error) {
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
	if param != nil && len(param.Header) > 0 {
		for k, v := range param.Header {
			req.Header.Add(k, v)
		}
	}

	hc := &http.Client{}
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

type Params struct {
	// Header 参数会被添加到请求header当中
	Header map[string]string
	// Form 参数会被格式化到url当中
	Form map[string]string
	// Body 参数会被序列化成json字符串
	Body interface{}
}

func Post(url string, params *Params) ([]byte, error) {
	return request("POST", url, params)
}

func Get(url string, params *Params) ([]byte, error) {
	return request("GET", url, params)
}

func Put(url string, params *Params) ([]byte, error) {
	return request("PUT", url, params)
}

func Delete(url string, params *Params) ([]byte, error) {
	return request("DELETE", url, params)
}
