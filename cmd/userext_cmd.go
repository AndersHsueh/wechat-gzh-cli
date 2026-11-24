package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"tx-gongzhonghao-cli/internal/client"
)

// newUserExtCmds 用户信息扩展子命令。
func newUserExtCmds() []*cobra.Command {
	var cmds []*cobra.Command

	// batch-info
	{
		var jsonOut bool
		cmd := &cobra.Command{
			Use:   "batch-info --openids o1,o2 [--lang zh_CN]",
			Short: "批量获取用户基本信息",
			Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
			RunE: func(c *cobra.Command, _ []string) error {
				if agentHelpShortcut(c) {
					return nil
				}
				ids, _ := c.Flags().GetString("openids")
				lang, _ := c.Flags().GetString("lang")
				if ids == "" {
					return fmt.Errorf("--openids 必填(逗号分隔)")
				}
				if lang == "" {
					lang = "zh_CN"
				}
				var userList []map[string]any
				for _, id := range strings.Split(ids, ",") {
					if id = strings.TrimSpace(id); id != "" {
						userList = append(userList, map[string]any{"openid": id, "lang": lang})
					}
				}
				if len(userList) == 0 {
					return fmt.Errorf("无有效 openid")
				}
				if len(userList) > 100 {
					userList = userList[:100]
				}
				return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
					list, err := cl.BatchUserInfo(ctx, userList)
					if err != nil {
						return err
					}
					if jsonOut {
						return writeJSON(map[string]any{"user_info_list": list})
					}
					for _, u := range list {
						fmt.Fprintf(c.OutOrStdout(), "%s\t%s\t%d\n", u.OpenID, u.Nickname, u.Subscribe)
					}
					return nil
				})
			},
		}
		cmd.Flags().String("openids", "", "逗号分隔的 openid(最多 100)")
		cmd.Flags().String("lang", "zh_CN", "语言")
		cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
		addAgentHelpFlag(cmd)
		registerAgentHelp("wechat-gzh-cli user batch-info", RenderAgentHelp(AgentHelpDoc{
			Command:        "wechat-gzh-cli user batch-info --openids o1,o2 --lang zh_CN",
			Summary:        "批量获取用户基本信息(最多 100 条)。",
			PreConditions:  []string{"已配置鉴权。", "--openids 逗号分隔,≤100。"},
			PostConditions: []string{"stdout 输出 openid\\tnickname\\tsubscribe。"},
			Idempotency:    "YES. 只读。",
			Atomicity:      "Atomic.",
			SideEffects:    []string{"无。"},
			CommonErrors:   []ErrSpec{{Code: "40003", Desc: "invalid openid", Fix: "检查 openid。"}},
			Scope:          "开发者权限。",
			RateLimit:      "常规。",
			RaceConditions: []string{"无。"},
			Orchestration:  "批量做用户画像/定向。",
			ParsingHints:   "--json 输出 user_info_list。",
			Examples:       []string{"wechat-gzh-cli user batch-info --openids o1,o2 --json"},
		}))
		cmds = append(cmds, cmd)
	}

	// remark
	{
		cmd := &cobra.Command{
			Use:   "remark --openid o_xxx --remark 备注",
			Short: "设置用户备注名",
			Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
			RunE: func(c *cobra.Command, _ []string) error {
				if agentHelpShortcut(c) {
					return nil
				}
				id, _ := c.Flags().GetString("openid")
				remark, _ := c.Flags().GetString("remark")
				if id == "" || remark == "" {
					return fmt.Errorf("--openid 与 --remark 必填")
				}
				return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
					if err := cl.UpdateRemark(ctx, id, remark); err != nil {
						return err
					}
					fmt.Fprintln(c.OutOrStdout(), "ok")
					return nil
				})
			},
		}
		cmd.Flags().String("openid", "", "用户 openid")
		cmd.Flags().String("remark", "", "备注名(最长 30 字符)")
		addAgentHelpFlag(cmd)
		registerAgentHelp("wechat-gzh-cli user remark", RenderAgentHelp(AgentHelpDoc{
			Command:        "wechat-gzh-cli user remark --openid o_xxx --remark 金卡会员",
			Summary:        "设置用户备注名(仅微信认证的服务号开放)。",
			PreConditions:  []string{"已配置鉴权。", "认证服务号。", "--openid/--remark 必填。"},
			PostConditions: []string{"stdout 输出 ok;备注已保存。"},
			Idempotency:    "NO.",
			Atomicity:      "Atomic.",
			SideEffects:    []string{"用户备注被修改。"},
			CommonErrors:   []ErrSpec{},
			Scope:          "开发者权限(认证服务号)。",
			RateLimit:      "常规。",
			RaceConditions: []string{"无。"},
			Orchestration:  "用于人工打标/运营备注。",
			ParsingHints:   "文本输出 ok。",
			Examples:       []string{"wechat-gzh-cli user remark --openid o_xxx --remark 金卡会员"},
		}))
		cmds = append(cmds, cmd)
	}

	return cmds
}

