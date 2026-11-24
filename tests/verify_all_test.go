package tests

import "testing"

// TestVerifyAll 是总测试:逐个调起各 API 的独立测试单元,做全量回归。
// 每新增一个 API client 测试,只需在此注册即可被 verify_all 全量覆盖。
func TestVerifyAll(t *testing.T) {
	units := map[string]func(*testing.T){
		// 鉴权
		"access_token":              TestClientAccessToken,
		"access_token_errcode":      TestClientAccessTokenErrcode,
		"access_token_cached":       TestClientAccessTokenCached,
		// 基础
		"get_api_domain_ip":         TestGetApiDomainIP,
		"get_material_count":        TestGetMaterialCount,
		"get_temp_media":            TestGetTempMedia,
		"stable_access_token":       TestStableAccessToken,
		"get_api_quota":             TestGetApiQuota,
		"get_callback_ip":           TestGetCallbackIP,
		// 素材
		"batch_get_material":        TestBatchGetMaterial,
		"upload_image":              TestUploadImage,
		// 客服
		"send_customer_message":     TestSendCustomerMessage,
		// 标签
		"list_tags":                 TestListTags,
		"tag_get_fans":              TestTagGetFans,
		"batch_tagging":             TestBatchTagging,
		// 菜单
		"add_conditional_menu":      TestAddConditionalMenu,
		// 消息
		"subscribe_send":            TestSubscribeSend,
		"get_autoreply":             TestGetAutoreply,
		// 草稿
		"draft_add":                 TestDraftAdd,
		"draft_batch_get":           TestDraftBatchGet,
		// 发布
		"publish_submit":            TestPublishSubmit,
		"publish_get":               TestPublishGet,
	}

	for name, fn := range units {
		t.Run(name, fn)
	}
}
