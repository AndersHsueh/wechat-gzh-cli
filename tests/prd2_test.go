package tests

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ---- base 扩展 ----

// TestGetMaterialCount 验证素材总数解析。
func TestGetMaterialCount(t *testing.T) {
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/material/get_materialcount" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"voice_count": 1, "video_count": 2, "image_count": 3, "news_count": 4})
	}))
	defer srv.Close()

	mc, err := newClientWithServer(t, srv).GetMaterialCount(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if mc.ImageCount != 3 || mc.NewsCount != 4 {
		t.Errorf("mc = %+v", mc)
	}
}

// TestGetTempMedia 验证临时素材下载返回原始字节并带 media_id 查询参数。
func TestGetTempMedia(t *testing.T) {
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/media/get" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("media_id") != "m_tmp" {
			t.Errorf("media_id = %q", r.URL.Query().Get("media_id"))
		}
		w.Write([]byte("rawfile"))
	}))
	defer srv.Close()

	data, err := newClientWithServer(t, srv).GetTempMedia(context.Background(), "m_tmp")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if string(data) != "rawfile" {
		t.Errorf("data = %q", data)
	}
}

// TestStableAccessToken 验证 stable_token 请求体与解析。
func TestStableAccessToken(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/stable_token" {
			t.Errorf("path = %q", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"access_token": "ST", "expires_in": 7200})
	}))
	defer srv.Close()

	tok, exp, err := newClientWithServer(t, srv).StableAccessToken(context.Background(), nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if gotBody["grant_type"] != "client_credential" || gotBody["appid"] != "a" {
		t.Errorf("body = %v", gotBody)
	}
	if tok != "ST" || exp != 7200 {
		t.Errorf("tok=%q exp=%d", tok, exp)
	}
}

// TestGetApiQuota 验证接口额度的请求体与解析。
func TestGetApiQuota(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/openapi/quota/get" {
			t.Errorf("path = %q", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"quota": map[string]any{"daily_limit": 1000, "used": 10, "remain": 990}})
	}))
	defer srv.Close()

	q, err := newClientWithServer(t, srv).GetApiQuota(context.Background(), "/cgi-bin/message/custom/send")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if gotBody["cgi_path"] != "/cgi-bin/message/custom/send" {
		t.Errorf("cgi_path = %v", gotBody["cgi_path"])
	}
	if q.Remain != 990 || q.DailyLimit != 1000 {
		t.Errorf("q = %+v", q)
	}
}

// TestGetCallbackIP 验证推送服务器 IP 解析。
func TestGetCallbackIP(t *testing.T) {
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/getcallbackip" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"ip_list": []string{"10.0.0.1"}})
	}))
	defer srv.Close()

	ips, err := newClientWithServer(t, srv).GetCallbackIP(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(ips) != 1 || ips[0] != "10.0.0.1" {
		t.Errorf("ips = %v", ips)
	}
}

// ---- 菜单扩展 ----

// TestAddConditionalMenu 验证个性化菜单 body 与 menuid 解析。
func TestAddConditionalMenu(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/menu/addconditional" {
			t.Errorf("path = %q", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"menuid": "208"})
	}))
	defer srv.Close()

	res, err := newClientWithServer(t, srv).AddConditionalMenu(context.Background(), map[string]any{"button": []any{}})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if gotBody["button"] == nil {
		t.Errorf("body = %v", gotBody)
	}
	if res.MenuID != "208" {
		t.Errorf("menuid = %q", res.MenuID)
	}
}

// ---- 消息扩展 ----

// TestSubscribeSend 验证订阅消息 body 转发。
func TestSubscribeSend(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/message/template/subscribe" {
			t.Errorf("path = %q", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"errcode": 0, "errmsg": "ok"})
	}))
	defer srv.Close()

	err := newClientWithServer(t, srv).SubscribeSend(context.Background(), map[string]any{"touser": "o_1", "template_id": "T"})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if gotBody["touser"] != "o_1" || gotBody["template_id"] != "T" {
		t.Errorf("body = %v", gotBody)
	}
}