// newUserBlacklistCmds 黑名单管理子命令。
func newUserBlacklistCmds() []*cobra.Command {
	var cmds []*cobra.Command

	// blacklist-get
	{
		var jsonOut bool
		cmd := &cobra.Command{
			Use:   "blacklist-get [--begin-openid xxx]",
			Short: "获取用户黑名单",
			Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
			RunE: func(c *cobra.Command, _ []string) error {
				if agentHelpShortcut(c) {
					return nil
				}
				begin, _ := c.Flags().GetString("begin-openid")
				return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
					ids, next, err := cl.GetBlacklist(ctx, begin)
					if err != nil {
						return err
					}
					if jsonOut {
						return writeJSON(map[string]any{"openid_list": ids, "next_openid": next})
					}
					for _, id := range ids {
						fmt.Fprintln(c.OutOrStdout(), id)
					}
					if next != "" {
						fmt.Fprintf(c.OutOrStdout(), "next_openid=%s\n", next)
					}
					return nil
				})
			},
		}
		cmd.Flags().String("begin-openid", "", "翻页游标")
		cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
		addAgentHelpFlag(cmd)
		registerAgentHelp("wechat-gzh-cli user blacklist-get", RenderAgentHelp(AgentHelpDoc{
			Command:        "wechat-gzh-cli user blacklist-get --begin-openid xxx",
			Summary:        "获取公众号黑名单用户 openid 列表。",
			PreConditions:  []string{"已配置鉴权。"},
			PostConditions: []string{"stdout 输出 openid(逐行)。"},
			Idempotency:    "YES. 只读。",
			Atomicity:      "Atomic.",
			SideEffects:    []string{"无。"},
			CommonErrors:   []ErrSpec{},
			Scope:          "开发者权限。",
			RateLimit:      "常规。",
			RaceConditions: []string{"分页快照。"},
			Orchestration:  "用 next_openid 翻页拉全量。",
			ParsingHints:   "--json 输出 openid_list/next_openid。",
			Examples:       []string{"wechat-gzh-cli user blacklist-get --json"},
		}))
		cmds = append(cmds, cmd)
	}

	// blacklist-add
	{
		cmd := &cobra.Command{
			Use:   "blacklist-add --openids o1,o2",
			Short: "拉黑用户",
			Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
			RunE: func(c *cobra.Command, _ []string) error {
				if agentHelpShortcut(c) {
					return nil
				}
				ids, _ := c.Flags().GetString("openids")
				openids, err := parseOpenIDs(ids)
				if err != nil {
					return err
				}
				return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
					if err := cl.BatchBlacklist(ctx, openids); err != nil {
						return err
					}
					fmt.Fprintln(c.OutOrStdout(), "ok")
					return nil
				})
			},
		}
		cmd.Flags().String("openids", "", "逗号分隔的 openid")
		addAgentHelpFlag(cmd)
		registerAgentHelp("wechat-gzh-cli user blacklist-add", RenderAgentHelp(AgentHelpDoc{
			Command:        "wechat-gzh-cli user blacklist-add --openids o1,o2",
			Summary:        "将用户加入黑名单(不再收到消息)。",
			PreConditions:  []string{"已配置鉴权。", "--openids 必填。"},
			PostConditions: []string{"stdout 输出 ok;用户被拉黑。"},
			Idempotency:    "YES.",
			Atomicity:      "Atomic.",
			SideEffects:    []string{"用户被拉黑。"},
			CommonErrors:   []ErrSpec{{Code: "40003", Desc: "invalid openid", Fix: "检查 openid。"}},
			Scope:          "开发者权限。",
			RateLimit:      "常规。",
			RaceConditions: []string{"无。"},
			Orchestration:  "拉黑前用 blacklist-get 确认。",
			ParsingHints:   "文本输出 ok。",
			Examples:       []string{"wechat-gzh-cli user blacklist-add --openids o1,o2"},
		}))
		cmds = append(cmds, cmd)
	}

	// blacklist-del
	{
		cmd := &cobra.Command{
			Use:   "blacklist-del --openids o1,o2",
			Short: "取消拉黑用户",
			Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
			RunE: func(c *cobra.Command, _ []string) error {
				if agentHelpShortcut(c) {
					return nil
				}
				ids, _ := c.Flags().GetString("openids")
				openids, err := parseOpenIDs(ids)
				if err != nil {
					return err
				}
				return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
					if err := cl.BatchUnblacklist(ctx, openids); err != nil {
						return err
					}
					fmt.Fprintln(c.OutOrStdout(), "ok")
					return nil
				})
			},
		}
		cmd.Flags().String("openids", "", "逗号分隔的 openid")
		addAgentHelpFlag(cmd)
		registerAgentHelp("wechat-gzh-cli user blacklist-del", RenderAgentHelp(AgentHelpDoc{
			Command:        "wechat-gzh-cli user blacklist-del --openids o1,o2",
			Summary:        "将用户移出黑名单。",
			PreConditions:  []string{"已配置鉴权。", "--openids 必填。"},
			PostConditions: []string{"stdout 输出 ok;用户解除拉黑。"},
			Idempotency:    "YES.",
			Atomicity:      "Atomic.",
			SideEffects:    []string{"用户被移出黑名单。"},
			CommonErrors:   []ErrSpec{},
			Scope:          "开发者权限。",
			RateLimit:      "常规。",
			RaceConditions: []string{"无。"},
			Orchestration:  "解除误拉黑。",
			ParsingHints:   "文本输出 ok。",
			Examples:       []string{"wechat-gzh-cli user blacklist-del --openids o1,o2"},
		}))
		cmds = append(cmds, cmd)
	}

	// change-openid
	{
		var jsonOut bool
		cmd := &cobra.Command{
			Use:   "change-openid --openids o1,o2",
			Short: "转换 openid(迁移)",
			Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
			RunE: func(c *cobra.Command, _ []string) error {
				if agentHelpShortcut(c) {
					return nil
				}
				ids, _ := c.Flags().GetString("openids")
				openids, err := parseOpenIDs(ids)
				if err != nil {
					return err
				}
				return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
					list, err := cl.ChangeOpenID(ctx, openids)
					if err != nil {
						return err
					}
					if jsonOut {
						return writeJSON(map[string]any{"result_list": list})
					}
					for _, r := range list {
						if r.ErrCode != 0 {
							fmt.Fprintf(c.OutOrStdout(), "%s -> (err %d: %s)\n", r.OpenID, r.ErrCode, r.ErrMsg)
						} else {
							fmt.Fprintf(c.OutOrStdout(), "%s -> %s\n", r.OpenID, r.NewOpenID)
						}
					}
					return nil
				})
			},
		}
		cmd.Flags().String("openids", "", "逗号分隔的 openid")
		cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
		addAgentHelpFlag(cmd)
		registerAgentHelp("wechat-gzh-cli user change-openid", RenderAgentHelp(AgentHelpDoc{
			Command:        "wechat-gzh-cli user change-openid --openids o_old",
			Summary:        "账号迁移后把旧 openid 转换为新 openid。",
			PreConditions:  []string{"已配置鉴权。", "--openids 必填(旧账号 openid)。"},
			PostConditions: []string{"stdout 输出新旧 openid 映射。"},
			Idempotency:    "YES.",
			Atomicity:      "Atomic.",
			SideEffects:    []string{"无。"},
			CommonErrors:   []ErrSpec{},
			Scope:          "开发者权限。",
			RateLimit:      "常规。",
			RaceConditions: []string{"仅迁移后可转换。"},
			Orchestration:  "迁移后用新 openid 继续业务。",
			ParsingHints:   "--json 输出 result_list。",
			Examples:       []string{"wechat-gzh-cli user change-openid --openids o_old --json"},
		}))
		cmds = append(cmds, cmd)
	}

	return cmds
}

// parseOpenIDs 解析逗号分隔的 openid 列表。
func parseOpenIDs(s string) ([]string, error) {
	if s == "" {
		return nil, fmt.Errorf("--openids 必填")
	}
	var openids []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			openids = append(openids, p)
		}
	}
	if len(openids) == 0 {
		return nil, fmt.Errorf("无有效 openid")
	}
	return openids, nil
}
