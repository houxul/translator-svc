package provider

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

func request(method, url string, queryString, headers map[string]string, body interface{}) ([]byte, error) {
	// defer func() {
	// 	// TODO 记录日志
	// }()

	var bod io.Reader
	if body != nil {
		bs, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bod = bytes.NewReader(bs)
	}

	req, err := http.NewRequest(method, url, bod)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for k, v := range queryString {
		q.Set(k, v)
	}
	req.URL.RawQuery = q.Encode()

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bs, nil
}
