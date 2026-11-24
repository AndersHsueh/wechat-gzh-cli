package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

// NewsItem 图文素材项(与公众号文档字段一致)。
type NewsItem struct {
	Title            string `json:"title"`
	Author           string `json:"author"`
	Digest           string `json:"digest"`
	Content          string `json:"content"`
	ContentSourceURL string `json:"content_source_url"`
	ThumbMediaID     string `json:"thumb_media_id"`
	ShowCoverPic     int    `json:"show_cover_pic"`
	URL              string `json:"url"`
}

// MaterialItem 素材列表项。
type MaterialItem struct {
	MediaID    string `json:"media_id"`
	Name       string `json:"name"`
	URL        string `json:"url"`
	UpdateTime int64  `json:"update_time"`
	Content    *struct {
		NewsItem []NewsItem `json:"news_item"`
	} `json:"content,omitempty"`
}

// AddMaterialResult 上传永久素材的返回。
type AddMaterialResult struct {
	MediaID string `json:"media_id"`
	URL     string `json:"url"`
}

// BatchGetMaterialResult 素材列表返回。
type BatchGetMaterialResult struct {
	TotalCount int            `json:"total_count"`
	ItemCount  int            `json:"item_count"`
	Item       []MaterialItem `json:"item"`
}

// doUpload 用 multipart/form-data 上传文件到指定 path。
func (c *Client) doUpload(ctx context.Context, path string, query url.Values, field, filename string, data io.Reader, out any) error {
	token, err := c.AccessToken(ctx)
	if err != nil {
		return err
	}
	if query == nil {
		query = url.Values{}
	}
	query.Set("access_token", token)

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, err := mw.CreateFormFile(field, filename)
	if err != nil {
		return err
	}
	if _, err := io.Copy(fw, data); err != nil {
		return err
	}
	if err := mw.Close(); err != nil {
		return err
	}

	endpoint := c.cfg.BaseURL + path + "?" + query.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())

	resp, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if e := parseAPIError(raw); e != nil {
		return e
	}
	if out != nil {
		return jsonUnmarshal(raw, out)
	}
	return nil
}

// UploadImage 上传发表内容中的图片(POST /cgi-bin/media/uploadimg),返回图片 URL。
func (c *Client) UploadImage(ctx context.Context, filename string, data io.Reader) (string, error) {
	var resp struct {
		URL string `json:"url"`
	}
	err := c.doUpload(ctx, "/cgi-bin/media/uploadimg", nil, "media", filename, data, &resp)
	return resp.URL, err
}

// UploadTempMedia 新增临时素材(POST /cgi-bin/media/upload)。
func (c *Client) UploadTempMedia(ctx context.Context, mediaType, filename string, data io.Reader) (*UploadResult, error) {
	q := url.Values{}
	q.Set("type", mediaType)
	var resp UploadResult
	if err := c.doUpload(ctx, "/cgi-bin/media/upload", q, "media", filename, data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UploadResult 临时素材上传返回。
type UploadResult struct {
	Type      string `json:"type"`
	MediaID   string `json:"media_id"`
	CreatedAt int64  `json:"created_at"`
}

// AddMaterial 新增永久素材(POST /cgi-bin/material/add_material)。
func (c *Client) AddMaterial(ctx context.Context, mediaType, filename string, data io.Reader) (*AddMaterialResult, error) {
	q := url.Values{}
	q.Set("type", mediaType)
	var resp AddMaterialResult
	if err := c.doUpload(ctx, "/cgi-bin/material/add_material", q, "media", filename, data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetMaterial 获取永久素材(POST /cgi-bin/material/get_material)。
func (c *Client) GetMaterial(ctx context.Context, mediaID string) ([]NewsItem, error) {
	var resp struct {
		NewsItem []NewsItem `json:"news_item"`
	}
	err := c.doJSON(ctx, http.MethodPost, "/cgi-bin/material/get_material", nil, map[string]any{"media_id": mediaID}, &resp)
	return resp.NewsItem, err
}

// BatchGetMaterial 分类型获取永久素材列表(POST /cgi-bin/material/batchget_material)。
func (c *Client) BatchGetMaterial(ctx context.Context, mediaType string, offset, count int) (*BatchGetMaterialResult, error) {
	var resp BatchGetMaterialResult
	err := c.doJSON(ctx, http.MethodPost, "/cgi-bin/material/batchget_material", nil,
		map[string]any{"type": mediaType, "offset": offset, "count": count}, &resp)
	return &resp, err
}

// DelMaterial 删除永久素材(POST /cgi-bin/material/del_material)。
func (c *Client) DelMaterial(ctx context.Context, mediaID string) error {
	return c.doJSON(ctx, http.MethodPost, "/cgi-bin/material/del_material", nil, map[string]any{"media_id": mediaID}, nil)
}

// parseAPIError 从响应字节提取 errcode!=0 的 apiError。
func parseAPIError(raw []byte) error {
	var box struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	if err := json.Unmarshal(raw, &box); err == nil && box.ErrCode != 0 {
		return &apiError{ErrCode: box.ErrCode, ErrMsg: box.ErrMsg}
	}
	return nil
}

func jsonUnmarshal(raw []byte, v any) error {
	return json.Unmarshal(raw, v)
}
