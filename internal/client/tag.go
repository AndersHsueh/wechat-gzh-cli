package client

import (
	"context"
	"net/http"
)

// Tag 用户标签。
type Tag struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// ListTags 获取公众号已创建的标签列表(GET /cgi-bin/tags/get)。
func (c *Client) ListTags(ctx context.Context) ([]Tag, error) {
	var resp struct {
		Tags []Tag `json:"tags"`
	}
	err := c.doJSON(ctx, http.MethodGet, "/cgi-bin/tags/get", nil, nil, &resp)
	return resp.Tags, err
}

// CreateTag 创建公众号标签(POST /cgi-bin/tags/create)。
func (c *Client) CreateTag(ctx context.Context, name string) (*Tag, error) {
	var resp struct {
		Tag *Tag `json:"tag"`
	}
	err := c.doJSON(ctx, http.MethodPost, "/cgi-bin/tags/create", nil,
		map[string]any{"tag": map[string]any{"name": name}}, &resp)
	return resp.Tag, err
}
