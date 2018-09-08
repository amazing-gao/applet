package api

import (
	"time"
)

type (
	// WechatToken 微信API Token
	WechatToken struct {
		token    string
		expireIn time.Duration
		expireAt time.Time
		store    WechatTokenStore
	}

	// WechatTokenStore 微信token存储器
	WechatTokenStore interface {
		Get() (token string, exipresIn time.Duration, expireAt time.Time)
		Set(token string, exipresIn time.Duration, expireAt time.Time)
		Renew() (token string, exipresIn time.Duration, expireAt time.Time)
		IsUpdating() bool
	}
)

// NewWechatToken 新建微信token
func NewWechatToken(store WechatTokenStore) *WechatToken {
	return &WechatToken{
		store: store,
	}
}

// GetToken 获取一个有效的token值
func (token *WechatToken) GetToken() string {
	if token.IsValid() || token.store.IsUpdating() {
		return token.token
	}

	// 获取存储器中最新的token
	newToken, newExpireIn, newExpireAt := token.store.Get()

	// 如果存储器中的也无效，那么重新申请一个
	if !token.isValid(newToken, newExpireAt) {
		newToken, newExpireIn, newExpireAt = token.store.Renew()

		// 提前一分钟进行token的刷新
		newExpireAt = newExpireAt.Add(time.Minute)

		// 保存新刷新的token
		token.store.Set(newToken, newExpireIn, newExpireAt)
	}

	// 缓存到内存中
	token.token = newToken
	token.expireIn = newExpireIn
	token.expireAt = newExpireAt

	return newToken
}

// IsValid 返回当前的token是否过期
// token值存在并且没有到达有效期
func (token *WechatToken) IsValid() bool {
	return token.isValid(token.token, token.expireAt)
}

func (token *WechatToken) isValid(tokenValue string, expireAt time.Time) bool {
	return (tokenValue != "" && expireAt.After(time.Now()))
}