// TestGetAutoreply 验证自动回复规则解析。
func TestGetAutoreply(t *testing.T) {
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/get_current_autoreply_info" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"is_add_friend_reply_open": 1, "is_autoreply_open": 1, "keyword_autoreply_info": map[string]any{"list": []any{map[string]any{"rule_name": "r1"}}}})
	}))
	defer srv.Close()

	info, err := newClientWithServer(t, srv).GetAutoreply(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if info.IsAddFriendReplyOpen != 1 || len(info.KeywordAutoreplyInfo.List) != 1 {
		t.Errorf("info = %+v", info)
	}
}

// ---- 标签扩展 ----

// TestTagGetFans 验证标签下粉丝列表解析与请求体。
func TestTagGetFans(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/user/tag/get" {
			t.Errorf("path = %q", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"count": 1, "data": map[string]any{"openid": []string{"o1"}}, "next_openid": "o_next"})
	}))
	defer srv.Close()

	res, err := newClientWithServer(t, srv).TagGetFans(context.Background(), 100, "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if gotBody["tagid"].(float64) != 100 {
		t.Errorf("tagid = %v", gotBody["tagid"])
	}
	if res.Count != 1 || len(res.Data.OpenID) != 1 || res.NextOpenID != "o_next" {
		t.Errorf("res = %+v", res)
	}
}

// TestBatchTagging 验证批量打标签请求体。
func TestBatchTagging(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/tags/members/batchtagging" {
			t.Errorf("path = %q", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"errcode": 0})
	}))
	defer srv.Close()

	err := newClientWithServer(t, srv).BatchTagging(context.Background(), []string{"o1", "o2"}, 7)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if gotBody["tagid"].(float64) != 7 {
		t.Errorf("tagid = %v", gotBody["tagid"])
	}
	if arr, _ := gotBody["openid_list"].([]any); len(arr) != 2 {
		t.Errorf("openid_list = %v", gotBody["openid_list"])
	}
}

// ---- 草稿 ----

// TestDraftAdd 验证新增草稿请求体与 media_id 解析。
func TestDraftAdd(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/draft/add" {
			t.Errorf("path = %q", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"media_id": "m_draft"})
	}))
	defer srv.Close()

	mid, err := newClientWithServer(t, srv).DraftAdd(context.Background(), []map[string]any{{"title": "t"}})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if gotBody["articles"] == nil {
		t.Errorf("body = %v", gotBody)
	}
	if mid != "m_draft" {
		t.Errorf("mid = %q", mid)
	}
}

// TestDraftBatchGet 验证草稿列表解析。
func TestDraftBatchGet(t *testing.T) {
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/draft/batchget" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"total_count": 1, "item_count": 1, "item": []map[string]any{{"media_id": "m1", "update_time": int64(123)}}})
	}))
	defer srv.Close()

	res, err := newClientWithServer(t, srv).DraftBatchGet(context.Background(), 0, 20)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if res.TotalCount != 1 || len(res.Item) != 1 || res.Item[0].MediaID != "m1" {
		t.Errorf("res = %+v", res)
	}
}

// ---- 发布 ----

// TestPublishSubmit 验证发布草稿请求体与 publish_id 解析。
func TestPublishSubmit(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/freepublish/submit" {
			t.Errorf("path = %q", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"publish_id": "100"})
	}))
	defer srv.Close()

	pid, err := newClientWithServer(t, srv).PublishSubmit(context.Background(), "m_draft")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if gotBody["media_id"] != "m_draft" {
		t.Errorf("media_id = %v", gotBody["media_id"])
	}
	if pid != "100" {
		t.Errorf("pid = %q", pid)
	}
}

// TestPublishGet 验证发布状态查询解析。
func TestPublishGet(t *testing.T) {
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/freepublish/get" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"publish_id": 100, "publish_status": 0})
	}))
	defer srv.Close()

	st, err := newClientWithServer(t, srv).PublishGet(context.Background(), "100")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if st.PublishState != 0 || st.PublishID != 100 {
		t.Errorf("st = %+v", st)
	}
}
