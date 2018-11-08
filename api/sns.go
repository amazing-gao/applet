package api

type (
	// RespCode2Session 响应结果
	RespCode2Session struct {
		UnionID    string `json:"unionid"`
		OpenID     string `json:"openid"`
		SessionKey string `json:"session_key"`
		ExpiresIn  int    `json:"expires_in"`
	}
)

// Code2Session oauth code转换为unionId,openId,sessionKey
func (api *WechatAPI) Code2Session(code string) (*RespCode2Session, *WechatResp, []error) {
	respData := &RespCode2Session{}
	resp, errs := api.Request(&option{
		method:    "GET",
		url:       "/sns/jscode2session",
		withToken: false,
		query: map[string]string{
			"appid":      api.appID,
			"secret":     api.appKey,
			"js_code":    code,
			"grant_type": "authorization_code",
		},
	}, &respData)

	return respData, resp, errs
}
