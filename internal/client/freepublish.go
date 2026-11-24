package client

import (
	"context"
	"net/http"
)

// PublishSubmit 发布草稿(POST /cgi-bin/freepublish/submit)。
func (c *Client) PublishSubmit(ctx context.Context, mediaID string) (string, error) {
	var resp struct {
		PublishID string `json:"publish_id"`
	}
	err := c.doJSON(ctx, http.MethodPost, "/cgi-bin/freepublish/submit", nil, map[string]any{"media_id": mediaID}, &resp)
	return resp.PublishID, err
}

// PublishStatus 发布状态下发的文章详情。
type PublishStatus struct {
	PublishID     int64 `json:"publish_id"`
	PublishState  int   `json:"publish_status"`
	ArticleDetail *struct {
		Count int              `json:"count"`
		Items []map[string]any `json:"item_detail_list"`
	} `json:"article_detail,omitempty"`
}

// PublishGet 发布状态查询(POST /cgi-bin/freepublish/get)。
func (c *Client) PublishGet(ctx context.Context, publishID string) (*PublishStatus, error) {
	var resp PublishStatus
	err := c.doJSON(ctx, http.MethodPost, "/cgi-bin/freepublish/get", nil, map[string]any{"publish_id": publishID}, &resp)
	return &resp, err
}

// PublishBatchResult 已发布消息列表返回。
type PublishBatchResult struct {
	TotalCount int `json:"total_count"`
	ItemCount  int `json:"item_count"`
	Item       []struct {
		ArticleID  string         `json:"article_id"`
		Content    map[string]any `json:"content,omitempty"`
		UpdateTime int64          `json:"update_time"`
	} `json:"item"`
}

// PublishBatchGet 获取已发布的消息列表(POST /cgi-bin/freepublish/batchget)。
func (c *Client) PublishBatchGet(ctx context.Context, offset, count int, noContent bool) (*PublishBatchResult, error) {
	body := map[string]any{"offset": offset, "count": count}
	if noContent {
		body["no_content"] = 1
	}
	var resp PublishBatchResult
	err := c.doJSON(ctx, http.MethodPost, "/cgi-bin/freepublish/batchget", nil, body, &resp)
	return &resp, err
}

// PublishDelete 删除发布文章(POST /cgi-bin/freepublish/delete)。
func (c *Client) PublishDelete(ctx context.Context, articleID string) error {
	return c.doJSON(ctx, http.MethodPost, "/cgi-bin/freepublish/delete", nil, map[string]any{"article_id": articleID}, nil)
}

// GetPublishedArticle 获取已发布图文信息(POST /cgi-bin/freepublish/getarticle)。
func (c *Client) GetPublishedArticle(ctx context.Context, articleID string) ([]map[string]any, error) {
	var resp struct {
		NewsItem []map[string]any `json:"news_item"`
	}
	err := c.doJSON(ctx, http.MethodPost, "/cgi-bin/freepublish/getarticle", nil, map[string]any{"article_id": articleID}, &resp)
	return resp.NewsItem, err
}
