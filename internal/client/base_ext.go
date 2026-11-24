package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// CallbackCheckResult 网络通信检测返回。
type CallbackCheckResult struct {
	DNS  []DNSItem  `json:"dns"`
	Ping []PingItem `json:"ping"`
}

// DNSItem DNS 解析结果项。
type DNSItem struct {
	IP           string `json:"ip"`
	RealOperator string `json:"real_operator"`
}

// PingItem ping 检测结果项。
type PingItem struct {
	IP           string `json:"ip"`
	FromOperator string `json:"from_operator"`
	PackageLoss  string `json:"package_loss"`
	Time         string `json:"time"`
}

// StableAccessToken 获取稳定版接口调用凭据(POST /cgi-bin/stable_token)。
// 不依赖 access_token(使用 AppID/AppSecret)。forceRefresh 为空则走普通模式。
func (c *Client) StableAccessToken(ctx context.Context, forceRefresh *bool) (string, int, error) {
	body := map[string]any{
		"grant_type": "client_credential",
		"appid":      c.cfg.AppID,
		"secret":     c.cfg.AppSecret,
	}
	if forceRefresh != nil {
		body["force_refresh"] = *forceRefresh
	}
	payload, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.BaseURL+"/cgi-bin/stable_token", bytes.NewReader(payload))
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.hc.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}
	var out struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return "", 0, err
	}
	return out.AccessToken, out.ExpiresIn, nil
}

// CallbackCheck 网络通信检测(POST /cgi-bin/callback/check)。
func (c *Client) CallbackCheck(ctx context.Context, action, checkOperator string) (*CallbackCheckResult, error) {
	var resp CallbackCheckResult
	err := c.doJSON(ctx, http.MethodPost, "/cgi-bin/callback/check", nil,
		map[string]any{"action": action, "check_operator": checkOperator}, &resp)
	return &resp, err
}

// GetCallbackIP 获取微信推送服务器 IP(GET /cgi-bin/getcallbackip)。
func (c *Client) GetCallbackIP(ctx context.Context) ([]string, error) {
	var resp struct {
		IPList []string `json:"ip_list"`
	}
	if err := c.doJSON(ctx, http.MethodGet, "/cgi-bin/getcallbackip", nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp.IPList, nil
}

// ResetQuota 重置 API 每日调用次数(POST /cgi-bin/clear_quota)。
func (c *Client) ResetQuota(ctx context.Context, appid string) error {
	return c.doJSON(ctx, http.MethodPost, "/cgi-bin/clear_quota", nil, map[string]any{"appid": appid}, nil)
}

// ApiQuota 接口调用额度。
type ApiQuota struct {
	DailyLimit int `json:"daily_limit"`
	Used       int `json:"used"`
	Remain     int `json:"remain"`
}

// GetApiQuota 查询接口调用额度(POST /cgi-bin/openapi/quota/get)。
func (c *Client) GetApiQuota(ctx context.Context, cgiPath string) (*ApiQuota, error) {
	var resp struct {
		Quota ApiQuota `json:"quota"`
	}
	err := c.doJSON(ctx, http.MethodPost, "/cgi-bin/openapi/quota/get", nil, map[string]any{"cgi_path": cgiPath}, &resp)
	return &resp.Quota, err
}

// RidRequest 查询 rid 信息返回的请求详情。
type RidRequest struct {
	InvokeTime   int    `json:"invoke_time"`
	CostInMS     int    `json:"cost_in_ms"`
	RequestURL   string `json:"request_url"`
	RequestBody  string `json:"request_body"`
	ResponseBody string `json:"response_body"`
	ClientIP     string `json:"client_ip"`
}

// GetRidInfo 查询 rid 信息(POST /cgi-bin/openapi/rid/get)。
func (c *Client) GetRidInfo(ctx context.Context, rid string) (*RidRequest, error) {
	var resp struct {
		Request RidRequest `json:"request"`
	}
	err := c.doJSON(ctx, http.MethodPost, "/cgi-bin/openapi/rid/get", nil, map[string]any{"rid": rid}, &resp)
	return &resp.Request, err
}

// ClearApiQuota 重置指定 API 调用次数(POST /cgi-bin/openapi/quota/clear)。
func (c *Client) ClearApiQuota(ctx context.Context, cgiPath string) error {
	return c.doJSON(ctx, http.MethodPost, "/cgi-bin/openapi/quota/clear", nil, map[string]any{"cgi_path": cgiPath}, nil)
}
