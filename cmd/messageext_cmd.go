package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"tx-gongzhonghao-cli/internal/client"
)

func newMessageMassDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mass-delete [--msg-id ID] [--url URL] [--idx N]",
		Short: "删除群发消息",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			msgID, _ := c.Flags().GetInt64("msg-id")
			url, _ := c.Flags().GetString("url")
			idx, _ := c.Flags().GetInt("idx")
			if msgID == 0 && url == "" {
				return fmt.Errorf("--msg-id 或 --url 至少一个")
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				if err := cl.MassDelete(ctx, msgID, idx, url); err != nil {
					return err
				}
				fmt.Fprintln(c.OutOrStdout(), "deleted")
				return nil
			})
		},
	}
	cmd.Flags().Int64("msg-id", 0, "群发消息 ID")
	cmd.Flags().String("url", "", "要删除内容的 url")
	cmd.Flags().Int("idx", 0, "第一篇为1;0 删除全部(随 msg-id)")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli message mass-delete", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli message mass-delete --msg-id 12345",
		Summary:        "删除通过 api 群发且已成功的文章/视频消息。",
		PreConditions:  []string{"已配置鉴权。", "--msg-id 或 --url 至少一个。"},
		PostConditions: []string{"stdout 输出 deleted;消息正文失效(已收到用户本地仍可见)。"},
		Idempotency:    "NO. 删除有副作用。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"文章/视频内容删除;同一内容的多次群发可能一并失效。"},
		CommonErrors:   []ErrSpec{{Code: "40059", Desc: "invalid msg id", Fix: "检查 msg-id。"}},
		Scope:          "开发者权限(群发)。",
		RateLimit:      "常规。",
		RaceConditions: []string{"只能删 api 发送且成功的文章/视频。"},
		Orchestration:  "删除后结合 mass-status 确认。",
		ParsingHints:   "文本输出 deleted。",
		Examples:       []string{"wechat-gzh-cli message mass-delete --msg-id 30124"},
	}))
	return cmd
}

func newMessageSpeedGetCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "speed-get",
		Short: "获取群发速度",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				s, err := cl.MassSpeedGet(ctx)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(s)
				}
				fmt.Fprintf(c.OutOrStdout(), "speed=%d realspeed=%d\n", s.Speed, s.RealSpeed)
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli message speed-get", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli message speed-get",
		Summary:        "获取当前群发速度级别与真实值。",
		PreConditions:  []string{"已配置鉴权。"},
		PostConditions: []string{"stdout 输出 speed/realspeed。"},
		Idempotency:    "YES. 只读。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限(群发)。",
		RateLimit:      "常规。",
		RaceConditions: []string{"无。"},
		Orchestration:  "群发前确认速度,过大时用 speed-set 调低。",
		ParsingHints:   "--json 输出 {\"speed\":..,\"realspeed\":..}。",
		Examples:       []string{"wechat-gzh-cli message speed-get --json"},
	}))
	return cmd
}

func newMessageSpeedSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "speed-set --speed N",
		Short: "设置群发速度",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			speed, _ := c.Flags().GetInt("speed")
			if speed < 0 || speed > 4 {
				return fmt.Errorf("--speed 取值 0-4")
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				if err := cl.MassSpeedSet(ctx, speed); err != nil {
					return err
				}
				fmt.Fprintln(c.OutOrStdout(), "ok")
				return nil
			})
		},
	}
	cmd.Flags().Int("speed", 0, "级别 0-4(0=80w/分钟..4=10w/分钟)")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli message speed-set", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli message speed-set --speed 3",
		Summary:        "设置群发速度级别(0-4)。",
		PreConditions:  []string{"已配置鉴权。", "--speed 0-4。"},
		PostConditions: []string{"stdout 输出 ok;群发速度生效。"},
		Idempotency:    "NO. 会改变群发速度设置。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"影响后续群发速度。"},
		CommonErrors:   []ErrSpec{{Code: "45083", Desc: "speed 越界", Fix: "检查 --speed 0-4。"}},
		Scope:          "开发者权限(群发)。",
		RateLimit:      "常规。",
		RaceConditions: []string{"无。"},
		Orchestration:  "群发前设置,完成后用 speed-get 确认。",
		ParsingHints:   "文本输出 ok。",
		Examples:       []string{"wechat-gzh-cli message speed-set --speed 3"},
	}))
	return cmd
}

