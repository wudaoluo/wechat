package oauth2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/wudaoluo/wechat/internal/debug/api"
	"github.com/wudaoluo/wechat/util"
	"net/http"
	"net/url"
)

type generateSchemeResp struct {
	ErrCode  int    `json:"errcode"`
	ErrMsg   string `json:"errmsg"`
	Openlink string `json:"openlink"`
}

type generateSchemeJumpWxa struct {
	Path       string `json:"path"`
	Query      string `json:"query"`
	EnvVersion string `json:"env_version"`
}

type generateSchemeReq struct {
	JumpWxa  generateSchemeJumpWxa `json:"jump_wxa"`
	IsExpire bool                  `json:"is_expire"`
}

// 获取 URL Scheme
func GetUrlScheme(accessToken string, path, query, env_version string) (info string, err error) {
	httpClient := util.DefaultHttpClient

	body := generateSchemeReq{
		JumpWxa: generateSchemeJumpWxa{
			Path:       path,
			Query:      query,
			EnvVersion: env_version,
		},
	}
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(body)

	_url := "https://api.weixin.qq.com/wxa/generatescheme?access_token=" + url.QueryEscape(accessToken)
	httpResp, err := httpClient.Post(_url, "application/json; charset=utf-8", bf)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http.Status: %s", httpResp.Status)
		return
	}

	var result generateSchemeResp
	if err = api.DecodeJSONHttpResponse(httpResp.Body, &result); err != nil {
		return
	}

	return result.Openlink, nil
}
