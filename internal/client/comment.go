package client

import (
	"context"
	"net/http"
)

// CommentItem 评论项。
type CommentItem struct {
	CommentID  int64  `json:"user_comment_id"`
	OpenID     string `json:"openid"`
	CreateTime int64  `json:"create_time"`
	Content    string `json:"content"`
	Reply      struct {
		Content string `json:"content"`
	} `json:"reply"`
	IsTop   int `json:"is_top"`
	IsElect int `json:"is_elect"`
}

// CommentList 评论列表返回。
type CommentList struct {
	Total   int           `json:"total"`
	Comment []CommentItem `json:"comment"`
}

// OpenComment 打开已群发文章评论(POST /cgi-bin/comment/open)。
func (c *Client) OpenComment(ctx context.Context, msgDataID, index int64) error {
	return c.doJSON(ctx, http.MethodPost, "/cgi-bin/comment/open", nil,
		map[string]any{"msg_data_id": msgDataID, "index": index}, nil)
}

// CloseComment 关闭已群发文章评论(POST /cgi-bin/comment/close)。
func (c *Client) CloseComment(ctx context.Context, msgDataID, index int64) error {
	return c.doJSON(ctx, http.MethodPost, "/cgi-bin/comment/close", nil,
		map[string]any{"msg_data_id": msgDataID, "index": index}, nil)
}

// ListComment 查看指定文章的评论数据(POST /cgi-bin/comment/list)。
// type: 0 全部,1 仅精选,2 仅未精选。
func (c *Client) ListComment(ctx context.Context, msgDataID, index int64, begin, count, commentType int) (*CommentList, error) {
	var resp CommentList
	err := c.doJSON(ctx, http.MethodPost, "/cgi-bin/comment/list", nil,
		map[string]any{"msg_data_id": msgDataID, "index": index, "begin": begin, "count": count, "type": commentType}, &resp)
	return &resp, err
}

// MarkElect 评论标记精选(POST /cgi-bin/comment/markelect)。
func (c *Client) MarkElect(ctx context.Context, commentID int64) error {
	return c.doJSON(ctx, http.MethodPost, "/cgi-bin/comment/markelect", nil,
		map[string]any{"comment_id": commentID}, nil)
}

// UnmarkElect 评论取消精选(POST /cgi-bin/comment/unmarkelect)。
func (c *Client) UnmarkElect(ctx context.Context, commentID int64) error {
	return c.doJSON(ctx, http.MethodPost, "/cgi-bin/comment/unmarkelect", nil,
		map[string]any{"comment_id": commentID}, nil)
}

// DeleteComment 删除评论(POST /cgi-bin/comment/delete)。
func (c *Client) DeleteComment(ctx context.Context, commentID int64) error {
	return c.doJSON(ctx, http.MethodPost, "/cgi-bin/comment/delete", nil,
		map[string]any{"comment_id": commentID}, nil)
}

// ReplyComment 回复评论(POST /cgi-bin/comment/reply/add)。
func (c *Client) ReplyComment(ctx context.Context, commentID int64, content string) error {
	return c.doJSON(ctx, http.MethodPost, "/cgi-bin/comment/reply/add", nil,
		map[string]any{"comment_id": commentID, "content": content}, nil)
}

// DeleteReply 删除回复(POST /cgi-bin/comment/reply/delete)。
func (c *Client) DeleteReply(ctx context.Context, commentID int64) error {
	return c.doJSON(ctx, http.MethodPost, "/cgi-bin/comment/reply/delete", nil,
		map[string]any{"comment_id": commentID}, nil)
}
