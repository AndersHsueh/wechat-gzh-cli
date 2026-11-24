package client

import (
	"context"
	"net/http"
	"net/url"
)

// UserInfo 用户基本信息(与公众号文档字段一致)。
type UserInfo struct {
	Subscribe      int     `json:"subscribe"`
	OpenID         string  `json:"openid"`
	Nickname       string  `json:"nickname"`
	Sex            int     `json:"sex"`
	City           string  `json:"city"`
	Province       string  `json:"province"`
	Country        string  `json:"country"`
	HeadImgURL     string  `json:"headimgurl"`
	SubscribeTime  int64   `json:"subscribe_time"`
	UnionID        string  `json:"unionid"`
	Remark         string  `json:"remark"`
	GroupID        int64   `json:"groupid"`
	TagIDList      []int64 `json:"tagid_list"`
	SubscribeScene string  `json:"subscribe_scene"`
}

// FollowerResult 关注用户列表返回。
type FollowerResult struct {
	Total int `json:"total"`
	Count int `json:"count"`
	Data  struct {
		OpenID []string `json:"openid"`
	} `json:"data"`
	NextOpenID string `json:"next_openid"`
}

// UserInfo 根据 OpenID 获取用户基本信息(GET /cgi-bin/user/info)。
func (c *Client) GetUserInfo(ctx context.Context, openid, lang string) (*UserInfo, error) {
	if lang == "" {
		lang = "zh_CN"
	}
	q := url.Values{}
	q.Set("openid", openid)
	q.Set("lang", lang)
	var resp UserInfo
	err := c.doJSON(ctx, http.MethodGet, "/cgi-bin/user/info", q, nil, &resp)
	return &resp, err
}

// GetFollowers 获取关注用户列表(GET /cgi-bin/user/get)。
func (c *Client) GetFollowers(ctx context.Context, nextOpenID string) (*FollowerResult, error) {
	q := url.Values{}
	if nextOpenID != "" {
		q.Set("next_openid", nextOpenID)
	}
	var resp FollowerResult
	err := c.doJSON(ctx, http.MethodGet, "/cgi-bin/user/get", q, nil, &resp)
	return &resp, err
}
