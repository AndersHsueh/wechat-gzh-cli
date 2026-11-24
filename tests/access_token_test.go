package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"tx-gongzhonghao-cli/internal/client"
)

// 模拟微信服务器:处理 GET /cgi-bin/token,返回 access_token。
func mockAccessTokenServer(t *testing.T, token string, errcode int, errmsg string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/token" {
			t.Errorf("unexpected path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		q := r.URL.Query()
		// 校验必传参数
		if q.Get("grant_type") != "client_credential" {
			t.Errorf("grant_type=%q want client_credential", q.Get("grant_type"))
		}
		if q.Get("appid") == "" || q.Get("secret") == "" {
			t.Errorf("appid/secret missing: %v", q)
		}
		w.Header().Set("Content-Type", "application/json")
		if errcode != 0 {
			json.NewEncoder(w).Encode(map[string]any{"errcode": errcode, "errmsg": errmsg})
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"access_token": token, "expires_in": 7200})
	}))
}

// 测试:AccessToken 正确解析返回的 token 与有效期(先红后绿)。
func TestClientAccessToken(t *testing.T) {
	token := "ACCESS_TOKEN_ABC"
	srv := mockAccessTokenServer(t, token, 0, "")
	defer srv.Close()

	c := client.New(client.Config{
		AppID:    "wx_appid",
		AppSecret: "secret",
		BaseURL:   srv.URL,
	})

	got, err := c.AccessToken(context.Background())
	if err != nil {
		t.Fatalf("AccessToken() error = %v", err)
	}
	if got != token {
		t.Errorf("AccessToken() = %q, want %q", got, token)
	}
}

// 测试:微信返回错误码时,AccessToken 返回错误(不把成功当空串)。
func TestClientAccessTokenErrcode(t *testing.T) {
	srv := mockAccessTokenServer(t, "", 40013, "invalid appid")
	defer srv.Close()

	c := client.New(client.Config{
		AppID:    "wx_appid",
		AppSecret: "secret",
		BaseURL:   srv.URL,
	})

	if _, err := c.AccessToken(context.Background()); err == nil {
		t.Fatal("AccessToken() expected error for errcode 40013, got nil")
	}
}

// 测试:access_token 在有效期内缓存,不全新的 HTTP 请求。
func TestClientAccessTokenCached(t *testing.T) {
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"access_token": "CACHED", "expires_in": 7200})
	}))
	defer srv.Close()

	c := client.New(client.Config{AppID: "a", AppSecret: "s", BaseURL: srv.URL})
	ctx := context.Background()
	if _, err := c.AccessToken(ctx); err != nil {
		t.Fatalf("first token err: %v", err)
	}
	if _, err := c.AccessToken(ctx); err != nil {
		t.Fatalf("second token err: %v", err)
	}
	if callCount != 1 {
		t.Errorf("expected token cached (1 request), got %d", callCount)
	}
}
