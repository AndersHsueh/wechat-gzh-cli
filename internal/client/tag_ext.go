package client

import (
	"context"
	"net/http"
)

// TagFansResult 标签下粉丝列表返回。
type TagFansResult struct {
	Count  int    `json:"count"`
	Data   struct {
		OpenID []string `json:"openid"`
	} `json:"data"`
	NextOpenID string `json:"next_openid"`
}

// TagGetFans 获取标签下粉丝列表(POST /cgi-bin/user/tag/get)。
func (c *Client) TagGetFans(ctx context.Context, tagID int, nextOpenID string) (*TagFansResult, error) {
	body := map[string]any{"tagid": tagID}
	if nextOpenID != "" {
		body["next_openid"] = nextOpenID
	}
	var resp TagFansResult
	err := c.doJSON(ctx, http.MethodPost, "/cgi-bin/user/tag/get", nil, body, &resp)
	return &resp, err
}

// TagUpdate 编辑标签(POST /cgi-bin/tags/update)。
func (c *Client) TagUpdate(ctx context.Context, tagID int, name string) error {
	return c.doJSON(ctx, http.MethodPost, "/cgi-bin/tags/update", nil,
		map[string]any{"tag": map[string]any{"id": tagID, "name": name}}, nil)
}

// TagDelete 删除标签(POST /cgi-bin/tags/delete)。
func (c *Client) TagDelete(ctx context.Context, tagID int) error {
	return c.doJSON(ctx, http.MethodPost, "/cgi-bin/tags/delete", nil, map[string]any{"tag": map[string]any{"id": tagID}}, nil)
}

// BatchTagging 批量为用户打标签(POST /cgi-bin/tags/members/batchtagging)。
func (c *Client) BatchTagging(ctx context.Context, openids []string, tagID int) error {
	return c.doJSON(ctx, http.MethodPost, "/cgi-bin/tags/members/batchtagging", nil,
		map[string]any{"openid_list": openids, "tagid": tagID}, nil)
}

// BatchUntagging 批量为用户取消标签(POST /cgi-bin/tags/members/batchuntagging)。
func (c *Client) BatchUntagging(ctx context.Context, openids []string, tagID int) error {
	return c.doJSON(ctx, http.MethodPost, "/cgi-bin/tags/members/batchuntagging", nil,
		map[string]any{"openid_list": openids, "tagid": tagID}, nil)
}

// GetUserTagList 获取用户的标签列表(POST /cgi-bin/tags/getidlist)。
func (c *Client) GetUserTagList(ctx context.Context, openid string) ([]int, error) {
	var resp struct {
		TagIDList []int `json:"tagid_list"`
	}
	err := c.doJSON(ctx, http.MethodPost, "/cgi-bin/tags/getidlist", nil, map[string]any{"openid": openid}, &resp)
	return resp.TagIDList, err
}
