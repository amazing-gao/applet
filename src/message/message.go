package message

import (
	"encoding/json"

	"github.com/BiteBit/applet/src/crypto"
)

type (
	// WechatMessenger 小程序消息推送信使
	WechatMessenger struct {
		crypto         *crypto.WechatCrypto
		messageHandler func(*Message) // 客服消息事件处理器
	}

	// Message 小程序消息推送
	Message struct {
		MsgID        int    `json:"MsgId"`        // 消息id，64位整型
		MsgType      string `json:"MsgType"`      // text image miniprogrampage event
		EncryptMsg   string `json:"Encrypt"`      // 加密后的消息
		ToUserName   string `json:"ToUserName"`   // 小程序的原始ID
		FromUserName string `json:"FromUserName"` // 发送者的openid
		CreateTime   int    `json:"CreateTime"`   // 事件创建时间(整型）
		Content      string `json:"Content"`      // text: 文本消息内容
		PicURL       string `json:"PicUrl"`       // image: 图片链接（由系统生成）
		MediaID      string `json:"MediaId"`      // image: 图片消息媒体id，可以调用获取临时素材接口拉取数据。
		Title        string `json:"Title"`        // miniprogrampage: 标题
		AppID        string `json:"AppId"`        // miniprogrampage: 小程序appid
		PagePath     string `json:"PagePath"`     // miniprogrampage: 小程序页面路径
		ThumbURL     string `json:"ThumbUrl"`     // miniprogrampage: 封面图片的临时cdn链接
		ThumbMediaID string `json:"ThumbMediaId"` // miniprogrampage: 封面图片的临时素材id
		Event        string `json:"Event"`        // event: 事件类型，user_enter_tempsession
		SessionFrom  string `json:"SessionFrom"`  // event: 开发者在客服会话按钮设置的session-from属性
	}

	// Handler 小程序消息推送处理器
	// Handler
)

// NewWechatMessager 新建一个微信消息信使
func NewWechatMessager(crypto *crypto.WechatCrypto) *WechatMessenger {
	return &WechatMessenger{
		crypto: crypto,
	}
}

// RegisterHandler 注册小程序消息推送处理器
func (mgr *WechatMessenger) RegisterHandler(messageHandler func(*Message)) *WechatMessenger {
	mgr.messageHandler = messageHandler

	return mgr
}

// MessageHandle 处理小程序消息推送
//                    -> 校验密文 -> 解析密文 \
//                  /                       \
// 解析明文 -----> 有密文 -----------------> 处理消息 ----> 响应腾讯服务器
//
func (mgr *WechatMessenger) MessageHandle(rawMsg []byte) *WechatMessenger {
	msg := &Message{}

	// 解析消息
	err := mgr.messageParser(rawMsg, msg)
	if err != nil {
		return mgr
	}

	// 处理消息
	mgr.messageHandler(msg)

	return mgr
}

// messageParser 小程序消息推送解析
func (mgr *WechatMessenger) messageParser(rawMsg []byte, msg *Message) (err error) {
	msg = &Message{}

	// 解析消息
	if err := json.Unmarshal(rawMsg, msg); err != nil {
		return err
	}

	// 如果消息加密，那么解密消息
	if msg.EncryptMsg != "" {
		decryptMsg := mgr.crypto.Decrypt(msg.EncryptMsg)

		if err := json.Unmarshal(([]byte)(decryptMsg), msg); err != nil {
			return err
		}
	}

	return nil
}
