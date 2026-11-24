package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"tx-gongzhonghao-cli/internal/client"
)

// newClientWithServer 构造指向 mock 服务器的 client。
func newClientWithServer(t *testing.T, srv *httptest.Server) *client.Client {
	t.Helper()
	return client.New(client.Config{AppID: "a", AppSecret: "s", BaseURL: srv.URL})
}

// mockTokenHandler 返回一个会先应答 token,再把后续请求路由到 next 的 handler。
func mockTokenHandler(t *testing.T, token string, next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/cgi-bin/token" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"access_token": token, "expires_in": 7200})
			return
		}
		next(w, r)
	})
}

// TestGetApiDomainIP 验证 GET /cgi-bin/get_api_domain_ip 的请求与解析。
func TestGetApiDomainIP(t *testing.T) {
	srv := httptest.NewServer(mockTokenHandler(t, "TOKEN", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/get_api_domain_ip" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if q := r.URL.Query().Get("access_token"); q != "TOKEN" {
			t.Errorf("access_token = %q", q)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"ip_list": []string{"1.1.1.1", "2.2.2.2"}})
	}))
	defer srv.Close()

	ips, err := newClientWithServer(t, srv).GetApiDomainIP(context.Background())
	if err != nil {
		t.Fatalf("GetApiDomainIP err: %v", err)
	}
	if len(ips) != 2 || ips[0] != "1.1.1.1" {
		t.Errorf("ips = %v", ips)
	}
}

// TestBatchGetMaterial 验证素材列表的请求体与解析。
func TestBatchGetMaterial(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/material/batchget_material" {
			t.Errorf("path = %q", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"total_count": 1, "item_count": 1,
			"item": []map[string]any{{"media_id": "m1", "name": "a.jpg"}},
		})
	}))
	defer srv.Close()

	res, err := newClientWithServer(t, srv).BatchGetMaterial(context.Background(), "image", 0, 20)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if gotBody["type"] != "image" {
		t.Errorf("type = %v, body = %v", gotBody["type"], gotBody)
	}
	if off, _ := gotBody["offset"].(float64); off != 0 {
		t.Errorf("offset = %v", gotBody["offset"])
	}
	if cnt, _ := gotBody["count"].(float64); cnt != 20 {
		t.Errorf("count = %v", gotBody["count"])
	}
	if res.TotalCount != 1 || len(res.Item) != 1 || res.Item[0].MediaID != "m1" {
		t.Errorf("res = %+v", res)
	}
}

// TestUploadImage 验证正文图上传(multipart)返回 URL。
func TestUploadImage(t *testing.T) {
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/media/uploadimg" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
			t.Errorf("content-type = %q", r.Header.Get("Content-Type"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"url": "https://mmbiz/a.jpg"})
	}))
	defer srv.Close()

	u, err := newClientWithServer(t, srv).UploadImage(context.Background(), "a.jpg", bytes.NewBufferString("data"))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if u != "https://mmbiz/a.jpg" {
		t.Errorf("url = %q", u)
	}
}

// TestSendCustomerMessage 验证客服消息请求体包含 touser/msgtype/content。
func TestSendCustomerMessage(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/message/custom/send" {
			t.Errorf("path = %q", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"errcode": 0, "errmsg": "ok"})
	}))
	defer srv.Close()

	err := newClientWithServer(t, srv).SendCustomerMessage(context.Background(), "o_1", "text", map[string]any{"content": "hi"})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if gotBody["touser"] != "o_1" || gotBody["msgtype"] != "text" {
		t.Errorf("body = %v", gotBody)
	}
	text := gotBody["text"].(map[string]any)
	if text["content"] != "hi" {
		t.Errorf("text = %v", text)
	}
}

// TestListTags 验证标签列表解析。
func TestListTags(t *testing.T) {
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/tags/get" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"tags": []map[string]any{{"id": 1, "name": "vip", "count": 5}}})
	}))
	defer srv.Close()

	tags, err := newClientWithServer(t, srv).ListTags(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(tags) != 1 || tags[0].ID != 1 || tags[0].Name != "vip" {
		t.Errorf("tags = %+v", tags)
	}
}
