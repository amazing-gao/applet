package message

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/amazing-gao/applet/crypto"
)

type (
	// WechatMessenger 小程序消息推送信使
	WechatMessenger struct {
		crypto         *crypto.WechatCrypto
		messageHandler Handler // 客服消息事件处理器
	}

	// Message 小程序消息推送
	Message struct {
		XMLName      xml.Name `json:"-" xml:"xml"`
		MsgID        int      `json:"MsgId" xml:"MsgId"`               // 消息id，64位整型
		MsgType      string   `json:"MsgType" xml:"MsgType"`           // text image miniprogrampage event
		EncryptMsg   string   `json:"Encrypt" xml:"Encrypt"`           // 加密后的消息
		ToUserName   string   `json:"ToUserName" xml:"ToUserName"`     // 小程序的原始ID
		FromUserName string   `json:"FromUserName" xml:"FromUserName"` // 发送者的openid
		CreateTime   int      `json:"CreateTime" xml:"CreateTime"`     // 事件创建时间(整型）
		Content      string   `json:"Content" xml:"Content"`           // text: 文本消息内容
		PicURL       string   `json:"PicUrl" xml:"PicUrl"`             // image: 图片链接（由系统生成）
		MediaID      string   `json:"MediaId" xml:"MediaId"`           // image: 图片消息媒体id，可以调用获取临时素材接口拉取数据。
		Title        string   `json:"Title" xml:"Title"`               // miniprogrampage: 标题
		AppID        string   `json:"AppId" xml:"AppId"`               // miniprogrampage: 小程序appid
		PagePath     string   `json:"PagePath" xml:"PagePath"`         // miniprogrampage: 小程序页面路径
		ThumbURL     string   `json:"ThumbUrl" xml:"ThumbUrl"`         // miniprogrampage: 封面图片的临时cdn链接
		ThumbMediaID string   `json:"ThumbMediaId" xml:"ThumbMediaId"` // miniprogrampage: 封面图片的临时素材id
		Event        string   `json:"Event" xml:"Event"`               // event: 事件类型，user_enter_tempsession
		SessionFrom  string   `json:"SessionFrom" xml:"SessionFrom"`   // event: 开发者在客服会话按钮设置的session-from属性
		Query        string   `json:"Query" xml:"Query"`               // 搜索内容
		Scene        int      `json:"Scene" xml:"Scene"`               // 场景值
	}

	// Handler 小程序消息推送处理器
	Handler func(*Message, error) string
)

// NewWechatMessager 新建一个微信消息信使
func NewWechatMessager(crypto *crypto.WechatCrypto) *WechatMessenger {
	return &WechatMessenger{
		crypto: crypto,
	}
}

// RegisterHandler 注册小程序消息推送处理器
func (mgr *WechatMessenger) RegisterHandler(messageHandler Handler) *WechatMessenger {
	mgr.messageHandler = messageHandler

	return mgr
}

// MessageHandleMiddleware 小程序消息处理中间件
// GET 验证消息的确来自微信服务器
// POST 处理客服消息
func (mgr *WechatMessenger) MessageHandleMiddleware(request *http.Request, writer http.ResponseWriter) {
	if request.Method == "GET" {
		mgr.MessageHandleValid(request, writer)
	} else if request.Method == "POST" {
		mgr.MessageHandle(request, writer)
	} else {
		mgr.MessageHandleNotSupport(request, writer)
	}
}

// MessageHandleValid 消息校验
func (mgr *WechatMessenger) MessageHandleValid(request *http.Request, writer http.ResponseWriter) {
	querys := request.URL.Query()
	nonce := querys.Get("nonce")
	echostr := querys.Get("echostr")
	signature := querys.Get("signature")
	timestamp := querys.Get("timestamp")

	if mgr.crypto.CheckSignature(timestamp, nonce, signature) {
		writer.Write([]byte(echostr))
	} else {
		writer.Write([]byte("invalid"))
	}
}

// MessageHandle 处理小程序消息推送
//                    -> 校验密文 -> 解析密文 \
//                  /                       \
// 解析明文 -----> 有密文 -----------------> 处理消息 ----> 响应腾讯服务器
//
func (mgr *WechatMessenger) MessageHandle(request *http.Request, writer http.ResponseWriter) {
	contentType := request.Header.Get("Content-Type")
	querys := request.URL.Query()
	nonce := querys.Get("nonce")
	timestamp := querys.Get("timestamp")
	encryptType := querys.Get("encrypt_type")
	msgSignature := querys.Get("msg_signature")

	var err error
	var ret string
	for index := 0; index < 1; index++ {
		msg := &Message{}

		rawMsg := []byte{}
		rawMsg, err = ioutil.ReadAll(request.Body)
		if err != nil {
			break
		}

		// 解析明文消息
		err = unmarshal(contentType, rawMsg, msg)
		if err != nil {
			break
		}

		// 如果消息已加密，先校验消息是否合法。如果合法，再解密
		if encryptType != "" && msgSignature != "" {
			// 如果校验失败，那么是非法消息
			if !mgr.crypto.CheckMsgSignature(timestamp, nonce, msg.EncryptMsg, msgSignature) {
				err = fmt.Errorf("invalid message")
				break
			}

			// 解密密文消息
			encryptMsg := mgr.crypto.Decrypt(msg.EncryptMsg)
			err = unmarshal(contentType, []byte(encryptMsg), msg)
			if err != nil {
				break
			}
		}

		// 处理消息
		ret = mgr.messageHandler(msg, err)
	}

	if err != nil {
		log.Println("Applet.MessageHandle.Error", err)
		ret = err.Error()
	} else if len(ret) == 0 {
		ret = "success"
	}

	writer.Write([]byte(ret))
}

// MessageHandleNotSupport 不支持的消息
func (mgr *WechatMessenger) MessageHandleNotSupport(request *http.Request, writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusMethodNotAllowed)
	writer.Write([]byte("Method Not Allowed"))
}
