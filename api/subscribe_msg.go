package api

type (
	// SubscribeMsg 小程序订阅消息
	SubscribeMsg struct {
		Touser     string      `json:"touser"`
		TemplateID string      `json:"template_id"`
		Page       string      `json:"page"`
		Data       interface{} `json:"data"`
	}
)

// SendSubscribeMessage 下发小程序订阅消息
func (api *WechatAPI) SendSubscribeMessage(msg *SubscribeMsg) (*WechatResp, []error) {
	return api.Request(&option{
		method:    "POST",
		url:       "/cgi-bin/message/subscribe/send",
		withToken: true,
		body:      msg,
	})
}
