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
		// 素材
		"batch_get_material":        TestBatchGetMaterial,
		"upload_image":              TestUploadImage,
		// 客服
		"send_customer_message":     TestSendCustomerMessage,
		// 标签
		"list_tags":                 TestListTags,
	}

	for name, fn := range units {
		t.Run(name, fn)
	}
}
