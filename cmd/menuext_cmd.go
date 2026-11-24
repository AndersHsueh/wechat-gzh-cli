package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"tx-gongzhonghao-cli/internal/client"
)

func newMenuSelfCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "self",
		Short: "查询自定义菜单信息",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				m, err := cl.GetSelfMenu(ctx)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(m)
				}
				fmt.Fprintf(c.OutOrStdout(), "menu_open=%d\n", m.IsMenuOpen)
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli menu self", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli menu self",
		Summary:        "查询当前使用的菜单(含官方后台设置),与 menu list 不同,能看到官网配置。",
		PreConditions:  []string{"已配置鉴权。"},
		PostConditions: []string{"stdout 输出菜单开启状态与结构。"},
		Idempotency:    "YES. 只读。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"无。"},
		Orchestration:  "用于检测菜单配置后决定是否需要 menu create。",
		ParsingHints:   "--json 输出 selfmenu_info。",
		Examples:       []string{"wechat-gzh-cli menu self --json"},
	}))
	return cmd
}

func newMenuAddConditionalCmd() *cobra.Command {
	var bodyJSON string
	cmd := &cobra.Command{
		Use:   "add-conditional --body '{...}'",
		Short: "创建个性化菜单",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			if bodyJSON == "" {
				return fmt.Errorf("--body 必填(JSON)")
			}
			var body map[string]any
			if err := json.Unmarshal([]byte(bodyJSON), &body); err != nil {
				return fmt.Errorf("--body 不是合法 JSON: %w", err)
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				res, err := cl.AddConditionalMenu(ctx, body)
				if err != nil {
					return err
				}
				fmt.Fprintln(c.OutOrStdout(), res.MenuID)
				return nil
			})
		},
	}
	cmd.Flags().String("body", "", "菜单 JSON(button + matchrule)")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli menu add-conditional", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli menu add-conditional --body '{\"button\":[...],\"matchrule\":{...}}'",
		Summary:        "创建个性化菜单,按 matchrule 让不同用户看到不同菜单。",
		PreConditions:  []string{"已配置鉴权。", "先创建了默认菜单。", "--body 含 button 与 matchrule,不包含隐私筛选字段。"},
		PostConditions: []string{"stdout 输出 menuid。"},
		Idempotency:    "NO. 每次创建新个性化菜单。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"新增一个个性化菜单;每日限 2000 次。"},
		CommonErrors:   []ErrSpec{{Code: "65320", Desc: "match rule violates privacy", Fix: "matchrule 去除性别/地区/语言等隐私字段。"}},
		Scope:          "开发者权限(消息与菜单)。",
		RateLimit:      "每日 2000 次。",
		RaceConditions: []string{"用户多标签时以最后打的标签匹配。"},
		Orchestration:  "创建后用 menu try-match 测试匹配结果。",
		ParsingHints:   "文本输出 menuid。",
		Examples:       []string{`wechat-gzh-cli menu add-conditional --body '{"button":[{"type":"click","name":"x","key":"k"}],"matchrule":{"tag_id":"2"}}'`},
	}))
	return cmd
}

func newMenuDelConditionalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "del-conditional --menuid ID",
		Short: "删除个性化菜单",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			id, _ := c.Flags().GetString("menuid")
			if id == "" {
				return fmt.Errorf("--menuid 必填")
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				if err := cl.DelConditionalMenu(ctx, id); err != nil {
					return err
				}
				fmt.Fprintln(c.OutOrStdout(), "deleted")
				return nil
			})
		},
	}
	cmd.Flags().String("menuid", "", "个性化菜单 ID")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli menu del-conditional", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli menu del-conditional --menuid ID",
		Summary:        "删除指定个性化菜单。",
		PreConditions:  []string{"已配置鉴权。", "--menuid 必填。"},
		PostConditions: []string{"stdout 输出 deleted。"},
		Idempotency:    "NO.",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"删除一个个性化菜单;每日限 2000 次。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限。",
		RateLimit:      "每日 2000 次。",
		RaceConditions: []string{"无。"},
		Orchestration:  "删除后对应用户回到默认菜单。",
		ParsingHints:   "文本输出 deleted。",
		Examples:       []string{"wechat-gzh-cli menu del-conditional --menuid 208379533"},
	}))
	return cmd
}

func newMenuTryMatchCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "try-match --userid xxx",
		Short: "测试个性化菜单匹配结果",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			id, _ := c.Flags().GetString("userid")
			if id == "" {
				return fmt.Errorf("--userid 必填")
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				btn, err := cl.TryMatchMenu(ctx, id)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(btn)
				}
				for _, b := range btn {
					if n, ok := b["name"]; ok {
						fmt.Fprintln(c.OutOrStdout(), n)
					}
				}
				return nil
			})
		},
	}
	cmd.Flags().String("userid", "", "用户 OpenID 或微信号")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli menu try-match", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli menu try-match --userid xxx",
		Summary:        "测试指定用户实际看到的菜单(个性化匹配结果)。",
		PreConditions:  []string{"已配置鉴权。", "--userid 必填。"},
		PostConditions: []string{"stdout 输出该用户看到的菜单按钮。"},
		Idempotency:    "YES. 只读。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限。",
		RateLimit:      "每日 20000 次。",
		RaceConditions: []string{"包含已废除字段的菜单不参与匹配。"},
		Orchestration:  "创建个性化菜单后验证匹配是否生效。",
		ParsingHints:   "--json 输出 button 数组。",
		Examples:       []string{"wechat-gzh-cli menu try-match --userid o_xxx --json"},
	}))
	return cmd
}
