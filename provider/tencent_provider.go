package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type TencentTransResp struct {
	Response struct {
		TencentErrResp
		TencentCommonResp
	} `json:"Response"`
}

type TencentErrResp struct {
	Error *struct {
		Code    string `json:"Code"`
		Message string `json:"Message"`
	} `json:"Error"`
}

type TencentCommonResp struct {
	Source     string `json:"Source"`
	Target     string `json:"Target"`
	TargetText string `json:"TargetText"`
}

func tencentTranslate(srcs []string, en2zh bool) ([]string, error) {
	secretId := "AKIDKOFpDLcqITSfePIrhMlv5zfIfohKk8JD"
	secretKey := "yg9LJ4kBkFUUDk3T9e4E9n7BCh0DBwQe"
	queryString := map[string]string{
		"Action":   "TextTranslate",
		"Nonce":    randString(4),
		"Region":   "ap-beijing",
		"Language": "zh-CN",
		//"Method":     "POST",
		"Timestamp":  fmt.Sprintf("%d", time.Now().Unix()),
		"Version":    "2018-03-21",
		"SecretId":   secretId,
		"SourceText": strings.Join(srcs, "/"),
		"Source":     "en",
		"Target":     "zh",
		"ProjectId":  "0",
	}
	if !en2zh {
		queryString["Source"] = "zh"
		queryString["Target"] = "en"
	}
	queryString["Signature"] = tencentSignature(secretKey, queryString)

	pairs := genPairs(queryString)
	body := []byte(strings.Join(pairs, "&"))
	bs, err := request(http.MethodPost, "https://tmt.tencentcloudapi.com", nil, nil, body)
	if err != nil {
		return srcs, fmt.Errorf("tencent request err:(%w)", err)
	}

	var resp TencentTransResp
	if err := json.Unmarshal(bs, &resp); err != nil {
		return srcs, fmt.Errorf("json Unmarshal err:(%w)", err)
	}

	if resp.Response.Error != nil {
		return srcs, fmt.Errorf("tencent request result:(%s:%s)", resp.Response.Error.Code, resp.Response.Error.Message)
	}

	return strings.Split(resp.Response.TargetText, "/"), nil
}
