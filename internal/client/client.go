// Package client 封装微信公众号服务端接口的 HTTP 调用与鉴权。
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

// 环境变量名。
const (
	EnvAppID     = "WX_GZ_AppID"
	EnvAppSecret = "WX_GZ_AppSecret"
)

// 微信 API 默认基础地址。
const defaultBaseURL = "https://api.weixin.qq.com"

// accessTokenTolerance 在到期前提前多少秒视为过期,用于安全刷新。
const accessTokenTolerance = 60

// Config 是 client 的配置。
type Config struct {
	// AppID 公众号 appid。
	AppID string
	// AppSecret 公众号 appsecret。
	AppSecret string
	// BaseURL 微信 API 基础地址,默认为 https://api.weixin.qq.com。
	BaseURL string
	// HTTPClient 底层 HTTP 客户端,默认为 http.DefaultClient。
	HTTPClient *http.Client
}

// Client 是对微信公众平台服务端接口的调用封装。
type Client struct {
	cfg Config
	hc  *http.Client

	mu          sync.Mutex
	token       string
	tokenExpiry time.Time
}

// New 构造 Client,并填充默认值。
func New(cfg Config) *Client {
	if cfg.BaseURL == "" {
		cfg.BaseURL = defaultBaseURL
	}
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = http.DefaultClient
	}
	return &Client{cfg: cfg, hc: cfg.HTTPClient}
}

// FromEnv 从环境变量构造 Client;缺失必需变量时返回错误。
func FromEnv() (*Client, error) {
	appID := os.Getenv(EnvAppID)
	if appID == "" {
		return nil, fmt.Errorf("missing env %s", EnvAppID)
	}
	secret := os.Getenv(EnvAppSecret)
	if secret == "" {
		return nil, fmt.Errorf("missing env %s", EnvAppSecret)
	}
	return New(Config{AppID: appID, AppSecret: secret}), nil
}

// apiError 表示微信接口返回的错误(errcode 非 0)。
type apiError struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (e *apiError) Error() string {
	return fmt.Sprintf("wechat api error %d: %s", e.ErrCode, e.ErrMsg)
}

// tokenResponse 是 /cgi-bin/token 的返回结构。
type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	apiError
}

// AccessToken 获取(或缓存)全局接口调用凭据 access_token。
// access_token 有效期 7200 秒;有效期内不会重复请求。
func (c *Client) AccessToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.token != "" && time.Now().Add(accessTokenTolerance*time.Second).Before(c.tokenExpiry) {
		return c.token, nil
	}

	endpoint := c.cfg.BaseURL + "/cgi-bin/token"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}
	q := req.URL.Query()
	q.Set("grant_type", "client_credential")
	q.Set("appid", c.cfg.AppID)
	q.Set("secret", c.cfg.AppSecret)
	req.URL.RawQuery = q.Encode()

	resp, err := c.hc.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tr tokenResponse
	if err := json.Unmarshal(body, &tr); err != nil {
		return "", fmt.Errorf("decode token response: %w", err)
	}
	if tr.ErrCode != 0 {
		return "", &tr.apiError
	}
	if tr.AccessToken == "" {
		return "", errors.New("access_token empty in response")
	}

	c.token = tr.AccessToken
	c.tokenExpiry = time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)
	return c.token, nil
}

// doJSON 发送 JSON 请求到指定 path(query 追加 access_token),并解析响应到 out。
// out 可以是包含 errcode/errmsg 的结构;若响应 errcode!=0 则返回 apiError。
func (c *Client) doJSON(ctx context.Context, method, path string, query url.Values, in, out any) error {
	token, err := c.AccessToken(ctx)
	if err != nil {
		return err
	}
	if query == nil {
		query = url.Values{}
	}
	query.Set("access_token", token)

	endpoint := c.cfg.BaseURL + path + "?" + query.Encode()

	var bodyReader io.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, bodyReader)
	if err != nil {
		return err
	}
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 通用检查:响应中若带 errcode 且非 0,视为业务错误。
	var errBox struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	if err := json.Unmarshal(raw, &errBox); err == nil && errBox.ErrCode != 0 {
		return &apiError{ErrCode: errBox.ErrCode, ErrMsg: errBox.ErrMsg}
	}

	if out != nil {
		if err := json.Unmarshal(raw, out); err != nil {
			return fmt.Errorf("decode %s response: %w", path, err)
		}
	}
	return nil
}
