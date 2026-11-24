package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"tx-gongzhonghao-cli/internal/client"
)

func newMessageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "message",
		Short: "群发与消息状态",
		Long:  "群发消息与查询群发状态。",
	}
	cmd.AddCommand(newMessageMassSendCmd())
	cmd.AddCommand(newMessageMassStatusCmd())
	cmd.AddCommand(newMessageMassDeleteCmd())
	cmd.AddCommand(newMessageSpeedGetCmd())
	cmd.AddCommand(newMessageSpeedSetCmd())
	cmd.AddCommand(newMessagePreviewCmd())
	return cmd
}

func newMessageMassSendCmd() *cobra.Command {
	var jsonOut bool
	var bodyJSON string
	cmd := &cobra.Command{
		Use:   "mass-send --body '{...}'",
		Short: "根据标签群发消息",
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
				res, err := cl.MassSend(ctx, body)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(res)
				}
				fmt.Fprintf(c.OutOrStdout(), "%d\n", res.MsgID)
				return nil
			})
		},
	}
	cmd.Flags().String("body", "", "群发请求体 JSON(含 filter/msgtype/消息体)")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli message mass-send", RenderAgentHelp(AgentHelpDoc{
		Command: "wechat-gzh-cli message mass-send --body '{\"filter\":{\"is_to_all\":true},\"msgtype\":\"text\",\"text\":{\"content\":\"hi\"}}'",
		Summary: "根据标签(或全部)向用户群发消息,支持 mpnews/text/voice/images/mpvideo/wxcard。",
		PreConditions: []string{
			"已配置鉴权。",
			"公众号仅认证的服务号拥有群发权限。",
			"--body 为符合公众号群发接口的 JSON 请求体。",
		},
		PostConditions: []string{"stdout 输出 msg_id;群发任务已提交(异步真正发送)。"},
		Idempotency:    "NO. 每次提交产生新群发任务;可用 --body 内 clientmsgid 防重复推送。",
		Atomicity:      "Not atomic. 提交成功仅代表任务入队,实际发送异步(可能部分失败)。",
		SideEffects:    []string{"向用户推送消息;可能触发管理员确认(若开启 API 群发保护)。"},
		CommonErrors: []ErrSpec{
			{Code: "45009", Desc: "reach max api daily quota", Fix: "当日发表额度用尽,明日重试。"},
			{Code: "45028", Desc: "no masssend quota", Fix: "周期内发表次数用完。"},
			{Code: "40007", Desc: "invalid media_id", Fix: "图文 media_id 无效。"},
		},
		Scope:        "需要公众号服务号的群发/消息权限。",
		RateLimit:    "受每日发表次数与频率限制(详见错误码 45009/45011)。",
		RaceConditions: []string{
			"群完结果经回调推送 MASSSENDJOBFINISH,勿轮询误判;msg_data_id 用于图文数据接口。",
		},
		Orchestration: "用 message mass-status 查询该 msg_id 的最终发送状态。",
		ParsingHints:  "文本输出 msg_id;--json 输出 {\"msg_id\":..,\"msg_data_id\":..}。",
		Examples: []string{
			`wechat-gzh-cli message mass-send --body '{"filter":{"is_to_all":true},"msgtype":"text","text":{"content":"hi"}}' --json`,
		},
	}))
	return cmd
}

func newMessageMassStatusCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "mass-status --msg-id 12345",
		Short: "查询群发消息发送状态",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			msgID, _ := c.Flags().GetInt64("msg-id")
			if msgID == 0 {
				return fmt.Errorf("--msg-id 必填")
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				st, err := cl.MassStatus(ctx, msgID)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(map[string]any{"msg_id": msgID, "msg_status": st})
				}
				fmt.Fprintln(c.OutOrStdout(), st)
				return nil
			})
		},
	}
	cmd.Flags().Int64("msg-id", 0, "群发任务 msg_id")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli message mass-status", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli message mass-status --msg-id 12345",
		Summary:        "查询群发消息发送状态(SEND_SUCCESS/SENDING/SEND_FAIL/DELETE)。",
		PreConditions:  []string{"已配置鉴权。", "--msg-id 为群发返回的 msg_id。"},
		PostConditions: []string{"stdout 输出 msg_status。"},
		Idempotency:    "YES. 只读查询。",
		Atomicity:      "Atomic. 单次 POST。",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{{Code: "40059", Desc: "invalid msg id", Fix: "检查 msg_id。"}},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"发送中(SENDING)状态会变化,需多次查询。"},
		Orchestration:  "群发后轮询本命令直至 SEND_SUCCESS 或 SEND_FAIL。",
		ParsingHints:   "文本输出状态字符串;--json 输出 {\"msg_id\":..,\"msg_status\":..}。",
		Examples:       []string{"wechat-gzh-cli message mass-status --msg-id 12345 --json"},
	}))
	return cmd
}
