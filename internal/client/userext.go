package client

import (
	"context"
	"net/http"
)

// BatchUserInfo 批量获取用户基本信息(POST /cgi-bin/user/info/batchget)。
// userList 为 [{openid, lang}, ...],最多 100 条。
func (c *Client) BatchUserInfo(ctx context.Context, userList []map[string]any) ([]UserInfo, error) {
	var resp struct {
		UserInfoList []UserInfo `json:"user_info_list"`
	}
	err := c.doJSON(ctx, http.MethodPost, "/cgi-bin/user/info/batchget", nil,
		map[string]any{"user_list": userList}, &resp)
	return resp.UserInfoList, err
}

// UpdateRemark 设置用户备注名(POST /cgi-bin/user/info/updateremark)。
func (c *Client) UpdateRemark(ctx context.Context, openid, remark string) error {
	return c.doJSON(ctx, http.MethodPost, "/cgi-bin/user/info/updateremark", nil,
		map[string]any{"openid": openid, "remark": remark}, nil)
}

// GetBlacklist 获取公众号的黑名单列表(POST /cgi-bin/tags/members/getblacklist)。
// beginOpenID 为空则从首个用户开始,返回 next_openid 用于翻页。
func (c *Client) GetBlacklist(ctx context.Context, beginOpenID string) (openids []string, next string, err error) {
	body := map[string]any{}
	if beginOpenID != "" {
		body["begin_openid"] = beginOpenID
	}
	var resp struct {
		Data struct {
			OpenID []string `json:"openid"`
		} `json:"data"`
		NextOpenID string `json:"next_openid"`
	}
	err = c.doJSON(ctx, http.MethodPost, "/cgi-bin/tags/members/getblacklist", nil, body, &resp)
	return resp.Data.OpenID, resp.NextOpenID, err
}

// BatchBlacklist 拉黑用户(POST /cgi-bin/tags/members/batchblacklist)。
func (c *Client) BatchBlacklist(ctx context.Context, openids []string) error {
	return c.doJSON(ctx, http.MethodPost, "/cgi-bin/tags/members/batchblacklist", nil,
		map[string]any{"openid_list": openids}, nil)
}

// BatchUnblacklist 取消拉黑用户(POST /cgi-bin/tags/members/batchunblacklist)。
func (c *Client) BatchUnblacklist(ctx context.Context, openids []string) error {
	return c.doJSON(ctx, http.MethodPost, "/cgi-bin/tags/members/batchunblacklist", nil,
		map[string]any{"openid_list": openids}, nil)
}

// ChangeOpenIDResult openid 转换结果项。
type ChangeOpenIDResult struct {
	OpenID    string `json:"ori_openid"`
	NewOpenID string `json:"new_openid"`
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
}

// ChangeOpenID 公众号迁移后转换 openid(POST /cgi-bin/changeopenid)。
func (c *Client) ChangeOpenID(ctx context.Context, openids []string) ([]ChangeOpenIDResult, error) {
	if len(openids) > 100 {
		openids = openids[:100]
	}
	var resp struct {
		ResultList []ChangeOpenIDResult `json:"result_list"`
	}
	err := c.doJSON(ctx, http.MethodPost, "/cgi-bin/changeopenid", nil,
		map[string]any{"openid_list": openids}, &resp)
	return resp.ResultList, err
}
