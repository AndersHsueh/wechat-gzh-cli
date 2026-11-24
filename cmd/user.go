package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"tx-gongzhonghao-cli/internal/client"
)

func newUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "用户管理",
		Long:  "用户信息查询与关注者列表。",
	}
	cmd.AddCommand(newUserInfoCmd())
	cmd.AddCommand(newUserFollowersCmd())
	return cmd
}

func newUserInfoCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "info --openid OPENID",
		Short: "获取用户基本信息",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			openid, _ := c.Flags().GetString("openid")
			lang, _ := c.Flags().GetString("lang")
			if openid == "" {
				return fmt.Errorf("--openid 必填")
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				info, err := cl.GetUserInfo(ctx, openid, lang)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(info)
				}
				fmt.Fprintf(c.OutOrStdout(), "%s\t%s\t%d\n", info.OpenID, info.Nickname, info.Subscribe)
				return nil
			})
		},
	}
	cmd.Flags().String("openid", "", "用户 openid")
	cmd.Flags().String("lang", "zh_CN", "语言")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli user info", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli user info --openid OPENID",
		Summary:        "根据 OpenID 获取用户基本信息(关注状态/昵称/地区/标签等)。",
		PreConditions:  []string{"已配置鉴权。", "--openid 必填且为该公众号下的用户。"},
		PostConditions: []string{"stdout 输出 openid/nickname/subscribe,或 --json 输出完整字段。"},
		Idempotency:    "YES. 只读。",
		Atomicity:      "Atomic. 单次 GET。",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{{Code: "40003", Desc: "invalid openid", Fix: "确认 openid 属于本公众号。"}},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"subscribe=0 表示已取消关注,此时只有 openid 字段。"},
		Orchestration:  "结合标签、消息做消息推送决策。",
		ParsingHints:   "文本输出 openid\\tnickname\\tsubscribe;--json 输出完整 UserInfo。",
		Examples:       []string{"wechat-gzh-cli user info --openid o_xxx --json"},
	}))
	return cmd
}

func newUserFollowersCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "followers [--next-openid OPENID]",
		Short: "获取关注用户列表",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			nextOpenID, _ := c.Flags().GetString("next-openid")
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				res, err := cl.GetFollowers(ctx, nextOpenID)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(res)
				}
				for _, id := range res.Data.OpenID {
					fmt.Fprintln(c.OutOrStdout(), id)
				}
				return nil
			})
		},
	}
	cmd.Flags().String("next-openid", "", "分页游标;首页为空")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli user followers", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli user followers --next-openid xxx",
		Summary:        "分页获取全部关注用户的 OpenID 列表。",
		PreConditions:  []string{"已配置鉴权。", "--next-openid 为空表示第一页。"},
		PostConditions: []string{"stdout 输出 openid 列表(每行一个),JSON 含 next_openid 游标。"},
		Idempotency:    "YES. 只读。",
		Atomicity:      "Atomic. 单次 GET。",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"翻页是点状快照,关注数可能在分页期间变化。"},
		Orchestration:  "配合 --next-openid 循环拉全量,再对每位用户做信息查询。",
		ParsingHints:   "--json 输出 {\"total\":..,\"data\":{\"openid\":[..]},\"next_openid\":..}。",
		Examples:       []string{"wechat-gzh-cli user followers --json"},
	}))
	return cmd
}
