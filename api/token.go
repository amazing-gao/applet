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
	token, exipresIn, expireAt := api.apiTokenStore.Get()
	log.Printf("Applet.GetToken appID:%s exipresIn:%d expireAt:%s", api.appID, exipresIn, expireAt)
	return token
}

// RenewToken 重新生成一个token
func (api *WechatAPI) RenewToken() (*WechatResp, []error) {
	defer api.locker.Unlock()

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
		apiToken := respToken.Token
		apiTokenExpireIn := respToken.ExpiresIn
		apiTokenExpireAt := time.Now().Add(time.Second * time.Duration(respToken.ExpiresIn-30)) // 提前30秒失效token

		// 保存新的token
		api.apiTokenStore.Set(apiToken, apiTokenExpireIn, apiTokenExpireAt)
	}

	api.locker.Lock()

	return resp, errs
}
