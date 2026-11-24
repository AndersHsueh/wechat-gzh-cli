package tests

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestTranslateContent 验证翻译请求 path/query 与解析。
func TestTranslateContent(t *testing.T) {
	var gotQuery string
	var gotBody map[string]any
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/media/voice/translatecontent" {
			t.Errorf("path = %q", r.URL.Path)
		}
		gotQuery = r.URL.RawQuery
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"tocontent": "Hello"})
	}))
	defer srv.Close()

	out, err := newClientWithServer(t, srv).TranslateContent(context.Background(), "zh_CN", "en_US", "你好")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !strings.Contains(gotQuery, "lfrom=zh_CN") || !strings.Contains(gotQuery, "lto=en_US") {
		t.Errorf("query = %q", gotQuery)
	}
	if gotBody["content"] != "你好" {
		t.Errorf("body = %v", gotBody)
	}
	if out != "Hello" {
		t.Errorf("out = %q", out)
	}
}

// TestOcrQuery 验证 OCR 请求带 img_url query 与解析。
func TestOcrQuery(t *testing.T) {
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cv/ocr/idcard" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("img_url") != "https://example.com/a.jpg" {
			t.Errorf("img_url = %q", r.URL.Query().Get("img_url"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"name": "张三"})
	}))
	defer srv.Close()

	res, err := newClientWithServer(t, srv).OcrQuery(context.Background(), "/cv/ocr/idcard", "https://example.com/a.jpg")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if res["name"] != "张三" {
		t.Errorf("res = %v", res)
	}
}

// TestVoiceToText 验证语音转写使用 multipart 上传。
func TestVoiceToText(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "voice.mp3")
	if err := os.WriteFile(fp, []byte("audio-data"), 0o644); err != nil {
		t.Fatal(err)
	}
	var gotQuery string
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/media/voice/addvoicetorecofortext" {
			t.Errorf("path = %q", r.URL.Path)
		}
		gotQuery = r.URL.RawQuery
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
			t.Errorf("content-type = %q", r.Header.Get("Content-Type"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"errcode": 0})
	}))
	defer srv.Close()

	err := newClientWithServer(t, srv).VoiceToText(context.Background(), "mp3", "v1", "zh_CN", fp)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !strings.Contains(gotQuery, "voice_id=v1") || !strings.Contains(gotQuery, "format=mp3") {
		t.Errorf("query = %q", gotQuery)
	}
}

// TestVoiceResult 验证识别结果查询带 voice_id query。
func TestVoiceResult(t *testing.T) {
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/media/voice/queryrecoresultfortext" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("voice_id") != "v1" {
			t.Errorf("voice_id = %q", r.URL.Query().Get("voice_id"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"result": "你好"})
	}))
	defer srv.Close()

	res, err := newClientWithServer(t, srv).VoiceResult(context.Background(), "v1", "zh_CN")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if res["result"] != "你好" {
		t.Errorf("res = %v", res)
	}
}

// TestGetJSSDKRawVoice 验证高清语音素材下载带 media_id query。
func TestGetJSSDKRawVoice(t *testing.T) {
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/media/get/jssdk" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("media_id") != "m_voice" {
			t.Errorf("media_id = %q", r.URL.Query().Get("media_id"))
		}
		w.Write([]byte("audio-raw"))
	}))
	defer srv.Close()

	data, err := newClientWithServer(t, srv).GetJSSDKRawVoice(context.Background(), "m_voice")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if string(data) != "audio-raw" {
		t.Errorf("data = %q", data)
	}
}

// TestWxaQuery 验证门店接口请求体与路径。
func TestWxaQuery(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/wxa/get_merchant_category" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"category": "餐饮", "errcode": 0})
	}))
	defer srv.Close()

	res, err := newClientWithServer(t, srv).WxaQuery(context.Background(), "/wxa/get_merchant_category", nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if res["category"] != "餐饮" {
		t.Errorf("res = %v", res)
	}
	_ = gotBody
}

// TestUploadKfHeadimg 验证客服头像上传(multipart + kf_account query)。
func TestUploadKfHeadimg(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "avatar.jpg")
	if err := os.WriteFile(fp, []byte("img"), 0o644); err != nil {
		t.Fatal(err)
	}
	var gotQuery string
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/customservice/kfaccount/uploadheadimg" {
			t.Errorf("path = %q", r.URL.Path)
		}
		gotQuery = r.URL.RawQuery
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
			t.Errorf("content-type = %q", r.Header.Get("Content-Type"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"errcode": 0})
	}))
	defer srv.Close()

	err := newClientWithServer(t, srv).UploadKfHeadimg(context.Background(), "test1@test", fp)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !strings.Contains(gotQuery, "kf_account=test1%40test") {
		t.Errorf("query = %q", gotQuery)
	}
}
