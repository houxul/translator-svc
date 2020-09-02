package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type LingocloudTransResp struct {
	Confidence float32  `json:"confidence"`
	Target     []string `json:"target"`
	Rc         int      `json:"rc"`
}

func lingocloudTranslate(srcs []string) ([]string, error) {
	token := "tg6jeai9s80m12anug0x"
	url := "http://api.interpreter.caiyunai.com/v1/translator"
	payload := map[string]interface{}{
		"source":     srcs,
		"trans_type": "en2zh",
		"request_id": "demo",
		"detect":     true,
	}

	headers := map[string]string{
		"content-type":    "application/json",
		"x-authorization": "token " + token,
	}

	bs, err := json.Marshal(payload)
	if err != nil {
		return srcs, fmt.Errorf("lingocloud request err:(%w)", err)
	}
	bod := bytes.NewReader(bs)

	bs, err = request(http.MethodPost, url, nil, headers, bod)
	if err != nil {
		return srcs, fmt.Errorf("lingocloud request err:(%w)", err)
	}

	var resp LingocloudTransResp
	if err := json.Unmarshal(bs, &resp); err != nil {
		return srcs, fmt.Errorf("json Unmarshal err:(%w)", err)
	}

	if resp.Rc != 0 {
		return srcs, fmt.Errorf("lingocloud request result:(%v)", resp)
	}

	return resp.Target, nil
}
