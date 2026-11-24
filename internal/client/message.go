package client

import (
	"context"
	"net/http"
)

// MassSendResult 根据标签群发消息的返回。
type MassSendResult struct {
	MsgID      int64 `json:"msg_id"`
	MsgDataID  int64 `json:"msg_data_id"`
}

// MassSend 根据标签群发消息(POST /cgi-bin/message/mass/sendall)。
// body 为完整请求负载(含 filter/msgtype 及对应消息体),按公众号文档填充。
func (c *Client) MassSend(ctx context.Context, body map[string]any) (*MassSendResult, error) {
	var resp MassSendResult
	err := c.doJSON(ctx, http.MethodPost, "/cgi-bin/message/mass/sendall", nil, body, &resp)
	return &resp, err
}

// MassStatus 查询群发消息发送状态(POST /cgi-bin/message/mass/get)。
func (c *Client) MassStatus(ctx context.Context, msgID int64) (string, error) {
	var resp struct {
		MsgID      int64  `json:"msg_id"`
		MsgStatus  string `json:"msg_status"`
	}
	err := c.doJSON(ctx, http.MethodPost, "/cgi-bin/message/mass/get", nil, map[string]any{"msg_id": msgID}, &resp)
	return resp.MsgStatus, err
}
