package api

import (
	"log"
	"time"
)

type (
	// WechatTokenStore 微信token存储器
	WechatTokenStore interface {
		Get() (token string, exipresIn int, expireAt time.Time)
		Set(token string, exipresIn int, expireAt time.Time)
	}

	// WechatRespToken 微信token响应结果
	WechatRespToken struct {
		Token     string `json:"access_token"`
		ExpiresIn int    `json:"expires_in"`
	}
)

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
	log.Println("RenewToken", err, resp)

	return api.apiToken
}

// RenewToken 重新生成一个token
func (api *WechatAPI) RenewToken() (*WechatResp, []error) {
	respToken := &WechatRespToken{}
	resp, errs := api.Request(&option{
		method:    "GET",
		url:       "/cgi-bin/token",
		withToken: false,
		query: map[string]string{
			"grant_type": "client_credential",
			"appid":      api.appID,
			"secret":     api.appKey,
		},
	}, respToken)

	// 请求成功，解析内容
	if resp.ErrCode == 0 {
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
