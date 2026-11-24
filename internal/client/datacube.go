package client

import (
	"context"
	"net/http"
)

// DatacubePath 定义 datacube 接口的请求路径、命令名与说明。
type DatacubePath struct {
	// Name 命令子命令名。
	Name string
	// Path datacube 接口路径(不含 /datacube 前缀)。
	Path string
	// Summary 一句话说明。
	Summary string
}

// DatacubeEndpoints 全部 datacube 接口。
var DatacubeEndpoints = []DatacubePath{
	// 用户数据
	{Name: "user-summary", Path: "/datacube/getusersummary", Summary: "获取用户增减数据(每日新增/取消关注,分渠道)"},
	{Name: "user-cumulate", Path: "/datacube/getusercumulate", Summary: "获取累计用户数据(截止每日累计关注)"},
	// 图文数据(新系列,建议使用)
	{Name: "article-read", Path: "/datacube/getarticleread", Summary: "发表内容每日阅读数据"},
	{Name: "article-share", Path: "/datacube/getarticleshare", Summary: "发表内容每日分享数据"},
	{Name: "biz-summary", Path: "/datacube/getbizsummary", Summary: "发表内容概况总数据"},
	{Name: "article-total-detail", Path: "/datacube/getarticletotaldetail", Summary: "发表内容发表详细数据"},
	// 图文数据(旧系列,官方已标注停止维护,兼容保留)
	{Name: "article-summary", Path: "/datacube/getarticlesummary", Summary: "图文群发每日数据(已停止维护)"},
	{Name: "article-total", Path: "/datacube/getarticletotal", Summary: "图文群发总数据(已停止维护)"},
	{Name: "user-read", Path: "/datacube/getuserread", Summary: "图文阅读概况数据(已停止维护)"},
	{Name: "user-read-hour", Path: "/datacube/getuserreadhour", Summary: "图文阅读分时数据(已停止维护)"},
	{Name: "user-share", Path: "/datacube/getusershare", Summary: "图文转发概况数据(已停止维护)"},
	{Name: "user-share-hour", Path: "/datacube/getusersharehour", Summary: "图文转发分时数据(已停止维护)"},
	// 消息数据
	{Name: "msg", Path: "/datacube/getupstreammsg", Summary: "消息发送概况数据"},
	{Name: "msg-week", Path: "/datacube/getupstreammsgweek", Summary: "消息发送周数据"},
	{Name: "msg-month", Path: "/datacube/getupstreammsgmonth", Summary: "消息发送月数据"},
	{Name: "msg-dist", Path: "/datacube/getupstreammsgdist", Summary: "消息发送分布数据"},
	{Name: "msg-dist-week", Path: "/datacube/getupstreammsgdistweek", Summary: "消息发送分布周数据"},
	{Name: "msg-dist-month", Path: "/datacube/getupstreammsgdistmonth", Summary: "消息发送分布月数据"},
	{Name: "msg-hour", Path: "/datacube/getupstreammsghour", Summary: "消息发送分时数据"},
	// 接口数据
	{Name: "interface-summary", Path: "/datacube/getinterfacesummary", Summary: "被动回复概要数据"},
	{Name: "interface-summary-hour", Path: "/datacube/getinterfacesummaryhour", Summary: "被动回复分布数据"},
}

// DatacubeQuery 查询指定 datacube 接口。
// begin/end 格式 YYYY-MM-DD。
func (c *Client) DatacubeQuery(ctx context.Context, path, begin, end string) ([]map[string]any, error) {
	var resp struct {
		List []map[string]any `json:"list"`
	}
	err := c.doJSON(ctx, http.MethodPost, path, nil,
		map[string]any{"begin_date": begin, "end_date": end}, &resp)
	if err != nil {
		return nil, err
	}
	return resp.List, nil
}
