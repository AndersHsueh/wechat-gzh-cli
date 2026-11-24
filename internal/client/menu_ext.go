package client

import (
	"context"
	"net/http"
)

// SelfMenu 查询自定义菜单信息返回(含官网后台设置的菜单)。
type SelfMenu struct {
	IsMenuOpen   int `json:"is_menu_open"`
	SelfmenuInfo struct {
		Button []map[string]any `json:"button"`
	} `json:"selfmenu_info"`
}

// GetSelfMenu 查询自定义菜单信息(GET /cgi-bin/get_current_selfmenu_info)。
func (c *Client) GetSelfMenu(ctx context.Context) (*SelfMenu, error) {
	var resp SelfMenu
	err := c.doJSON(ctx, http.MethodGet, "/cgi-bin/get_current_selfmenu_info", nil, nil, &resp)
	return &resp, err
}

// ConditionalMenuResult 创建个性化菜单的返回。
type ConditionalMenuResult struct {
	MenuID string `json:"menuid"`
}

// AddConditionalMenu 创建个性化菜单(POST /cgi-bin/menu/addconditional)。
// body 应含 button 与 matchrule。
func (c *Client) AddConditionalMenu(ctx context.Context, body map[string]any) (*ConditionalMenuResult, error) {
	var resp ConditionalMenuResult
	err := c.doJSON(ctx, http.MethodPost, "/cgi-bin/menu/addconditional", nil, body, &resp)
	return &resp, err
}

// DelConditionalMenu 删除个性化菜单(POST /cgi-bin/menu/delconditional)。
func (c *Client) DelConditionalMenu(ctx context.Context, menuid string) error {
	return c.doJSON(ctx, http.MethodPost, "/cgi-bin/menu/delconditional", nil, map[string]any{"menuid": menuid}, nil)
}

// TryMatchMenu 测试个性化菜单匹配结果(POST /cgi-bin/menu/trymatch)。
func (c *Client) TryMatchMenu(ctx context.Context, userID string) ([]map[string]any, error) {
	var resp struct {
		Button []map[string]any `json:"button"`
	}
	err := c.doJSON(ctx, http.MethodPost, "/cgi-bin/menu/trymatch", nil, map[string]any{"user_id": userID}, &resp)
	return resp.Button, err
}
