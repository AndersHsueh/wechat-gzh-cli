package client

import (
	"context"
	"net/http"
)

// GetApiDomainIP 获取微信 API 服务器 IP 列表(GET /cgi-bin/get_api_domain_ip)。
func (c *Client) GetApiDomainIP(ctx context.Context) ([]string, error) {
	var resp struct {
		IPList []string `json:"ip_list"`
	}
	if err := c.doJSON(ctx, http.MethodGet, "/cgi-bin/get_api_domain_ip", nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp.IPList, nil
}
