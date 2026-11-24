package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"tx-gongzhonghao-cli/internal/client"
)

func newMenuCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "menu",
		Short: "自定义菜单",
		Long:  "自定义菜单的创建/查询/删除。",
	}
	cmd.AddCommand(newMenuCreateCmd())
	cmd.AddCommand(newMenuListCmd())
	cmd.AddCommand(newMenuDeleteCmd())
	return cmd
}

func newMenuCreateCmd() *cobra.Command {
	var bodyJSON string
	cmd := &cobra.Command{
		Use:   "create --body '{...}'",
		Short: "创建自定义菜单",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			if bodyJSON == "" {
				return fmt.Errorf("--body 必填(JSON)")
			}
			var menu map[string]any
			if err := json.Unmarshal([]byte(bodyJSON), &menu); err != nil {
				return fmt.Errorf("--body 不是合法 JSON: %w", err)
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				if err := cl.CreateMenu(ctx, menu); err != nil {
					return err
				}
				fmt.Fprintln(c.OutOrStdout(), "ok")
				return nil
			})
		},
	}
	cmd.Flags().String("body", "", "菜单 JSON(含 button 数组)")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli menu create", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli menu create --body '{\"button\":[...]}'",
		Summary:        "创建/替换默认自定义菜单(最多 3 个一级菜单,各含最多 5 个二级)。",
		PreConditions:  []string{"已配置鉴权。", "--body 为符合公众号菜单接口的 JSON(button 数组)。"},
		PostConditions: []string{"stdout 输出 ok;公众号默认菜单被替换。"},
		Idempotency:    "NO. 每次调用会整体替换当前默认菜单。",
		Atomicity:      "Atomic. 单次 POST;成功即替换,失败保持原样。",
		SideEffects:    []string{"公众号对话界面的默认菜单被更新。"},
		CommonErrors:   []ErrSpec{{Code: "40016", Desc: "invalid button size", Fix: "一级菜单数量超限(≤3)。"}, {Code: "40018", Desc: "invalid button name size", Fix: "菜单标题超长。"}},
		Scope:          "开发者权限(消息与菜单)。",
		RateLimit:      "常规。",
		RaceConditions: []string{"删除默认菜单会连带删除全部个性化菜单。"},
		Orchestration:  "创建后用 menu list 确认结构,或用 menu delete 回滚。",
		ParsingHints:   "文本输出 ok 表示成功。",
		Examples: []string{
			`wechat-gzh-cli menu create --body '{"button":[{"type":"click","name":"今日","key":"V1001"}]}'`,
		},
	}))
	return cmd
}

func newMenuListCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "获取自定义菜单配置",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				menu, err := cl.GetMenu(ctx)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(menu)
				}
				for _, b := range menu.Menu.Button {
					fmt.Fprintf(c.OutOrStdout(), "%v\n", b)
				}
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli menu list", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli menu list",
		Summary:        "获取默认菜单与全部个性化菜单配置。",
		PreConditions:  []string{"已配置鉴权。"},
		PostConditions: []string{"stdout 输出菜单按钮,或 --json 输出完整 menu 结构。"},
		Idempotency:    "YES. 只读。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"无。"},
		Orchestration:  "查看当前配置以决定是否重建/调整。",
		ParsingHints:   "--json 输出 {\"menu\":{\"button\":..}}。",
		Examples:       []string{"wechat-gzh-cli menu list --json"},
	}))
	return cmd
}

func newMenuDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "删除自定义菜单",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				if err := cl.DeleteMenu(ctx); err != nil {
					return err
				}
				fmt.Fprintln(c.OutOrStdout(), "deleted")
				return nil
			})
		},
	}
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli menu delete", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli menu delete",
		Summary:        "删除当前使用的自定义菜单(含全部个性化菜单)。",
		PreConditions:  []string{"已配置鉴权。"},
		PostConditions: []string{"stdout 输出 deleted;默认菜单与个性化菜单全部删除。"},
		Idempotency:    "NO. 有副作用;删除后菜单被清空。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"公众号菜单被清空。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"删除不可逆,误删需重建。"},
		Orchestration:  "删除前确认,或用 menu create 直接重建。",
		ParsingHints:   "文本输出 deleted 表示成功。",
		Examples:       []string{"wechat-gzh-cli menu delete"},
	}))
	return cmd
}
