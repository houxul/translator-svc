package provider

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type BaiduTransResp struct {
	BaiduErrResp
	BaiduCommonResp
}

type BaiduErrResp struct {
	ErrCode string `json:"error_code"`
	ErrMsg  string `json:"error_msg"`
}

type BaiduCommonResp struct {
	From        string `json:"from"`
	To          string `json:"to"`
	TransResult []struct {
		Src string `json:"src"`
		Dst string `json:"dst"`
	} `json:"trans_result"`
}

func baiduTranslate(srcs []string, en2zh bool) ([]string, error) {
	const (
		appid = "20200828000553619"
		key   = "6_VDtMyd41S_qX72AuON"
	)
	var (
		from = "en"
		to   = "zh"
	)
	if !en2zh {
		from, to = to, from
	}

	salt := fmt.Sprintf("%d", time.Now().Unix())
	query := strings.Join(srcs, "\n")
	sign := fmt.Sprintf("%x", md5.Sum([]byte(appid+query+salt+key)))
	url := "http://api.fanyi.baidu.com/api/trans/vip/translate"
	queryString := map[string]string{
		"q":     query,
		"from":  from,
		"to":    to,
		"appid": appid,
		"salt":  salt,
		"sign":  sign,
	}
	bs, err := request(http.MethodGet, url, queryString, nil, nil)
	if err != nil {
		return srcs, fmt.Errorf("baidu request err:(%w)", err)
	}

	var resp BaiduTransResp
	if err := json.Unmarshal(bs, &resp); err != nil {
		return srcs, fmt.Errorf("json Unmarshal err:(%w)", err)
	}

	if resp.ErrCode != "" {
		return srcs, fmt.Errorf("baidu request result:(%s:%s)", resp.ErrCode, resp.ErrMsg)
	}

	dsts := make([]string, len(resp.TransResult))
	for i, result := range resp.TransResult {
		dsts[i] = result.Dst
	}
	return dsts, nil
}
