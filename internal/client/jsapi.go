package client

import (
	"context"
	"net/http"
	"net/url"
)

// JSSDKTicket JS-SDK 临时票据返回。
type JSSDKTicket struct {
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}

// GetJSSDKTicket 获取 JS-SDK 临时票据(GET /cgi-bin/ticket/getticket?type=jsapi)。
func (c *Client) GetJSSDKTicket(ctx context.Context) (*JSSDKTicket, error) {
	q := url.Values{}
	q.Set("type", "jsapi")
	var resp JSSDKTicket
	err := c.doJSON(ctx, http.MethodGet, "/cgi-bin/ticket/getticket", q, nil, &resp)
	return &resp, err
}

// SendChannelMsg 就医助手消息推送接口(POST /cityservice/sendchannelmsg)。
// body 按医疗服务主体要求携带接收用户与消息内容。
func (c *Client) SendChannelMsg(ctx context.Context, body map[string]any) error {
	return c.doJSON(ctx, http.MethodPost, "/cityservice/sendchannelmsg", nil, body, nil)
}
