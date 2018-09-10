package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/parnurzeal/gorequest"
)

type (
	// WechatAPI 小程序API
	WechatAPI struct {
		appID            string
		appKey           string
		apiScheme        string
		apiDomain        string
		apiToken         string
		apiTokenExpireIn int
		apiTokenExpireAt time.Time
		apiTokenStore    WechatTokenStore
	}

	// WechatTokenStore 微信token存储器
	WechatTokenStore interface {
		Get() (token string, exipresIn int, expireAt time.Time)
		Set(token string, exipresIn int, expireAt time.Time)
	}

	// WechatResp 微信接口响应
	WechatResp struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
		Data    []byte
	}

	// WechatRespToken 微信token响应结果
	WechatRespToken struct {
		Token     string `json:"access_token"`
		ExpiresIn int    `json:"expires_in"`
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
func NewWechatAPI(appID, appKey string, tokenStore WechatTokenStore) *WechatAPI {
	return &WechatAPI{
		apiScheme:     "https",
		apiDomain:     "api.weixin.qq.com",
		appID:         appID,
		appKey:        appKey,
		apiTokenStore: tokenStore,
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
		req.Query(fmt.Sprintf("access_token=%s", api.GetToken()))
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
	wechatResp.Data = []byte(body)
	return wechatResp, nil
}

// GetToken 获取token
// 优先 内存获取
// 再次 store获取
// 再次 生成并保存
func (api *WechatAPI) GetToken() string {
	if api.IsValid() {
		return api.apiToken
	}

	token, exipresIn, expireAt := api.apiTokenStore.Get()
	if isValid(token, expireAt) {
		api.apiToken = token
		api.apiTokenExpireIn = exipresIn
		api.apiTokenExpireAt = expireAt

		return token
	}

	resp, err := api.RenewToken()
	fmt.Println("RenewToken", err, resp)

	return api.apiToken
}

// RenewToken 重新生成一个token
func (api *WechatAPI) RenewToken() (*WechatResp, []error) {
	query := map[string]string{}
	query["grant_type"] = "client_credential"
	query["appid"] = api.appID
	query["secret"] = api.appKey

	resp, errs := api.Request(&option{
		method:    "GET",
		url:       "/cgi-bin/token",
		withToken: false,
		query:     query,
	})

	// 请求成功，解析内容
	if resp.ErrCode == 0 {
		respToken := &WechatRespToken{}

		err := json.Unmarshal(resp.Data, respToken)
		if err != nil {
			return resp, []error{err}
		}

		api.apiToken = respToken.Token
		api.apiTokenExpireIn = respToken.ExpiresIn
		api.apiTokenExpireAt = time.Now().Add(time.Second * time.Duration(respToken.ExpiresIn-30)) // 提前30秒失效token

		// 保存新的token
		api.apiTokenStore.Set(api.apiToken, api.apiTokenExpireIn, api.apiTokenExpireAt)
	}

	return resp, errs
}

// IsValid 返回当前的token是否过期
// token值存在并且没有到达有效期
func (api *WechatAPI) IsValid() bool {
	return isValid(api.apiToken, api.apiTokenExpireAt)
}

func isValid(apiToken string, apiTokenExpireAt time.Time) bool {
	return (apiToken != "" && apiTokenExpireAt.After(time.Now()))
}
