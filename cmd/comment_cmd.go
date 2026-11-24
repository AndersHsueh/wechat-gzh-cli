package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"tx-gongzhonghao-cli/internal/client"
)

// newCommentCmd 留言管理(评论)命令组。
func newCommentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "comment",
		Short: "留言评论管理",
		Long:  "对已群发文章的评论进行开关、列表、精选、删除与回复。需公众号具备留言功能权限。",
	}
	cmd.AddCommand(newCommentOpenCmd())
	cmd.AddCommand(newCommentCloseCmd())
	cmd.AddCommand(newCommentListCmd())
	cmd.AddCommand(newCommentMarkelectCmd())
	cmd.AddCommand(newCommentUnmarkelectCmd())
	cmd.AddCommand(newCommentDeleteCmd())
	cmd.AddCommand(newCommentReplyAddCmd())
	cmd.AddCommand(newCommentReplyDelCmd())
	return cmd
}

// articleRefFlags 添加 msg_data_id/index 通用 flag 到 cmd。
func articleRefFlags(c *cobra.Command) {
	c.Flags().Int64("msg-data-id", 0, "群发消息 id(msg_data_id)")
	c.Flags().Int64("index", 0, "多图文时的文章序号(第一页为 0)")
}

// requireArticleRef 校验并返回 msg_data_id/index。
func requireArticleRef(c *cobra.Command) (int64, int64, error) {
	mid, _ := c.Flags().GetInt64("msg-data-id")
	idx, _ := c.Flags().GetInt64("index")
	if mid == 0 {
		return 0, 0, fmt.Errorf("--msg-data-id 必填")
	}
	return mid, idx, nil
}

func newCommentOpenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open --msg-data-id ID [--index 0]",
		Short: "打开文章评论",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			mid, idx, err := requireArticleRef(c)
			if err != nil {
				return err
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				if err := cl.OpenComment(ctx, mid, idx); err != nil {
					return err
				}
				fmt.Fprintln(c.OutOrStdout(), "ok")
				return nil
			})
		},
	}
	articleRefFlags(cmd)
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli comment open", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli comment open --msg-data-id 100 --index 0",
		Summary:        "打开已群发文章的评论功能。",
		PreConditions:  []string{"已配置鉴权。", "公众号具备留言权限。", "--msg-data-id 必填。"},
		PostConditions: []string{"stdout 输出 ok;评论对外可见。"},
		Idempotency:    "NO.",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"文章评论被打开。"},
		CommonErrors:   []ErrSpec{{Code: "45078", Desc: "无留言权限", Fix: "开通留言功能。"}},
		Scope:          "开发者权限(留言功能)。",
		RateLimit:      "常规。",
		RaceConditions: []string{"无。"},
		Orchestration:  "打开后可 comment list 查看评论。",
		ParsingHints:   "文本输出 ok。",
		Examples:       []string{"wechat-gzh-cli comment open --msg-data-id 100"},
	}))
	return cmd
}

func newCommentCloseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close --msg-data-id ID [--index 0]",
		Short: "关闭文章评论",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			mid, idx, err := requireArticleRef(c)
			if err != nil {
				return err
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				if err := cl.CloseComment(ctx, mid, idx); err != nil {
					return err
				}
				fmt.Fprintln(c.OutOrStdout(), "ok")
				return nil
			})
		},
	}
	articleRefFlags(cmd)
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli comment close", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli comment close --msg-data-id 100",
		Summary:        "关闭已群发文章的评论功能。",
		PreConditions:  []string{"已配置鉴权。", "--msg-data-id 必填。"},
		PostConditions: []string{"stdout 输出 ok;评论被关闭。"},
		Idempotency:    "NO.",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"文章评论被关闭。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限(留言功能)。",
		RateLimit:      "常规。",
		RaceConditions: []string{"无。"},
		Orchestration:  "控制评论区热度。",
		ParsingHints:   "文本输出 ok。",
		Examples:       []string{"wechat-gzh-cli comment close --msg-data-id 100"},
	}))
	return cmd
}

func newCommentListCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "list --msg-data-id ID [--index 0 --begin 0 --count 20 --type 0]",
		Short: "查看文章评论列表",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			mid, idx, err := requireArticleRef(c)
			if err != nil {
				return err
			}
			begin, _ := c.Flags().GetInt("begin")
			count, _ := c.Flags().GetInt("count")
			ctype, _ := c.Flags().GetInt("type")
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				res, err := cl.ListComment(ctx, mid, idx, begin, count, ctype)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(res)
				}
				for _, cm := range res.Comment {
					fmt.Fprintf(c.OutOrStdout(), "%d\t%s\t%s\n", cm.CommentID, cm.OpenID, cm.Content)
				}
				return nil
			})
		},
	}
	articleRefFlags(cmd)
	cmd.Flags().Int("begin", 0, "偏移")
	cmd.Flags().Int("count", 20, "数量(1-50)")
	cmd.Flags().Int("type", 0, "0 全部/1 仅精选/2 仅未精选")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli comment list", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli comment list --msg-data-id 100 --type 1",
		Summary:        "查看指定文章评论,可筛选精选/未精选。",
		PreConditions:  []string{"已配置鉴权。", "--msg-data-id 必填。"},
		PostConditions: []string{"stdout 输出 comment_id\\topenid\\tcontent。"},
		Idempotency:    "YES. 只读。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限(留言功能)。",
		RateLimit:      "常规。",
		RaceConditions: []string{"分页快照。"},
		Orchestration:  "用返回的 comment_id 标记精选/回复。",
		ParsingHints:   "--json 输出 comment 数组。",
		Examples:       []string{"wechat-gzh-cli comment list --msg-data-id 100 --count 50 --json"},
	}))
	return cmd
}

func newCommentMarkelectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "markelect --comment-id ID",
		Short: "评论标记精选",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			commentID, err := requireCommentID(c)
			if err != nil {
				return err
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				if err := cl.MarkElect(ctx, commentID); err != nil {
					return err
				}
				fmt.Fprintln(c.OutOrStdout(), "ok")
				return nil
			})
		},
	}
	commentIDFlag(cmd)
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli comment markelect", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli comment markelect --comment-id 123",
		Summary:        "将评论标记为精选(展示在精选位置)。",
		PreConditions:  []string{"已配置鉴权。", "--comment-id 必填。"},
		PostConditions: []string{"stdout 输出 ok;评论被精选。"},
		Idempotency:    "NO.",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"评论进入精选。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限(留言功能)。",
		RateLimit:      "常规。",
		RaceConditions: []string{"无。"},
		Orchestration:  "从 comment list 取 comment_id。",
		ParsingHints:   "文本输出 ok。",
		Examples:       []string{"wechat-gzh-cli comment markelect --comment-id 123"},
	}))
	return cmd
}

func newCommentUnmarkelectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unmarkelect --comment-id ID",
		Short: "评论取消精选",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			commentID, err := requireCommentID(c)
			if err != nil {
				return err
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				if err := cl.UnmarkElect(ctx, commentID); err != nil {
					return err
				}
				fmt.Fprintln(c.OutOrStdout(), "ok")
				return nil
			})
		},
	}
	commentIDFlag(cmd)
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli comment unmarkelect", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli comment unmarkelect --comment-id 123",
		Summary:        "将精选评论取消精选。",
		PreConditions:  []string{"已配置鉴权。", "--comment-id 必填。"},
		PostConditions: []string{"stdout 输出 ok。"},
		Idempotency:    "NO.",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"评论移出精选。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限(留言功能)。",
		RateLimit:      "常规。",
		RaceConditions: []string{"无。"},
		Orchestration:  "调整精选展示。",
		ParsingHints:   "文本输出 ok。",
		Examples:       []string{"wechat-gzh-cli comment unmarkelect --comment-id 123"},
	}))
	return cmd
}

func newCommentDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete --comment-id ID",
		Short: "删除评论",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			commentID, err := requireCommentID(c)
			if err != nil {
				return err
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				if err := cl.DeleteComment(ctx, commentID); err != nil {
					return err
				}
				fmt.Fprintln(c.OutOrStdout(), "deleted")
				return nil
			})
		},
	}
	commentIDFlag(cmd)
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli comment delete", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli comment delete --comment-id 123",
		Summary:        "删除某条评论(不可恢复)。",
		PreConditions:  []string{"已配置鉴权。", "--comment-id 必填。"},
		PostConditions: []string{"stdout 输出 deleted;评论被删除。"},
		Idempotency:    "NO.",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"评论被永久删除。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限(留言功能)。",
		RateLimit:      "常规。",
		RaceConditions: []string{"无。"},
		Orchestration:  "清理违规评论。",
		ParsingHints:   "文本输出 deleted。",
		Examples:       []string{"wechat-gzh-cli comment delete --comment-id 123"},
	}))
	return cmd
}

func newCommentReplyAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reply-add --comment-id ID --content 内容",
		Short: "回复评论",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			commentID, err := requireCommentID(c)
			if err != nil {
				return err
			}
			content, _ := c.Flags().GetString("content")
			if content == "" {
				return fmt.Errorf("--content 必填")
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				if err := cl.ReplyComment(ctx, commentID, content); err != nil {
					return err
				}
				fmt.Fprintln(c.OutOrStdout(), "ok")
				return nil
			})
		},
	}
	commentIDFlag(cmd)
	cmd.Flags().String("content", "", "回复内容")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli comment reply-add", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli comment reply-add --comment-id 123 --content 谢谢",
		Summary:        "回复某条评论。",
		PreConditions:  []string{"已配置鉴权。", "--comment-id/--content 必填。"},
		PostConditions: []string{"stdout 输出 ok;回复已发布。"},
		Idempotency:    "NO.",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"新增一条评论回复。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限(留言功能)。",
		RateLimit:      "常规。",
		RaceConditions: []string{"无。"},
		Orchestration:  "与用户互动运营。",
		ParsingHints:   "文本输出 ok。",
		Examples:       []string{"wechat-gzh-cli comment reply-add --comment-id 123 --content 谢谢"},
	}))
	return cmd
}

func newCommentReplyDelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reply-del --comment-id ID",
		Short: "删除评论回复",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			commentID, err := requireCommentID(c)
			if err != nil {
				return err
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				if err := cl.DeleteReply(ctx, commentID); err != nil {
					return err
				}
				fmt.Fprintln(c.OutOrStdout(), "deleted")
				return nil
			})
		},
	}
	commentIDFlag(cmd)
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli comment reply-del", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli comment reply-del --comment-id 123",
		Summary:        "删除某条评论下的官方回复。",
		PreConditions:  []string{"已配置鉴权。", "--comment-id 必填。"},
		PostConditions: []string{"stdout 输出 deleted;回复被删除。"},
		Idempotency:    "NO.",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"评论回复被删除。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限(留言功能)。",
		RateLimit:      "常规。",
		RaceConditions: []string{"无。"},
		Orchestration:  "删除不当回复。",
		ParsingHints:   "文本输出 deleted。",
		Examples:       []string{"wechat-gzh-cli comment reply-del --comment-id 123"},
	}))
	return cmd
}

// commentIDFlag 添加 comment-id flag。
func commentIDFlag(c *cobra.Command) {
	c.Flags().Int64("comment-id", 0, "评论 id(user_comment_id)")
}

// requireCommentID 校验并返回 comment-id。
func requireCommentID(c *cobra.Command) (int64, error) {
	id, _ := c.Flags().GetInt64("comment-id")
	if id == 0 {
		return 0, fmt.Errorf("--comment-id 必填")
	}
	return id, nil
}
