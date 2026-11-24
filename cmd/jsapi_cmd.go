package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"tx-gongzhonghao-cli/internal/client"
)

// newJSSDKTicketCmd 获取 JS-SDK 临时票据。
func newJSSDKTicketCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "jsapi-ticket",
		Short: "获取 JS-SDK 临时票据",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				tk, err := cl.GetJSSDKTicket(ctx)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(tk)
				}
				fmt.Fprintln(c.OutOrStdout(), tk.Ticket)
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli jsapi-ticket", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli jsapi-ticket",
		Summary:        "获取 JS-SDK 临时票据(ticket,7200s 有效),用于前端 jsapi 签名。",
		PreConditions:  []string{"已配置鉴权。"},
		PostConditions: []string{"stdout 输出 ticket。"},
		Idempotency:    "YES. 只读。",
		Atomicity:      "Atomic. 单次 GET。",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{{Code: "40097", Desc: "jsapi ticket 相关错误", Fix: "检查账号配置。"}},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"ticket 7200s 有效,建议缓存。"},
		Orchestration:  "配合前端 jsapi 签名(window.config)使用。",
		ParsingHints:   "文本输出 ticket;--json 输出 {\"ticket\":..,\"expires_in\":7200}。",
		Examples:       []string{"wechat-gzh-cli jsapi-ticket --json"},
	}))
	return cmd
}

// newHospitalMsgCmd 就医助手消息推送。
func newHospitalMsgCmd() *cobra.Command {
	var bodyJSON string
	cmd := &cobra.Command{
		Use:   "hospital-msg send --body '{...}'",
		Short: "就医助手消息推送",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			if bodyJSON == "" {
				return fmt.Errorf("--body 必填(JSON)")
			}
			body, err := parseJSONMap(bodyJSON)
			if err != nil {
				return err
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				if err := cl.SendChannelMsg(ctx, body); err != nil {
					return err
				}
				fmt.Fprintln(c.OutOrStdout(), "ok")
				return nil
			})
		},
	}
	cmd.Flags().String("body", "", "就医助手推送请求 JSON")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli hospital-msg send", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli hospital-msg send --body '{...}'",
		Summary:        "医疗服务公众号向用户推送就医相关消息(就诊提醒/报告等)。",
		PreConditions:  []string{"已配置鉴权。", "--body 为符合医疗主体要求的请求 JSON。"},
		PostConditions: []string{"stdout 输出 ok;用户收到推送。"},
		Idempotency:    "NO.",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"用户收到一条医疗消息。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "医疗服务权限。",
		RateLimit:      "受消息推送额度限制。",
		RaceConditions: []string{"无。"},
		Orchestration:  "医疗场景定向推送。",
		ParsingHints:   "文本输出 ok。",
		Examples:       []string{`wechat-gzh-cli hospital-msg send --body '{"touser":"o_xxx"}'`},
	}))
	return cmd
}
