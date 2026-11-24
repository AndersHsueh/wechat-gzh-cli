package tests

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestAddKfAccount 验证添加客服账号的请求体与路径。
func TestAddKfAccount(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/customservice/kfaccount/add" {
			t.Errorf("path = %q", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"errcode": 0})
	}))
	defer srv.Close()

	err := newClientWithServer(t, srv).AddKfAccount(context.Background(), "test1@test", "客服1")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if gotBody["kf_account"] != "test1@test" || gotBody["nickname"] != "客服1" {
		t.Errorf("body = %v", gotBody)
	}
}

// TestGetKfList 验证客服列表解析。
func TestGetKfList(t *testing.T) {
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/customservice/getkflist" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"kf_list": []map[string]any{{"kf_account": "a@b", "kf_nick": "n", "kf_id": 1}}})
	}))
	defer srv.Close()

	res, err := newClientWithServer(t, srv).GetKfList(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(res.KFList) != 1 || res.KFList[0].KfAccount != "a@b" || res.KFList[0].KfID != 1 {
		t.Errorf("res = %+v", res)
	}
}

// TestCreateKfSession 验证创建会话请求体。
func TestCreateKfSession(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/customservice/kfsession/create" {
			t.Errorf("path = %q", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"errcode": 0})
	}))
	defer srv.Close()

	err := newClientWithServer(t, srv).CreateKfSession(context.Background(), "k@b", "o_1")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if gotBody["openid"] != "o_1" || gotBody["kf_account"] != "k@b" {
		t.Errorf("body = %v", gotBody)
	}
}

// TestGetKfSession 验证会话状态查询带 openid query。
func TestGetKfSession(t *testing.T) {
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/customservice/kfsession/getsession" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("openid") != "o_1" {
			t.Errorf("openid = %q", r.URL.Query().Get("openid"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"kf_account": "k@b", "createtime": int64(123)})
	}))
	defer srv.Close()

	res, err := newClientWithServer(t, srv).GetKfSession(context.Background(), "o_1")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if res.KfAccount != "k@b" {
		t.Errorf("res = %+v", res)
	}
}

// TestGetMsgList 验证聊天记录请求体与解析。
func TestGetMsgList(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/customservice/msgrecord/getmsglist" {
			t.Errorf("path = %q", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"recordlist": []map[string]any{{"openid": "o_1"}}})
	}))
	defer srv.Close()

	list, err := newClientWithServer(t, srv).GetMsgList(context.Background(), 100, 200, 0, 10000)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if gotBody["number"].(float64) != 10000 {
		t.Errorf("number = %v", gotBody["number"])
	}
	if len(list) != 1 {
		t.Errorf("list = %v", list)
	}
}

// TestSendTyping 验证客服输入状态请求体。
func TestSendTyping(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/message/custom/typing" {
			t.Errorf("path = %q", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"errcode": 0})
	}))
	defer srv.Close()

	err := newClientWithServer(t, srv).SendTyping(context.Background(), "o_1", "Typing")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if gotBody["command"] != "Typing" || gotBody["touser"] != "o_1" {
		t.Errorf("body = %v", gotBody)
	}
}

// TestGetWaitCase 验证未接入会话列表解析。
func TestGetWaitCase(t *testing.T) {
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/customservice/kfsession/getwaitcase" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"waitcaselist": []map[string]any{{"openid": "o_1"}}})
	}))
	defer srv.Close()

	list, err := newClientWithServer(t, srv).GetWaitCase(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(list) != 1 || list[0].OpenID != "o_1" {
		t.Errorf("list = %v", list)
	}
}

// TestBatchUserInfo 验证批量用户信息请求体与解析。
func TestBatchUserInfo(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/user/info/batchget" {
			t.Errorf("path = %q", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"user_info_list": []map[string]any{{"openid": "o_1", "subscribe": 1}}})
	}))
	defer srv.Close()

	list, err := newClientWithServer(t, srv).BatchUserInfo(context.Background(), []map[string]any{{"openid": "o_1", "lang": "zh_CN"}})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if gotBody["user_list"] == nil {
		t.Errorf("body = %v", gotBody)
	}
	if len(list) != 1 || list[0].OpenID != "o_1" {
		t.Errorf("list = %+v", list)
	}
}

// TestGetBlacklist 验证黑名单列表解析。
func TestGetBlacklist(t *testing.T) {
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/tags/members/getblacklist" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"openid": []string{"o_1"}}, "next_openid": "o_next"})
	}))
	defer srv.Close()

	ids, next, err := newClientWithServer(t, srv).GetBlacklist(context.Background(), "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(ids) != 1 || next != "o_next" {
		t.Errorf("ids=%v next=%q", ids, next)
	}
}

