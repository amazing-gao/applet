package api

type (
	// UniformMsg 统一服务消息
	UniformMsg struct {
		OpenID string           `json:"touser"`
		WeApp  WeAppTemplateMsg `json:"weapp_template_msg"`
		Mp     MpTemplateMsg    `json:"mp_template_msg"`
	}

	// WeAppTemplateMsg 小程序模版消息
	WeAppTemplateMsg struct {
		TemplateID      string      `json:"template_id"`
		Page            string      `json:"page"`
		FormID          string      `json:"form_id"`
		Data            interface{} `json:"data"`
		EmphasisKeyword string      `json:"emphasis_keyword"`
	}

	// MpTemplateMsg 公众号模版消息
	MpTemplateMsg struct {
		AppID       string      `json:"appid"`
		TemplateID  string      `json:"template_id"`
		URL         string      `json:"url"`
		Miniprogram interface{} `json:"miniprogram"`
		Data        interface{} `json:"data"`
	}
)

// SendUniformMessage 下发小程序和公众号统一的服务消息
func (api *WechatAPI) SendUniformMessage(msg *UniformMsg) (*WechatResp, []error) {
	return api.Request(&option{
		method:    "POST",
		url:       "/cgi-bin/message/wxopen/template/uniform_send",
		withToken: true,
		body:      msg,
	})
}
