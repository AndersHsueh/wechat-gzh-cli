package client

import (
	"context"
	"net/http"
)

// MassDelete 删除群发消息(POST /cgi-bin/message/mass/delete)。
func (c *Client) MassDelete(ctx context.Context, msgID int64, articleIdx int, url string) error {
	body := map[string]any{}
	if msgID != 0 {
		body["msg_id"] = msgID
		body["article_idx"] = articleIdx
	}
	if url != "" {
		body["url"] = url
	}
	return c.doJSON(ctx, http.MethodPost, "/cgi-bin/message/mass/delete", nil, body, nil)
}

// MassSpeed 群发速度返回值。
type MassSpeed struct {
	Speed     int `json:"speed"`
	RealSpeed int `json:"realspeed"`
}

// MassSpeedGet 获取群发速度(POST /cgi-bin/message/mass/speed/get)。
func (c *Client) MassSpeedGet(ctx context.Context) (*MassSpeed, error) {
	var resp MassSpeed
	err := c.doJSON(ctx, http.MethodPost, "/cgi-bin/message/mass/speed/get", nil, map[string]any{}, &resp)
	return &resp, err
}

// MassSpeedSet 设置群发速度(POST /cgi-bin/message/mass/speed/set)。
func (c *Client) MassSpeedSet(ctx context.Context, speed int) error {
	return c.doJSON(ctx, http.MethodPost, "/cgi-bin/message/mass/speed/set", nil, map[string]any{"speed": speed}, nil)
}

// MassPreview 预览消息(POST /cgi-bin/message/mass/preview)。
// 底部用 touser/towxname 二选一。
func (c *Client) MassPreview(ctx context.Context, body map[string]any) error {
	return c.doJSON(ctx, http.MethodPost, "/cgi-bin/message/mass/preview", nil, body, nil)
}

// SubscribeSend 发送一次性订阅消息(POST /cgi-bin/message/template/subscribe)。
func (c *Client) SubscribeSend(ctx context.Context, body map[string]any) error {
	return c.doJSON(ctx, http.MethodPost, "/cgi-bin/message/template/subscribe", nil, body, nil)
}

// AutoreplyInfo 自动回复规则。
type AutoreplyInfo struct {
	IsAddFriendReplyOpen   int `json:"is_add_friend_reply_open"`
	IsAutoreplyOpen        int `json:"is_autoreply_open"`
	AddFriendAutoreplyInfo struct {
		Type    string `json:"type"`
		Content string `json:"content"`
	} `json:"add_friend_autoreply_info"`
	MessageDefaultAutoreplyInfo struct {
		Type    string `json:"type"`
		Content string `json:"content"`
	} `json:"message_default_autoreply_info"`
	KeywordAutoreplyInfo struct {
		List []map[string]any `json:"list"`
	} `json:"keyword_autoreply_info"`
}

// GetAutoreply 获取自动回复规则(GET /cgi-bin/get_current_autoreply_info)。
func (c *Client) GetAutoreply(ctx context.Context) (*AutoreplyInfo, error) {
	var resp AutoreplyInfo
	err := c.doJSON(ctx, http.MethodGet, "/cgi-bin/get_current_autoreply_info", nil, nil, &resp)
	return &resp, err
}
