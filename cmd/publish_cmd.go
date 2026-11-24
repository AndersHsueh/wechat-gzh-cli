package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"tx-gongzhonghao-cli/internal/client"
)

func newPublishCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publish",
		Short: "草稿发布能力",
		Long:  "发布草稿、状态查询、已发布列表与删除。",
	}
	cmd.AddCommand(newPublishSubmitCmd())
	cmd.AddCommand(newPublishGetCmd())
	cmd.AddCommand(newPublishListCmd())
	cmd.AddCommand(newPublishDeleteCmd())
	cmd.AddCommand(newPublishGetArticleCmd())
	return cmd
}

func newPublishSubmitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit --media-id MEDIA_ID",
		Short: "发布草稿",
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
				pid, err := cl.PublishSubmit(ctx, id)
				if err != nil {
					return err
				}
				fmt.Fprintln(c.OutOrStdout(), pid)
				return nil
			})
		},
	}
	cmd.Flags().String("media-id", "", "要发布的草稿 media_id")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli publish submit", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli publish submit --media-id MEDIA_ID",
		Summary:        "提交草稿发布,返回 publish_id(异步真正发布)。",
		PreConditions:  []string{"已配置鉴权。", "--media-id 为草稿箱中的草稿。"},
		PostConditions: []string{"stdout 输出 publish_id;发布任务已提交。"},
		Idempotency:    "NO. 每次提交新发布任务。",
		Atomicity:      "Not atomic. 提交成功即入队,真正发布异步,发布完成回调 PUBLISHJOBFINISH。",
		SideEffects:    []string{"向订阅用户发布内容。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限(发布)。",
		RateLimit:      "受发表次数限制。",
		RaceConditions: []string{"发布状态用 publish get 查询,勿轻信提交成功。"},
		Orchestration:  "提交后用 publish get 轮询 publish_status 至成功。",
		ParsingHints:   "文本输出 publish_id。",
		Examples:       []string{"wechat-gzh-cli publish submit --media-id m_xxx"},
	}))
	return cmd
}

func newPublishGetCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "get --publish-id ID",
		Short: "发布状态查询",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			id, _ := c.Flags().GetString("publish-id")
			if id == "" {
				return fmt.Errorf("--publish-id 必填")
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				st, err := cl.PublishGet(ctx, id)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(st)
				}
				fmt.Fprintf(c.OutOrStdout(), "publish_status=%d\n", st.PublishState)
				if st.ArticleDetail != nil {
					fmt.Fprintf(c.OutOrStdout(), "articles=%d\n", st.ArticleDetail.Count)
				}
				return nil
			})
		},
	}
	cmd.Flags().String("publish-id", "", "发布任务 id")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli publish get", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli publish get --publish-id ID",
		Summary:        "查询发布任务状态(PUBLISH_STATUS)。",
		PreConditions:  []string{"已配置鉴权。", "--publish-id 必填。"},
		PostConditions: []string{"stdout 输出 publish_status。"},
		Idempotency:    "YES. 只读。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限(发布)。",
		RateLimit:      "常规。",
		RaceConditions: []string{"状态会变化,需轮询。"},
		Orchestration:  "发布后轮询直到状态为 0(成功)或失败。",
		ParsingHints:   "--json 输出 publish_status/article_detail。",
		Examples:       []string{"wechat-gzh-cli publish get --publish-id 100000001 --json"},
	}))
	return cmd
}

func newPublishListCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "list [--offset 0 --count 20 --no-content]",
		Short: "获取已发布消息列表",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				offset, _ := c.Flags().GetInt("offset")
				count, _ := c.Flags().GetInt("count")
				nc, _ := c.Flags().GetBool("no-content")
				res, err := cl.PublishBatchGet(ctx, offset, count, nc)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(res)
				}
				for _, it := range res.Item {
					fmt.Fprintf(c.OutOrStdout(), "%s\t%d\n", it.ArticleID, it.UpdateTime)
				}
				return nil
			})
		},
	}
	cmd.Flags().Int("offset", 0, "偏移")
	cmd.Flags().Int("count", 20, "数量(1-20)")
	cmd.Flags().Bool("no-content", false, "不返回 content 字段")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli publish list", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli publish list",
		Summary:        "分页获取已成功发布的消息列表。",
		PreConditions:  []string{"已配置鉴权。"},
		PostConditions: []string{"stdout 输出已发布 article_id 列表。"},
		Idempotency:    "YES. 只读。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限(发布)。",
		RateLimit:      "常规。",
		RaceConditions: []string{"分页快照。"},
		Orchestration:  "管理已发布内容。",
		ParsingHints:   "--json 输出 item。",
		Examples:       []string{"wechat-gzh-cli publish list --json"},
	}))
	return cmd
}

func newPublishDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete --article-id ID",
		Short: "删除发布文章",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			id, _ := c.Flags().GetString("article-id")
			if id == "" {
				return fmt.Errorf("--article-id 必填")
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				if err := cl.PublishDelete(ctx, id); err != nil {
					return err
				}
				fmt.Fprintln(c.OutOrStdout(), "deleted")
				return nil
			})
		},
	}
	cmd.Flags().String("article-id", "", "已发布文章的 article_id")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli publish delete", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli publish delete --article-id ID",
		Summary:        "删除已发布文章(不可逆)。",
		PreConditions:  []string{"已配置鉴权。", "--article-id 必填。"},
		PostConditions: []string{"stdout 输出 deleted;文章下线。"},
		Idempotency:    "NO. 不可逆。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"文章被永久删除。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限(发布)。",
		RateLimit:      "常规。",
		RaceConditions: []string{"无。"},
		Orchestration:  "删除前确认。",
		ParsingHints:   "文本输出 deleted。",
		Examples:       []string{"wechat-gzh-cli publish delete --article-id a_xxx"},
	}))
	return cmd
}

func newPublishGetArticleCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "get-article --article-id ID",
		Short: "获取已发布图文信息",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			id, _ := c.Flags().GetString("article-id")
			if id == "" {
				return fmt.Errorf("--article-id 必填")
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				items, err := cl.GetPublishedArticle(ctx, id)
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
	cmd.Flags().String("article-id", "", "已发布文章 article_id")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli publish get-article", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli publish get-article --article-id ID",
		Summary:        "获取已发布图文的详细信息。",
		PreConditions:  []string{"已配置鉴权。", "--article-id 必填。"},
		PostConditions: []string{"stdout 输出标题列表。"},
		Idempotency:    "YES. 只读。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限(发布)。",
		RateLimit:      "常规。",
		RaceConditions: []string{"无。"},
		Orchestration:  "查看文章详情。",
		ParsingHints:   "--json 输出 news_item。",
		Examples:       []string{"wechat-gzh-cli publish get-article --article-id a_xxx --json"},
	}))
	return cmd
}
