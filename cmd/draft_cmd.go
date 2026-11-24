package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"tx-gongzhonghao-cli/internal/client"
)

func newDraftCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "draft",
		Short: "草稿管理",
		Long:  "图文草稿的增删改查。",
	}
	cmd.AddCommand(newDraftAddCmd())
	cmd.AddCommand(newDraftListCmd())
	cmd.AddCommand(newDraftGetCmd())
	cmd.AddCommand(newDraftCountCmd())
	cmd.AddCommand(newDraftUpdateCmd())
	cmd.AddCommand(newDraftDeleteCmd())
	return cmd
}

func newDraftAddCmd() *cobra.Command {
	var bodyJSON string
	cmd := &cobra.Command{
		Use:   "add --body '[{...}]'",
		Short: "新增草稿",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			if bodyJSON == "" {
				return fmt.Errorf("--body 必填(JSON 数组)")
			}
			var articles []map[string]any
			if err := json.Unmarshal([]byte(bodyJSON), &articles); err != nil {
				return fmt.Errorf("--body 不是合法 JSON 数组: %w", err)
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				mediaID, err := cl.DraftAdd(ctx, articles)
				if err != nil {
					return err
				}
				fmt.Fprintln(c.OutOrStdout(), mediaID)
				return nil
			})
		},
	}
	cmd.Flags().String("body", "", "图文数组 JSON(含 title/content/thumb_media_id 等)")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli draft add", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli draft add --body '[{...}]'",
		Summary:        "新增图文草稿到草稿箱,返回草稿 media_id。",
		PreConditions:  []string{"已配置鉴权。", "--body 为图文数组 JSON。"},
		PostConditions: []string{"stdout 输出草稿 media_id;草稿进入草稿箱。"},
		Idempotency:    "NO. 每次新增新草稿。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"草稿箱新增一条草稿。"},
		CommonErrors:   []ErrSpec{{Code: "45166", Desc: "invalid content", Fix: "content 超 2 万字符或含外部图片链接。"}},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"无。"},
		Orchestration:  "新增后可 publish submit 发布。",
		ParsingHints:   "文本输出 media_id。",
		Examples:       []string{`wechat-gzh-cli draft add --body '[{"title":"t","content":"c","thumb_media_id":"m"}]'`},
	}))
	return cmd
}

func newDraftListCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "list [--offset 0 --count 20]",
		Short: "获取草稿列表",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				offset, _ := c.Flags().GetInt("offset")
				count, _ := c.Flags().GetInt("count")
				res, err := cl.DraftBatchGet(ctx, offset, count)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(res)
				}
				for _, it := range res.Item {
					fmt.Fprintf(c.OutOrStdout(), "%s\t%d\n", it.MediaID, it.UpdateTime)
				}
				return nil
			})
		},
	}
	cmd.Flags().Int("offset", 0, "偏移")
	cmd.Flags().Int("count", 20, "数量(1-20)")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli draft list", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli draft list",
		Summary:        "分页获取草稿列表。",
		PreConditions:  []string{"已配置鉴权。"},
		PostConditions: []string{"stdout 输出草稿 media_id 列表。"},
		Idempotency:    "YES. 只读。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"分页快照。"},
		Orchestration:  "根据列表选择要发布的草稿。",
		ParsingHints:   "--json 输出 item 数组。",
		Examples:       []string{"wechat-gzh-cli draft list --json"},
	}))
	return cmd
}

func newDraftGetCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "get --media-id MEDIA_ID",
		Short: "获取草稿详情",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			id, _ := c.Flags().GetString("media-id")
			if id == "" {
				return fmt.Errorf("--media-id 必填")
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				items, err := cl.DraftGet(ctx, id)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(items)
				}
				for _, it := range items {
					if t, ok := it["title"]; ok {
						fmt.Fprintln(c.OutOrStdout(), t)
					}
				}
				return nil
			})
		},
	}
	cmd.Flags().String("media-id", "", "草稿 media_id")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli draft get", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli draft get --media-id MEDIA_ID",
		Summary:        "获取草稿详细内容。",
		PreConditions:  []string{"已配置鉴权。", "--media-id 必填。"},
		PostConditions: []string{"stdout 输出草稿标题/内容。"},
		Idempotency:    "YES. 只读。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"无。"},
		Orchestration:  "发布/编辑前查看草稿。",
		ParsingHints:   "--json 输出 news_item。",
		Examples:       []string{"wechat-gzh-cli draft get --media-id m_xxx --json"},
	}))
	return cmd
}

func newDraftCountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "count",
		Short: "获取草稿总数",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				n, err := cl.DraftCount(ctx)
				if err != nil {
					return err
				}
				fmt.Fprintf(c.OutOrStdout(), "%d\n", n)
				return nil
			})
		},
	}
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli draft count", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli draft count",
		Summary:        "获取草稿总数(仅数量)。",
		PreConditions:  []string{"已配置鉴权。"},
		PostConditions: []string{"stdout 输出数量。"},
		Idempotency:    "YES. 只读。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"无。"},
		Orchestration:  "监控草稿量。",
		ParsingHints:   "文本输出数字。",
		Examples:       []string{"wechat-gzh-cli draft count"},
	}))
	return cmd
}

func newDraftUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update --media-id M --index 0 --article '{...}'",
		Short: "更新草稿",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			id, _ := c.Flags().GetString("media-id")
			idx, _ := c.Flags().GetInt("index")
			article, _ := c.Flags().GetString("article")
			if id == "" || article == "" {
				return fmt.Errorf("--media-id 与 --article 必填")
			}
			var art map[string]any
			if err := json.Unmarshal([]byte(article), &art); err != nil {
				return fmt.Errorf("--article 不是合法 JSON: %w", err)
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				if err := cl.DraftUpdate(ctx, id, idx, art); err != nil {
					return err
				}
				fmt.Fprintln(c.OutOrStdout(), "ok")
				return nil
			})
		},
	}
	cmd.Flags().String("media-id", "", "草稿 media_id")
	cmd.Flags().Int("index", 0, "要更新的文章在图文中的位置(第一项为 0)")
	cmd.Flags().String("article", "", "图文 JSON")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli draft update", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli draft update --media-id M --index 0 --article '{...}'",
		Summary:        "修改草稿中的某篇文章。",
		PreConditions:  []string{"已配置鉴权。", "--media-id/--article 必填。"},
		PostConditions: []string{"stdout 输出 ok;草稿被更新。"},
		Idempotency:    "NO.",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"草稿内容被修改。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"无。"},
		Orchestration:  "编辑后再发布。",
		ParsingHints:   "文本输出 ok。",
		Examples:       []string{`wechat-gzh-cli draft update --media-id m --index 0 --article '{"title":"t"}'`},
	}))
	return cmd
}

func newDraftDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete --media-id MEDIA_ID",
		Short: "删除草稿",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			id, _ := c.Flags().GetString("media-id")
			if id == "" {
				return fmt.Errorf("--media-id 必填")
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				if err := cl.DraftDelete(ctx, id); err != nil {
					return err
				}
				fmt.Fprintln(c.OutOrStdout(), "deleted")
				return nil
			})
		},
	}
	cmd.Flags().String("media-id", "", "草稿 media_id")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli draft delete", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli draft delete --media-id MEDIA_ID",
		Summary:        "删除指定草稿(不可撤销)。",
		PreConditions:  []string{"已配置鉴权。", "--media-id 必填。"},
		PostConditions: []string{"stdout 输出 deleted;草稿被删除。"},
		Idempotency:    "NO. 有副作用。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"草稿被永久删除。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"无。"},
		Orchestration:  "删除前确认。",
		ParsingHints:   "文本输出 deleted。",
		Examples:       []string{"wechat-gzh-cli draft delete --media-id m_xxx"},
	}))
	return cmd
}
