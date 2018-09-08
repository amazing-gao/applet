package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/parnurzeal/gorequest"
)

type (
	// WechatAPI 小程序API
	WechatAPI struct {
		appID     string
		appKey    string
		apiScheme string
		apiDomain string
		apiToken  *WechatToken
	}

	// WechatResp 微信接口响应
	WechatResp struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
		Data    interface{}
	}

	option struct {
		method    string
		url       string
		withToken bool
		query     interface{}
		body      interface{}
	}
)

// NewWechatAPI 生成一个api
func NewWechatAPI(appID, appKey string, apiToken *WechatToken) *WechatAPI {
	return &WechatAPI{
		apiScheme: "https",
		apiDomain: "api.weixin.qq.com",
		appID:     appID,
		appKey:    appKey,
		apiToken:  apiToken,
	}
}

// Request 请求
func (api *WechatAPI) Request(opt *option) (*WechatResp, []error) {
	req := gorequest.New()

	u := &url.URL{
		Scheme: api.apiScheme,
		Host:   api.apiDomain,
		Path:   opt.url,
	}

	if opt.method == "GET" {
		req.Get(u.String())
	} else if opt.method == "POST" {
		req.Post(u.String())
	}

	req.Query(opt.query).Send(opt.body)

	if opt.withToken {
		req.Query(fmt.Sprintf("access_token=%s", api.apiToken.GetToken()))
	}

	_, body, errs := req.End()

	if len(errs) != 0 {
		return nil, errs
	}

	wechatResp := &WechatResp{}
	if err := json.Unmarshal([]byte(body), wechatResp); err != nil {
		return nil, []error{err}
	}

	// 业务错误，直接返回
	if wechatResp.ErrCode != 0 || wechatResp.ErrMsg != "" {
		return wechatResp, nil
	}

	// 业务成功
	wechatResp.Data = body
	return wechatResp, nil
}
