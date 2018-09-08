package applet

import (
	"github.com/BiteBit/applet/src/api"
	"github.com/BiteBit/applet/src/crypto"
	"github.com/BiteBit/applet/src/message"
)

type (
	// Applet 小程序
	Applet struct {
		appID          string                   // 小程序id
		appKey         string                   // 小程序秘钥
		appToken       string                   // 小程序token
		encodingAESKey string                   // 小程序加密秘钥
		API            *api.WechatAPI           // 小程序接口
		Token          *api.WechatToken         // 小程序token服务
		Crypto         *crypto.WechatCrypto     // 微信加密解密工具
		Messager       *message.WechatMessenger // 小程序消息信使
	}
)

// NewApplet 新建一个小程序实例
func NewApplet(appID, appKey, appToken, encodingAESKey string, tokenStore api.WechatTokenStore) *Applet {
	token := api.NewWechatToken(tokenStore)
	apix := api.NewWechatAPI(appID, appKey, token)
	crypto := crypto.NewWechatCrypto(appID, appToken, encodingAESKey)
	messager := message.NewWechatMessager(crypto)

	return &Applet{
		appID:          appID,
		appKey:         appKey,
		appToken:       appToken,
		encodingAESKey: encodingAESKey,
		API:            apix,
		Token:          token,
		Crypto:         crypto,
		Messager:       messager,
	}
}
