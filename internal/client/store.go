package client

import (
	"context"
	"net/http"
)

// WxaEndpoint 门店(小程序/POI)接口定义。
type WxaEndpoint struct {
	// Name 子命令名。
	Name string
	// Path wxa 接口路径(不含 /wxa 前缀)。
	Path string
	// Summary 说明。
	Summary string
	// NeedBody 是否需要请求体(申请书类接口为 true)。
	NeedBody bool
	// BodyHint 请求体 JSON 示例说明,空则为空 body。
	BodyHint string
}

// WxaEndpoints 全部门店小程序与门店 POI 接口。
var WxaEndpoints = []WxaEndpoint{
	// 门店小程序
	{Name: "shop-category", Path: "/wxa/get_merchant_category", Summary: "拉取门店小程序类目", NeedBody: false},
	{Name: "shop-apply", Path: "/wxa/apply_merchant", Summary: "创建门店小程序(提交入驻申请)", NeedBody: true, BodyHint: "申请材料 JSON"},
	{Name: "shop-audit", Path: "/wxa/get_merchant_audit_info", Summary: "获取门店小程序审核结果", NeedBody: false},
	{Name: "shop-modify", Path: "/wxa/modify_merchant", Summary: "修改门店小程序信息", NeedBody: true, BodyHint: "待修改字段 JSON"},
	{Name: "shop-district", Path: "/wxa/get_district", Summary: "获取省市区信息", NeedBody: false},
	{Name: "shop-map-poi", Path: "/wxa/search_map_poi", Summary: "搜索门店地图信息", NeedBody: true, BodyHint: "搜索条件 JSON"},
	// 门店 POI
	{Name: "store-add", Path: "/wxa/add_store", Summary: "新增线下实体门店", NeedBody: true, BodyHint: "门店信息 JSON"},
	{Name: "store-get", Path: "/wxa/get_store_info", Summary: "获取门店详情", NeedBody: true, BodyHint: "含 poi_id"},
	{Name: "store-list", Path: "/wxa/get_store_list", Summary: "获取门店列表", NeedBody: true, BodyHint: "分页参数 JSON"},
	{Name: "store-del", Path: "/wxa/del_store", Summary: "删除门店", NeedBody: true, BodyHint: "含 poi_id"},
	{Name: "store-update", Path: "/wxa/update_store", Summary: "更新门店信息", NeedBody: true, BodyHint: "门店信息 JSON"},
	{Name: "store-map-create", Path: "/wxa/create_map_poi", Summary: "在地图中创建门店 POI", NeedBody: true, BodyHint: "地图 POI 信息 JSON"},
}

// WxaQuery 调用门店(wxa)接口。body 可为 nil(仅当接口不需请求体)。
func (c *Client) WxaQuery(ctx context.Context, path string, body map[string]any) (map[string]any, error) {
	var in map[string]any
	if len(body) > 0 {
		in = body
	}
	var resp map[string]any
	err := c.doJSON(ctx, http.MethodPost, path, nil, in, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
