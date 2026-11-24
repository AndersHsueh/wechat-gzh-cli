package client

import (
	"context"
	"net/http"
)

// CreateMenu 创建自定义菜单(POST /cgi-bin/menu/create)。
// menu 为完整菜单负载({"button": [...]})。
func (c *Client) CreateMenu(ctx context.Context, menu map[string]any) error {
	return c.doJSON(ctx, http.MethodPost, "/cgi-bin/menu/create", nil, menu, nil)
}

// FullMenu 自定义菜单配置(默认 + 个性化)。
type FullMenu struct {
	Menu struct {
		Button []map[string]any `json:"button"`
		MenuID int64            `json:"menuid"`
	} `json:"menu"`
	ConditionalMenu []map[string]any `json:"conditionalmenu"`
}

// GetMenu 获取自定义菜单配置(GET /cgi-bin/menu/get)。
func (c *Client) GetMenu(ctx context.Context) (*FullMenu, error) {
	var resp FullMenu
	err := c.doJSON(ctx, http.MethodGet, "/cgi-bin/menu/get", nil, nil, &resp)
	return &resp, err
}

// DeleteMenu 删除自定义菜单(GET /cgi-bin/menu/delete)。
func (c *Client) DeleteMenu(ctx context.Context) error {
	return c.doJSON(ctx, http.MethodGet, "/cgi-bin/menu/delete", nil, nil, nil)
}
