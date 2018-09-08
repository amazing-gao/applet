package api

type (
	// CustomerMsg 客服消息主体
	CustomerMsg struct {
		OpenID  string             `json:"touser"`
		MsgType string             `json:"msgtype"` // text,image,link,miniprogrampage
		Text    *CustomerMsgText   `json:"text,omitempty"`
		Image   *CustomerMsgImage  `json:"image,omitempty"`
		Link    *CustomerMsgLink   `json:"link,omitempty"`
		Applet  *CustomerMsgApplet `json:"miniprogrampage,omitempty"`
	}

	// 文本
	CustomerMsgText struct {
		Content string `json:"content"`
	}

	// 图片
	CustomerMsgImage struct {
		MediaID string `json:"media_id"`
	}

	// 图文链接
	CustomerMsgLink struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		URL         string `json:"url"`
		ThumbURL    string `json:"thumb_url"`
	}

	// 小程序卡片
	CustomerMsgApplet struct {
		Title        string `json:"title"`
		PagePath     string `json:"pagepath"`
		ThumbMediaID string `json:"thumb_media_id"`
	}
)

// SendMessage 发送客服消息
func (api *WechatAPI) SendMessage(msg *CustomerMsg) (*WechatResp, []error) {
	return api.Request(&option{
		method:    "POST",
		url:       "/cgi-bin/message/custom/send",
		withToken: true,
		body:      msg,
	})
}

// SendText 发送文本消息
func (api *WechatAPI) SendText(openid, content string) (*WechatResp, []error) {
	return api.SendMessage(&CustomerMsg{
		OpenID:  openid,
		MsgType: "text",
		Text: &CustomerMsgText{
			Content: content,
		},
	})
}

// SendImage 发送图片消息
func (api *WechatAPI) SendImage(openid, mediaId string) (*WechatResp, []error) {
	return api.SendMessage(&CustomerMsg{
		OpenID:  openid,
		MsgType: "image",
		Image: &CustomerMsgImage{
			MediaID: mediaId,
		},
	})
}

// SendLink 发送图文消息
func (api *WechatAPI) SendLink(openid string, link *CustomerMsgLink) (*WechatResp, []error) {
	return api.SendMessage(&CustomerMsg{
		OpenID:  openid,
		MsgType: "link",
		Link:    link,
	})
}

// SendApplet 发送小程序消息
func (api *WechatAPI) SendApplet(openid string, applet *CustomerMsgApplet) (*WechatResp, []error) {
	return api.SendMessage(&CustomerMsg{
		OpenID:  openid,
		MsgType: "miniprogrampage",
		Applet:  applet,
	})
}
