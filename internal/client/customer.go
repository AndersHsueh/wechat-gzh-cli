package client

import (
	"context"
	"net/http"
)

// SendCustomerMessage 发送客服消息(POST /cgi-bin/message/custom/send)。
// msgtype 支持 text/image/voice/video/music/news/mpnews/mpnewsarticle/msgmenu/wxcard/miniprogrampage。
// payload 为对应消息类型的具体内容对象,会放入 body[msgtype] 下。
func (c *Client) SendCustomerMessage(ctx context.Context, touser, msgtype string, payload map[string]any) error {
	body := map[string]any{"touser": touser, "msgtype": msgtype}
	body[msgtype] = payload
	return c.doJSON(ctx, http.MethodPost, "/cgi-bin/message/custom/send", nil, body, nil)
}
