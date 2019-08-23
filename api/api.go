package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"sync"

	"github.com/parnurzeal/gorequest"
)

type (
	// WechatAPI 小程序API
	WechatAPI struct {
		appID         string
		appKey        string
		apiScheme     string
		apiDomain     string
		apiBasePath   string
		apiTokenStore WechatTokenStore
		before        Before
		after         After
		locker        *sync.Mutex
	}

	// WechatResp 微信接口响应
	WechatResp struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}

	// Before request
	Before func(*gorequest.SuperAgent)

	// After request
	After func(*gorequest.SuperAgent, []error, string, *gorequest.Response)

	option struct {
		method    string
		url       string
		withToken bool
		query     interface{}
		body      interface{}
	}
)

// NewWechatAPI 生成一个api
func NewWechatAPI(appID, appKey string, tokenStore WechatTokenStore) *WechatAPI {
	return &WechatAPI{
		apiScheme:     "https",
		apiDomain:     "api.weixin.qq.com",
		apiBasePath:   "/",
		appID:         appID,
		appKey:        appKey,
		apiTokenStore: tokenStore,
		locker:        &sync.Mutex{},
	}
}

// SetBefore 设置请求前hook
func (api *WechatAPI) SetBefore(be Before) {
	api.before = be
}

// SetAftre 设置请求后hook
func (api *WechatAPI) SetAftre(af After) {
	api.after = af
}

// SetDomain 设置domain，覆盖默认的api.weixin.qq.com
func (api *WechatAPI) SetDomain(domain string) {
	api.apiDomain = domain
}

// SetBasePath 设置请求的BasePath，覆盖默认的/
func (api *WechatAPI) SetBasePath(basePath string) {
	api.apiBasePath = basePath
}

// Request 请求
func (api *WechatAPI) Request(opt *option, respData ...interface{}) (*WechatResp, []error) {
	req := gorequest.New()

	u := &url.URL{
		Scheme: api.apiScheme,
		Host:   api.apiDomain,
		Path:   path.Join(api.apiBasePath, opt.url),
	}

	if opt.method == "GET" {
		req.Get(u.String())
	} else if opt.method == "POST" {
		req.Post(u.String())
	}

	req.Query(opt.query).Send(opt.body)

	if opt.withToken {
		req.Query(fmt.Sprintf("access_token=%s", api.GetToken()))
	}

	if api.before != nil {
		api.before(req)
	}

	resp, body, errs := req.End()

	if api.after != nil {
		api.after(req, errs, body, &resp)
	}

	if len(errs) != 0 {
		return nil, errs
	}

	wechatResp := &WechatResp{}
	if err := json.Unmarshal([]byte(body), wechatResp); err != nil {
		return nil, []error{err}
	}

	// 业务错误，直接返回
	if wechatResp.ErrCode != 0 || wechatResp.ErrMsg != "" {
		return wechatResp, nil
	}

	// 业务成功，获取返回的数据
	if len(respData) != 0 {
		if err := json.Unmarshal([]byte(body), respData[0]); err != nil {
			return wechatResp, []error{err}
		}
	}

	return wechatResp, nil
}
