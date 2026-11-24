# Changelog

`wechat-gzh-cli` 的版本变更记录。版本规则:从 `v0.0.1` 起,仅在最小位(末位)递增。

## [v0.0.5] - 2026-11-25

### 质量加固(simplify 全量简化)

- 全库 `staticcheck` + `go vet` 零警告:无死代码、无未使用变量。
- `store` 命令组:内联一次性辅助函数(`preCondForStore`/`sideEffectForStore`),移除为单次调用而抽的抽象。
- `hospital-msg`:`Use` 与 `CommandPath`/`--agent-help` 注册路径不一致的 bug 修复,重构为标准父(`hospital-msg`) + 子(`send`) 结构。
- 新增 `tests/agenthelp_reach_test.go`:全量可达性回归测试,遍历全部 130 个叶子命令,断言无环境变量下 `--agent-help` 文档可见,并注册进 `TestVerifyAll`,固化三条可达性法则。

## [v0.0.4] - 2026-11-24

### 新增(def 约 26 个接口)

- **智能接口 `ai` 命令组**:`translate`(微信翻译)、`voice-to-text`(语音转写,multipart 上传)、`voice-result`(识别结果)、`ocr-menu/comm/driving/bankcard/bizlicense/driverlicense/idcard`(7 种 OCR,基于 img_url)、`img-aicrop`/`img-qrcode`(图像处理)。
- **微信门店 `store` 命令组**:`shop-category/apply/audit/modify/district/map-poi`(门店小程序)、`store-add/get/list/del/update/map-create`(门店 POI)。
- **客服头像**:`customer-msg account-upload-headimg`(multipart 上传)。
- **素材**:`material get-jssdk-voice`(高清语音素材)。
- 废弃的草稿箱开关仅文档登记,不实现;本轮无不实现接口,不产生 prd-5。

## [v0.0.3] - 2026-11-24

### 新增(约 51 个接口)

- **客服扩展**:`customer-msg account-add/update/del/list`、`online-list`、`inviteworker`、`session-create/close/get/list/waiting`、`msg-record`、`typing`。
- **用户扩展**:`user batch-info`、`remark`、`blacklist-get`、`blacklist-add`、`blacklist-del`、`change-openid`。
- **数据统计**:`datacube` 命令组(21 个子命令:用户/图文/消息/接口分析)。
- **JS-SDK**:`jsapi-ticket`。
- **就医助手**:`hospital-msg send`。
- **留言管理**:`comment` 命令组(8 个:`open/close/list/markelect/unmarkelect/delete/reply-add/reply-del`)。

## [v0.0.2] - 2026-11-24

### 新增(36 个接口)

- **基础**:`stable-token`、`callback-check`、`push-server-ip`。
- **openApi**:`openapi reset`、`openapi quota`、`openapi rid`、`openapi reset-api`。
- **菜单扩展**:`menu self`、`add-conditional`、`del-conditional`、`try-match`。
- **消息扩展**:`message mass-delete`、`speed-get`、`speed-set`、`preview`、`subscribe send`、`autoreply get`。
- **素材扩展**:`material count`、`get-temp`。
- **草稿**:`draft add/update/list/count/delete/get`。
- **发布**:`publish submit/get/list/delete/get-article`。
- **标签扩展**:`tag getfans/update/delete/batch-tag/batch-untag/of-user`。

## [v0.0.1] - 2026-11-24

### 新增(核心 MVP,18 个命令)

- **鉴权**:`access-token`、`server-ip`。
- **素材**:`material upload/list/get/delete/upload-temp/upload-img`。
- **消息**:`message mass-send/mass-status`。
- **用户**:`user info/followers`。
- **标签**:`tag list/create`。
- **菜单**:`menu create/list/delete`。
- **客服**:`customer-msg send`。

### 基础设施

- Cobra 双界面 CLI:`--help`(人) + `--agent-help`(11 维度,供 AI Agent),满足三条可达性法则。
- 从环境变量 `WX_GZ_AppID`/`WX_GZ_AppSecret` 鉴权,access_token 7200s 缓存。
- TDD:各接口独立测试单元 + `TestVerifyAll` 全量回归。
