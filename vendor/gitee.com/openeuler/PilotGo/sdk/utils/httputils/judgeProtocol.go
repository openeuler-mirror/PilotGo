package httputils

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
)

func ServerIsHttp(rawurl string) (bool, error) {
	url, err := url.Parse(rawurl)
	if err != nil {
		return false, err
	}

	url_string := fmt.Sprintf("http://%s", net.JoinHostPort(url.Hostname(), url.Port()))
	req, err := http.NewRequest("GET", url_string, nil)
	if err != nil {
		return false, err
	}

	hc := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	resp, err := hc.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	respbytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != 200 && strings.Contains(string(respbytes), "Client sent an HTTP request to an HTTPS server") {
		return false, nil
	}
	return true, nil
}
