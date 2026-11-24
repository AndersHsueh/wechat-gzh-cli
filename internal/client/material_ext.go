package client

import (
	"context"
	"io"
	"net/http"
	"net/url"
)

// MaterialCount 永久素材总数。
type MaterialCount struct {
	VoiceCount int `json:"voice_count"`
	VideoCount int `json:"video_count"`
	ImageCount int `json:"image_count"`
	NewsCount  int `json:"news_count"`
}

// GetMaterialCount 获取永久素材总数(GET /cgi-bin/material/get_materialcount)。
func (c *Client) GetMaterialCount(ctx context.Context) (*MaterialCount, error) {
	var resp MaterialCount
	err := c.doJSON(ctx, http.MethodGet, "/cgi-bin/material/get_materialcount", nil, nil, &resp)
	return &resp, err
}

// GetTempMedia 获取临时素材(GET /cgi-bin/media/get),返回原始内容与下载 URL。
// 对视频素材,响应体为 {"video_url": "..."};其它类型为原始文件字节流。
func (c *Client) GetTempMedia(ctx context.Context, mediaID string) ([]byte, error) {
	token, err := c.AccessToken(ctx)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("access_token", token)
	q.Set("media_id", mediaID)
	endpoint := c.cfg.BaseURL + "/cgi-bin/media/get?" + q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
