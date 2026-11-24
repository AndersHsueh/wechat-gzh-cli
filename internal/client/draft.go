package client

import (
	"context"
	"net/http"
)

// DraftItem 草稿列表项。
type DraftItem struct {
	MediaID    string `json:"media_id"`
	UpdateTime int64  `json:"update_time"`
}

// DraftBatchResult 草稿列表返回。
type DraftBatchResult struct {
	TotalCount int         `json:"total_count"`
	ItemCount  int         `json:"item_count"`
	Item       []DraftItem `json:"item"`
}

// DraftAdd 新增草稿(POST /cgi-bin/draft/add)。
func (c *Client) DraftAdd(ctx context.Context, articles []map[string]any) (string, error) {
	var resp struct {
		MediaID string `json:"media_id"`
	}
	err := c.doJSON(ctx, http.MethodPost, "/cgi-bin/draft/add", nil, map[string]any{"articles": articles}, &resp)
	return resp.MediaID, err
}

// DraftUpdate 更新草稿(POST /cgi-bin/draft/update)。
func (c *Client) DraftUpdate(ctx context.Context, mediaID string, index int, article map[string]any) error {
	body := map[string]any{"media_id": mediaID, "index": index, "articles": article}
	return c.doJSON(ctx, http.MethodPost, "/cgi-bin/draft/update", nil, body, nil)
}

// DraftBatchGet 获取草稿列表(POST /cgi-bin/draft/batchget)。
func (c *Client) DraftBatchGet(ctx context.Context, offset, count int) (*DraftBatchResult, error) {
	var resp DraftBatchResult
	err := c.doJSON(ctx, http.MethodPost, "/cgi-bin/draft/batchget", nil,
		map[string]any{"offset": offset, "count": count}, &resp)
	return &resp, err
}

// DraftCount 获取草稿总数(POST /cgi-bin/draft/count)。
func (c *Client) DraftCount(ctx context.Context) (int, error) {
	var resp struct {
		TotalCount int `json:"total_count"`
	}
	err := c.doJSON(ctx, http.MethodPost, "/cgi-bin/draft/count", nil, map[string]any{}, &resp)
	return resp.TotalCount, err
}

// DraftDelete 删除草稿(POST /cgi-bin/draft/delete)。
func (c *Client) DraftDelete(ctx context.Context, mediaID string) error {
	return c.doJSON(ctx, http.MethodPost, "/cgi-bin/draft/delete", nil, map[string]any{"media_id": mediaID}, nil)
}

// DraftGet 获取草稿详情(POST /cgi-bin/draft/get)。mediaID 传 "" 表示不取 content。
func (c *Client) DraftGet(ctx context.Context, mediaID string) ([]map[string]any, error) {
	body := map[string]any{}
	if mediaID != "" {
		body["media_id"] = mediaID
	}
	var resp struct {
		NewsItem []map[string]any `json:"news_item"`
	}
	err := c.doJSON(ctx, http.MethodPost, "/cgi-bin/draft/get", nil, body, &resp)
	return resp.NewsItem, err
}
