package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"tx-gongzhonghao-cli/internal/client"
)

func newCustomerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "customer-msg",
		Short: "客服消息与客服管理",
		Long:  "客服消息发送、客服账号管理、会话控制与聊天记录。",
	}
	cmd.AddCommand(newCustomerSendCmd())
	for _, sub := range newCustomerAccountCmds() {
		cmd.AddCommand(sub)
	}
	for _, sub := range newCustomerSessionCmds() {
		cmd.AddCommand(sub)
	}
	for _, sub := range newCustomerMsgExtCmds() {
		cmd.AddCommand(sub)
	}
	return cmd
}

func newCustomerSendCmd() *cobra.Command {
	var jsonOut bool
	var payloadJSON string
	cmd := &cobra.Command{
		Use:   "send --openid xxx --type text --payload '{\"content\":\"hi\"}'",
		Short: "发送客服消息",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			openid, _ := c.Flags().GetString("openid")
			msgType, _ := c.Flags().GetString("type")
			if openid == "" || msgType == "" {
				return fmt.Errorf("--openid 与 --type 必填")
			}
			var payload map[string]any
			if payloadJSON != "" {
				if err := json.Unmarshal([]byte(payloadJSON), &payload); err != nil {
					return fmt.Errorf("--payload 不是合法 JSON: %w", err)
				}
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				if err := cl.SendCustomerMessage(ctx, openid, msgType, payload); err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(map[string]bool{"ok": true})
				}
				fmt.Fprintln(c.OutOrStdout(), "ok")
				return nil
			})
		},
	}
	cmd.Flags().String("openid", "", "接收消息用户的 openid")
	cmd.Flags().String("type", "", "消息类型: text/image/voice/video/music/news/mpnews/mpnewsarticle/msgmenu/wxcard/miniprogrampage")
	cmd.Flags().String("payload", "", "消息内容 JSON(如 text 的 {\"content\":\"hi\"})")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli customer-msg send", RenderAgentHelp(AgentHelpDoc{
		Command: "wechat-gzh-cli customer-msg send --openid o_xxx --type text --payload '{\"content\":\"hi\"}'",
		Summary: "向指定用户发送客服消息(需用户近期交互过)。",
		PreConditions: []string{
			"已配置鉴权。",
			"--openid 为该公众号用户。",
			"用户与公众号 48 小时内有过交互(否则需用 get-session... 或客服会话)。",
		},
		PostConditions: []string{"stdout 输出 ok;消息推送到用户会话。"},
		Idempotency:    "NO. 每次调用发送新消息。",
		Atomicity:      "Atomic. 单次 POST。",
		SideEffects:    []string{"用户收到一条客服消息。"},
		CommonErrors: []ErrSpec{
			{Code: "45015", Desc: "response out of time limit", Fix: "用户 48 小时窗口外,无法发送。"},
			{Code: "45047", Desc: "out of response count", Fix: "客服接口调用频繁,限速。"},
			{Code: "40003", Desc: "invalid openid", Fix: "openid 错误。"},
		},
		Scope:          "需客服/消息权限(服务号或认证公众号)。",
		RateLimit:      "受客服消息发送频率限制(每用户/账号窗口,见 45047)。",
		RaceConditions: []string{"无独立竞态;注意 48 小时窗口。"},
		Orchestration:  "发送前可用客服输入状态命令(未内置)提示“正在输入”。",
		ParsingHints:   "文本输出 ok 表示成功。",
		Examples: []string{
			`wechat-gzh-cli customer-msg send --openid o_xxx --type text --payload '{"content":"hello"}' --json`,
		},
	}))
	return cmd
}
