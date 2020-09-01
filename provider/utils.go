package provider

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

func request(method, url string, queryString, headers map[string]string, bod io.Reader) ([]byte, error) {
	// defer func() {
	// 	// TODO 记录日志
	// }()

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

func genPairs(m map[string]string) []string {
	var pairs []string
	for k, v := range m {
		pairs = append(pairs, k+"="+v)
	}
	return pairs
}

func tencentSignature(key string, queryString map[string]string) string {
	pairs := genPairs(queryString)
	sort.StringSlice(pairs).Sort()
	pairsStr := "POSTtmt.tencentcloudapi.com/?" + strings.Join(pairs, "&")
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(pairsStr))
	baseStr := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return url.QueryEscape(baseStr)
}

var letterRunes = []rune("0123456789")

func randString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
