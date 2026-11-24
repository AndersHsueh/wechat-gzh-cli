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

// KfAccount 客服账号。
type KfAccount struct {
	KfAccount     string `json:"kf_account"`
	KfNick        string `json:"kf_nick"`
	KfHeadImg     string `json:"kf_headimgurl"`
	InviteWx      string `json:"invite_wx"`
	InviteExpTime int64  `json:"invite_expire_time"`
	KfID          int    `json:"kf_id"`
}

// KfListResult 客服列表返回。
type KfListResult struct {
	KFList []KfAccount `json:"kf_list"`
}

// OnlineKfAccount 在线客服。
type OnlineKfAccount struct {
	KfAccount    string `json:"kf_account"`
	Status       int    `json:"status"`
	KfID         int    `json:"kf_id"`
	AutoAccept   int    `json:"auto_accept"`
	AcceptedCase int    `json:"accepted_case"`
}

// OnlineKfListResult 在线客服列表返回。
type OnlineKfListResult struct {
	KFOnlineList []OnlineKfAccount `json:"kf_online_list"`
}

// AddKfAccount 添加客服账号(POST /customservice/kfaccount/add)。
func (c *Client) AddKfAccount(ctx context.Context, kfAccount, nickname string) error {
	return c.doJSON(ctx, http.MethodPost, "/customservice/kfaccount/add", nil,
		map[string]any{"kf_account": kfAccount, "nickname": nickname}, nil)
}

// UpdateKfAccount 修改客服账号(POST /customservice/kfaccount/update)。
func (c *Client) UpdateKfAccount(ctx context.Context, kfAccount, nickname string) error {
	return c.doJSON(ctx, http.MethodPost, "/customservice/kfaccount/update", nil,
		map[string]any{"kf_account": kfAccount, "nickname": nickname}, nil)
}

// DelKfAccount 删除客服账号(POST /customservice/kfaccount/del)。
func (c *Client) DelKfAccount(ctx context.Context, kfAccount string) error {
	return c.doJSON(ctx, http.MethodPost, "/customservice/kfaccount/del", nil,
		map[string]any{"kf_account": kfAccount}, nil)
}

// UploadKfHeadimg 设置客服账号头像(POST /customservice/kfaccount/uploadheadimg)。
// 以 multipart/form-data 上传 media 字段,且 kf_account 需带在 query。
func (c *Client) UploadKfHeadimg(ctx context.Context, kfAccount, filepath string) error {
	token, err := c.AccessToken(ctx)
	if err != nil {
		return err
	}
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	fw, err := w.CreateFormFile("media", filepath)
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
	q.Set("kf_account", kfAccount)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.BaseURL+"/customservice/kfaccount/uploadheadimg?"+q.Encode(), body)
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

// GetKfList 获取所有客服账号(GET /cgi-bin/customservice/getkflist)。
func (c *Client) GetKfList(ctx context.Context) (*KfListResult, error) {
	var resp KfListResult
	err := c.doJSON(ctx, http.MethodGet, "/cgi-bin/customservice/getkflist", nil, nil, &resp)
	return &resp, err
}

// GetOnlineKfList 获取在线客服列表(GET /cgi-bin/customservice/getonlinekflist)。
func (c *Client) GetOnlineKfList(ctx context.Context) (*OnlineKfListResult, error) {
	var resp OnlineKfListResult
	err := c.doJSON(ctx, http.MethodGet, "/cgi-bin/customservice/getonlinekflist", nil, nil, &resp)
	return &resp, err
}

// InviteKfWorker 邀请绑定客服账号(POST /customservice/kfaccount/inviteworker)。
func (c *Client) InviteKfWorker(ctx context.Context, kfAccount, inviteWx string) error {
	return c.doJSON(ctx, http.MethodPost, "/customservice/kfaccount/inviteworker", nil,
		map[string]any{"kf_account": kfAccount, "invite_wx": inviteWx}, nil)
}

// CreateKfSession 创建会话(POST /customservice/kfsession/create)。
func (c *Client) CreateKfSession(ctx context.Context, kfAccount, openid string) error {
	return c.doJSON(ctx, http.MethodPost, "/customservice/kfsession/create", nil,
		map[string]any{"kf_account": kfAccount, "openid": openid}, nil)
}

// CloseKfSession 关闭会话(POST /customservice/kfsession/close)。
func (c *Client) CloseKfSession(ctx context.Context, kfAccount, openid string) error {
	return c.doJSON(ctx, http.MethodPost, "/customservice/kfsession/close", nil,
		map[string]any{"kf_account": kfAccount, "openid": openid}, nil)
}

// KfSessionResult 客户会话状态。
type KfSessionResult struct {
	KfAccount  string `json:"kf_account"`
	CreateTime int64  `json:"createtime"`
}

// GetKfSession 获取客户会话状态(GET /customservice/kfsession/getsession)。
func (c *Client) GetKfSession(ctx context.Context, openid string) (*KfSessionResult, error) {
	q := url.Values{}
	q.Set("openid", openid)
	var resp KfSessionResult
	err := c.doJSON(ctx, http.MethodGet, "/customservice/kfsession/getsession", q, nil, &resp)
	return &resp, err
}

// SessionListItem 会话列表现。
type SessionListItem struct {
	OpenID     string `json:"openid"`
	CreateTime int64  `json:"createtime"`
}

// GetKfSessionList 获取客服会话列表(GET /customservice/kfsession/getsessionlist)。
func (c *Client) GetKfSessionList(ctx context.Context, kfAccount string) ([]SessionListItem, error) {
	q := url.Values{}
	q.Set("kf_account", kfAccount)
	var resp struct {
		SessionList []SessionListItem `json:"sessionlist"`
	}
	err := c.doJSON(ctx, http.MethodGet, "/customservice/kfsession/getsessionlist", q, nil, &resp)
	return resp.SessionList, err
}

// GetWaitCase 获取未接入会话列表(GET /customservice/kfsession/getwaitcase)。
func (c *Client) GetWaitCase(ctx context.Context) ([]SessionListItem, error) {
	var resp struct {
		WaitCaseList []SessionListItem `json:"waitcaselist"`
	}
	err := c.doJSON(ctx, http.MethodGet, "/customservice/kfsession/getwaitcase", nil, nil, &resp)
	return resp.WaitCaseList, err
}

// GetMsgList 获取聊天记录(POST /customservice/msgrecord/getmsglist)。
func (c *Client) GetMsgList(ctx context.Context, startTime, endTime, msgID int64, number int) ([]map[string]any, error) {
	var resp struct {
		RecordList []map[string]any `json:"recordlist"`
	}
	err := c.doJSON(ctx, http.MethodPost, "/customservice/msgrecord/getmsglist", nil,
		map[string]any{"starttime": startTime, "endtime": endTime, "msgid": msgID, "number": number}, &resp)
	return resp.RecordList, err
}

// SendTyping 客服输入状态(POST /cgi-bin/message/custom/typing)。
// command 为 Typing 或 CancelTyping。
func (c *Client) SendTyping(ctx context.Context, touser, command string) error {
	return c.doJSON(ctx, http.MethodPost, "/cgi-bin/message/custom/typing", nil,
		map[string]any{"touser": touser, "command": command}, nil)
}
