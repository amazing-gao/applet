package api

type (
	// WidgetImportData 抽样数据
	WidgetImportData struct {
		Lifespan uint32 `json:"lifespan"`
		Query    string `json:"query"`
		Scene    int32  `json:"scene"`
		Data     string `json:"data"`
	}

	// WidgetImportDataData 抽样数据
	WidgetImportDataData struct {
		Items     []WidgetImportDataItem    `json:"items"`
		Attribute WidgetImportDataAttribute `json:"attribute"`
	}

	// WidgetImportDataItem 抽样数据内容
	WidgetImportDataItem struct{}

	// WidgetImportDataAttribute 抽样数据属性
	WidgetImportDataAttribute struct {
		Count      int `json:"count"`
		TotalCount int `json:"totalcount"`
		ID         int `json:"id"`
		Seq        int `json:"seq"`
	}
)

// WidgetSetDynamicData 微信搜一搜，自定义模版导入抽样数据
func (api *WechatAPI) WidgetSetDynamicData(data *WidgetImportData) (*WechatResp, []error) {
	return api.Request(&option{
		method:    "POST",
		url:       "/wxa/setdynamicdata",
		withToken: true,
		body:      data,
	})
}
