package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"tx-gongzhonghao-cli/internal/client"
)

func newTagCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag",
		Short: "用户标签管理",
		Long:  "标签列表与创建。",
	}
	cmd.AddCommand(newTagListCmd())
	cmd.AddCommand(newTagCreateCmd())
	return cmd
}

func newTagListCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "获取标签列表",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				tags, err := cl.ListTags(ctx)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(tags)
				}
				for _, t := range tags {
					fmt.Fprintf(c.OutOrStdout(), "%d\t%s\t%d\n", t.ID, t.Name, t.Count)
				}
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli tag list", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli tag list",
		Summary:        "获取公众号全部标签及其用户数。",
		PreConditions:  []string{"已配置鉴权。"},
		PostConditions: []string{"stdout 输出 id\\tname\\tcount,或 --json 输出完整列表。"},
		Idempotency:    "YES.",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"无。"},
		Orchestration:  "标签 id 用于群发 filter.tag_id 或个性化菜单。",
		ParsingHints:   "--json 输出 [{\"id\":..,\"name\":..,\"count\":..}]。",
		Examples:       []string{"wechat-gzh-cli tag list --json"},
	}))
	return cmd
}

func newTagCreateCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "create --name 标签名",
		Short: "创建标签",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			name, _ := c.Flags().GetString("name")
			if name == "" {
				return fmt.Errorf("--name 必填")
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				tag, err := cl.CreateTag(ctx, name)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(tag)
				}
				fmt.Fprintf(c.OutOrStdout(), "%d\n", tag.ID)
				return nil
			})
		},
	}
	cmd.Flags().String("name", "", "标签名称")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli tag create", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli tag create --name 标签",
		Summary:        "创建用户标签,返回标签 id。",
		PreConditions:  []string{"已配置鉴权。", "--name 必填(标签名)。"},
		PostConditions: []string{"stdout 输出新标签 id,或 --json 输出 {\"id\":..}。"},
		Idempotency:    "NO. 每次创建都会新建标签;同名标签会重复创建。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"公众号新增一个标签。"},
		CommonErrors:   []ErrSpec{{Code: "45157", Desc: "tag name too long", Fix: "标签名限 30 字符。"}, {Code: "45158", Desc: "tag name invalid", Fix: "含非法字符。"}, {Code: "45058", Desc: "too many tags", Fix: "已超 100 个标签上限。"}},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"并发创建同名标签会产生多个。"},
		Orchestration:  "创建后可在 tag list 中查看 id,用于群发/菜单定向。",
		ParsingHints:   "文本输出 tag id;--json 输出 {\"id\":..,\"name\":..}。",
		Examples:       []string{"wechat-gzh-cli tag create --name vip --json"},
	}))
	return cmd
}
