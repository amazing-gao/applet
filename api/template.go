package api

type (
	// RespTemplate 模版内容
	RespTemplate struct {
		TemplateID string `json:"template_id"`
		Title      string `json:"title"`
		Content    string `json:"content"`
		Example    string `json:"example"`
	}

	// RespTemplateList 模版列表
	RespTemplateList []RespTemplate
)

// GetTemplateList 获取帐号下已存在的模板列表
func (api *WechatAPI) GetTemplateList(offset, count uint) (RespTemplateList, *WechatResp, []error) {
	respData := RespTemplateList{}
	resp, errs := api.Request(&option{
		method:    "GET",
		url:       "/cgi-bin/wxopen/template/list",
		withToken: true,
		query: map[string]uint{
			"offset": offset,
			"count":  count,
		},
	}, &respData)

	return respData, resp, errs
}
