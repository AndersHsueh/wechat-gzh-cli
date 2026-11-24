package client

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

// TranslateContent 微信翻译(POST /cgi-bin/media/voice/translatecontent)。
// lfrom/lto 为源/目标语言,如 zh_CN/en_US。content 为要翻译的文本。
func (c *Client) TranslateContent(ctx context.Context, lfrom, lto, content string) (string, error) {
	q := url.Values{}
	q.Set("lfrom", lfrom)
	q.Set("lto", lto)
	var resp struct {
		ToContent string `json:"tocontent"`
	}
	err := c.doJSON(ctx, http.MethodPost, "/cgi-bin/media/voice/translatecontent", q,
		map[string]any{"content": content}, &resp)
	return resp.ToContent, err
}

// VoiceToText 上传语音文件转文字(POST /cgi-bin/media/voice/addvoicetorecofortext)。
// voiceID 用于后续查询识别结果。
func (c *Client) VoiceToText(ctx context.Context, format, voiceID, lang, filepath string) error {
	token, err := c.AccessToken(ctx)
	if err != nil {
		return err
	}
	if lang == "" {
		lang = "zh_CN"
	}
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	fw, err := w.CreateFormFile("voice_data", filepath)
	if err != nil {
		return err
	}
	if _, err := io.Copy(fw, f); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	q := url.Values{}
	q.Set("access_token", token)
	q.Set("format", format)
	q.Set("voice_id", voiceID)
	q.Set("lang", lang)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.BaseURL+"/cgi-bin/media/voice/addvoicetorecofortext?"+q.Encode(), body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return parseAPIError(raw)
}

// VoiceResult 查询语音识别结果(POST /cgi-bin/media/voice/queryrecoresultfortext)。
func (c *Client) VoiceResult(ctx context.Context, voiceID, lang string) (map[string]any, error) {
	if lang == "" {
		lang = "zh_CN"
	}
	q := url.Values{}
	q.Set("voice_id", voiceID)
	q.Set("lang", lang)
	var resp map[string]any
	err := c.doJSON(ctx, http.MethodPost, "/cgi-bin/media/voice/queryrecoresultfortext", q, nil, &resp)
	return resp, err
}

// OcrEndpoint OCR 接口定义。
type OcrEndpoint struct {
	// Name 子命令名。
	Name string
	// Path OCR 接口路径(含 /cv 前缀)。
	Path string
	// Summary 说明。
	Summary string
}

// OcrEndpoints 全部 OCR 与图像处理接口(基于 img_url)。
var OcrEndpoints = []OcrEndpoint{
	{Name: "ocr-menu", Path: "/cv/ocr/menu", Summary: "菜单识别"},
	{Name: "ocr-comm", Path: "/cv/ocr/comm", Summary: "通用印刷体识别"},
	{Name: "ocr-driving", Path: "/cv/ocr/driving", Summary: "行驶证识别"},
	{Name: "ocr-bankcard", Path: "/cv/ocr/bankcard", Summary: "银行卡识别"},
	{Name: "ocr-bizlicense", Path: "/cv/ocr/bizlicense", Summary: "营业执照识别"},
	{Name: "ocr-driverlicense", Path: "/cv/ocr/drivinglicense", Summary: "驾驶证识别"},
	{Name: "ocr-idcard", Path: "/cv/ocr/idcard", Summary: "身份证识别"},
	{Name: "img-aicrop", Path: "/cv/img/aicrop", Summary: "图片智能裁剪"},
	{Name: "img-qrcode", Path: "/cv/img/qrcode", Summary: "二维码/条码识别"},
}

// OcrQuery 基于 img_url 调用 OCR/图像处理接口,返回解析后的结构字段。
func (c *Client) OcrQuery(ctx context.Context, path, imgURL string) (map[string]any, error) {
	q := url.Values{}
	q.Set("img_url", imgURL)
	var resp map[string]any
	err := c.doJSON(ctx, http.MethodPost, path, q, nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetJSSDKRawVoice 获取高质量语音素材(GET /cgi-bin/media/get/jssdk),返回原始字节。
func (c *Client) GetJSSDKRawVoice(ctx context.Context, mediaID string) ([]byte, error) {
	token, err := c.AccessToken(ctx)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("access_token", token)
	q.Set("media_id", mediaID)
	endpoint := c.cfg.BaseURL + "/cgi-bin/media/get/jssdk?" + q.Encode()
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
