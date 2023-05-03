package utils

import (
	"io"
	"net/http"
)

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