func newMessagePreviewCmd() *cobra.Command {
	var bodyJSON string
	cmd := &cobra.Command{
		Use:   "preview --body '{...}'",
		Short: "预览消息",
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
				if err := cl.MassPreview(ctx, body); err != nil {
					return err
				}
				fmt.Fprintln(c.OutOrStdout(), "ok")
				return nil
			})
		},
	}
	cmd.Flags().String("body", "", "预览请求 JSON")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli message preview", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli message preview --body '{...}'",
		Summary:        "向指定用户发送预览,在手机端检查样式排版。",
		PreConditions:  []string{"已配置鉴权。", "--body 含接受者(touser/towxname)与消息内容。"},
		PostConditions: []string{"stdout 输出 ok;指定用户收到预览消息。"},
		Idempotency:    "NO.",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"向指定用户发送一条预览;按微信号预览每日限 100 次。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限(群发)。",
		RateLimit:      "按微信号预览每日 100 次。",
		RaceConditions: []string{"无。"},
		Orchestration:  "群发前先 preview 校验样式。",
		ParsingHints:   "文本输出 ok。",
		Examples:       []string{`wechat-gzh-cli message preview --body '{"towxname":"微信号","msgtype":"text","text":{"content":"hi"}}'`},
	}))
	return cmd
}

func newSubscribeSendCmd() *cobra.Command {
	var bodyJSON string
	cmd := &cobra.Command{
		Use:   "send --body '{...}'",
		Short: "发送一次性订阅消息",
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
				if err := cl.SubscribeSend(ctx, body); err != nil {
					return err
				}
				fmt.Fprintln(c.OutOrStdout(), "ok")
				return nil
			})
		},
	}
	cmd.Flags().String("body", "", "订阅消息请求 JSON")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli subscribe send", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli subscribe send --body '{...}'",
		Summary:        "向授权用户发送一次性订阅模板消息。",
		PreConditions:  []string{"已配置鉴权。", "--body 含 touser/template_id/scene/title/data。"},
		PostConditions: []string{"stdout 输出 ok;用户收到订阅通知。"},
		Idempotency:    "NO.",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"用户收到一条订阅消息。"},
		CommonErrors:   []ErrSpec{{Code: "40037", Desc: "invalid template_id", Fix: "检查模板 ID。"}, {Code: "40062", Desc: "title 超长", Fix: "标题 15 字内。"}},
		Scope:          "需要一次性订阅消息权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"已关注用户走会话,未关注走服务通知。"},
		Orchestration:  "推送后无特殊后续。",
		ParsingHints:   "文本输出 ok。",
		Examples:       []string{`wechat-gzh-cli subscribe send --body '{"touser":"o_1","template_id":"T","scene":"1000","title":"通知","data":{"content":{"value":"hi"}}}'`},
	}))
	return cmd
}

func newAutoreplyGetCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "get",
		Short: "获取自动回复规则",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				info, err := cl.GetAutoreply(ctx)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(info)
				}
				fmt.Fprintf(c.OutOrStdout(), "add_friend_open=%d msg_open=%d keywords=%d\n",
					info.IsAddFriendReplyOpen, info.IsAutoreplyOpen, len(info.KeywordAutoreplyInfo.List))
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli autoreply get", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli autoreply get",
		Summary:        "获取公众号当前自动回复规则(关注/消息/关键词)。",
		PreConditions:  []string{"已配置鉴权。", "仅在官方后台配置的自动回复可被查询。"},
		PostConditions: []string{"stdout 输出开关状态与规则条数。"},
		Idempotency:    "YES. 只读。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"自研回复不可查询。"},
		Orchestration:  "查看配置后决定是否调整。",
		ParsingHints:   "--json 输出完整 autoreply 结构。",
		Examples:       []string{"wechat-gzh-cli autoreply get --json"},
	}))
	return cmd
}