// TestBatchBlacklist 验证拉黑请求体。
func TestBatchBlacklist(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/tags/members/batchblacklist" {
			t.Errorf("path = %q", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"errcode": 0})
	}))
	defer srv.Close()

	err := newClientWithServer(t, srv).BatchBlacklist(context.Background(), []string{"o_1"})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if arr, _ := gotBody["openid_list"].([]any); len(arr) != 1 {
		t.Errorf("openid_list = %v", gotBody["openid_list"])
	}
}

// TestChangeOpenID 验证 openid 转换请求体与解析。
func TestChangeOpenID(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/changeopenid" {
			t.Errorf("path = %q", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"result_list": []map[string]any{{"ori_openid": "o_old", "new_openid": "o_new"}}})
	}))
	defer srv.Close()

	list, err := newClientWithServer(t, srv).ChangeOpenID(context.Background(), []string{"o_old"})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(list) != 1 || list[0].NewOpenID != "o_new" {
		t.Errorf("list = %+v", list)
	}
}

// TestDatacubeQuery 验证数据统计请求体与解析。
func TestDatacubeQuery(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/datacube/getusersummary" {
			t.Errorf("path = %q", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"list": []map[string]any{{"ref_date": "2026-11-01", "new_user": 5}}})
	}))
	defer srv.Close()

	list, err := newClientWithServer(t, srv).DatacubeQuery(context.Background(), "/datacube/getusersummary", "2026-11-01", "2026-11-02")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if gotBody["begin_date"] != "2026-11-01" || gotBody["end_date"] != "2026-11-02" {
		t.Errorf("body = %v", gotBody)
	}
	if len(list) != 1 {
		t.Errorf("list = %v", list)
	}
}

// TestOpenComment 验证打开评论请求体。
func TestOpenComment(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/comment/open" {
			t.Errorf("path = %q", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"errcode": 0})
	}))
	defer srv.Close()

	err := newClientWithServer(t, srv).OpenComment(context.Background(), 100, 0)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if gotBody["msg_data_id"].(float64) != 100 {
		t.Errorf("body = %v", gotBody)
	}
}

// TestListComment 验证评论列表解析。
func TestListComment(t *testing.T) {
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/comment/list" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"total": 1, "comment": []map[string]any{{"user_comment_id": 5, "content": "hi"}}})
	}))
	defer srv.Close()

	res, err := newClientWithServer(t, srv).ListComment(context.Background(), 100, 0, 0, 20, 0)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if res.Total != 1 || len(res.Comment) != 1 || res.Comment[0].Content != "hi" {
		t.Errorf("res = %+v", res)
	}
}

// TestMarkElect 验证评论精选请求体。
func TestMarkElect(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/comment/markelect" {
			t.Errorf("path = %q", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"errcode": 0})
	}))
	defer srv.Close()

	err := newClientWithServer(t, srv).MarkElect(context.Background(), 5)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if gotBody["comment_id"].(float64) != 5 {
		t.Errorf("body = %v", gotBody)
	}
}

// TestReplyComment 验证回复请求体。
func TestReplyComment(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/comment/reply/add" {
			t.Errorf("path = %q", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"errcode": 0})
	}))
	defer srv.Close()

	err := newClientWithServer(t, srv).ReplyComment(context.Background(), 5, "谢谢")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if gotBody["content"] != "谢谢" {
		t.Errorf("body = %v", gotBody)
	}
}

// TestGetJSSDKTicket 验证 jsapi 票据解析与 type=jsapi query。
func TestGetJSSDKTicket(t *testing.T) {
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cgi-bin/ticket/getticket" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("type") != "jsapi" {
			t.Errorf("type = %q", r.URL.Query().Get("type"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"ticket": "TK", "expires_in": 7200})
	}))
	defer srv.Close()

	tk, err := newClientWithServer(t, srv).GetJSSDKTicket(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if tk.Ticket != "TK" || tk.ExpiresIn != 7200 {
		t.Errorf("tk = %+v", tk)
	}
}

// TestSendChannelMsg 验证就医助手推送请求体与路径。
func TestSendChannelMsg(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(mockTokenHandler(t, "T", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cityservice/sendchannelmsg" {
			t.Errorf("path = %q", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"errcode": 0})
	}))
	defer srv.Close()

	err := newClientWithServer(t, srv).SendChannelMsg(context.Background(), map[string]any{"touser": "o_1"})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if gotBody["touser"] != "o_1" {
		t.Errorf("body = %v", gotBody)
	}
}
