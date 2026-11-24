package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"tx-gongzhonghao-cli/internal/client"
)

func newOpenapiCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "openapi",
		Short: "OpenApi 配额与调用治理",
		Long:  "查询/重置 API 调用配额,查询 rid 信息。",
	}
	cmd.AddCommand(newOpenapiQuotaCmd())
	cmd.AddCommand(newOpenapiResetCmd())
	cmd.AddCommand(newOpenapiRidCmd())
	cmd.AddCommand(newOpenapiResetApiCmd())
	return cmd
}

func newOpenapiQuotaCmd() *cobra.Command {
	var jsonOut bool
	var path string
	cmd := &cobra.Command{
		Use:   "quota --path /cgi-bin/xxx",
		Short: "查询 API 调用额度",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			if path == "" {
				return fmt.Errorf("--path 必填(接口请求路径,如 /cgi-bin/message/custom/send)")
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				q, err := cl.GetApiQuota(ctx, path)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(q)
				}
				fmt.Fprintf(c.OutOrStdout(), "daily=%d used=%d remain=%d\n", q.DailyLimit, q.Used, q.Remain)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&path, "path", "", "接口请求路径")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli openapi quota", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli openapi quota --path /cgi-bin/xxx",
		Summary:        "查询指定接口每日调用额度/已用/剩余与频率限制。",
		PreConditions:  []string{"已配置鉴权。", "--path 必填,以 / 开头,如 /cgi-bin/message/custom/send。"},
		PostConditions: []string{"stdout 输出 daily/used/remain。"},
		Idempotency:    "YES. 只读。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{{Code: "76021", Desc: "cgi_path 填错", Fix: "检查 --path 格式。"}, {Code: "76022", Desc: "token 与账号不符", Fix: "检查接口归属。"}},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"额度是实时快照。"},
		Orchestration:  "调用频繁出 45009 前,可先查询额度判断。",
		ParsingHints:   "--json 输出 {\"daily_limit\":..,\"used\":..,\"remain\":..}。",
		Examples:       []string{"wechat-gzh-cli openapi quota --path /cgi-bin/message/custom/send --json"},
	}))
	return cmd
}

func newOpenapiResetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset --appid APPID",
		Short: "重置 API 每日调用次数",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			appid, _ := c.Flags().GetString("appid")
			if appid == "" {
				return fmt.Errorf("--appid 必填")
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				if err := cl.ResetQuota(ctx, appid); err != nil {
					return err
				}
				fmt.Fprintln(c.OutOrStdout(), "ok")
				return nil
			})
		},
	}
	cmd.Flags().String("appid", "", "要清零的账号 appid")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli openapi reset", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli openapi reset --appid APPID",
		Summary:        "清空服务端接口每日调用次数(每账号每月限 10 次)。",
		PreConditions:  []string{"已配置鉴权。", "--appid 必填。"},
		PostConditions: []string{"stdout 输出 ok;每日配额被清零。"},
		Idempotency:    "NO. 清零有副作用,消耗月度机会。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"消耗 1 次月度清零机会。"},
		CommonErrors:   []ErrSpec{{Code: "48006", Desc: "已达清零次数上限", Fix: "下月再试。"}},
		Scope:          "开发者权限。",
		RateLimit:      "每月 10 次。",
		RaceConditions: []string{"清零立即生效。"},
		Orchestration:  "接口触发 45009 后用于应急恢复。",
		ParsingHints:   "文本输出 ok。",
		Examples:       []string{"wechat-gzh-cli openapi reset --appid wx_xxx"},
	}))
	return cmd
}

func newOpenapiRidCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "rid --rid RID",
		Short: "查询 rid 信息",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			rid, _ := c.Flags().GetString("rid")
			if rid == "" {
				return fmt.Errorf("--rid 必填")
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				r, err := cl.GetRidInfo(ctx, rid)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(r)
				}
				fmt.Fprintf(c.OutOrStdout(), "request_url=%s\ncost_ms=%d\nclient_ip=%s\n", r.RequestURL, r.CostInMS, r.ClientIP)
				fmt.Fprintf(c.OutOrStdout(), "response=%s\n", r.ResponseBody)
				return nil
			})
		},
	}
	cmd.Flags().String("rid", "", "接口报错返回的 rid")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli openapi rid", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli openapi rid --rid xxx",
		Summary:        "按 rid 查询某次接口报错的完整请求/响应详情,用于定位问题。",
		PreConditions:  []string{"已配置鉴权。", "--rid 必填且为 7 天内产生的 rid。"},
		PostConditions: []string{"stdout 输出请求 URL/耗时/响应等。"},
		Idempotency:    "YES. 只读。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{{Code: "76001", Desc: "rid 不存在/过期", Fix: "rid 仅 7 天有效。"}},
		Scope:          "开发者权限(同账号)。",
		RateLimit:      "常规。",
		RaceConditions: []string{"rid 过期后无法查询。"},
		Orchestration:  "对应某接口报错后,取报错里的 rid 排查。",
		ParsingHints:   "--json 输出 request 对象。",
		Examples:       []string{"wechat-gzh-cli openapi rid --rid 61725984-xxx --json"},
	}))
	return cmd
}

func newOpenapiResetApiCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset-api --path /channels/ec/...",
		Short: "重置指定 API 调用次数",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			path, _ := c.Flags().GetString("path")
			if path == "" {
				return fmt.Errorf("--path 必填(以 / 开头)")
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				if err := cl.ClearApiQuota(ctx, path); err != nil {
					return err
				}
				fmt.Fprintln(c.OutOrStdout(), "ok")
				return nil
			})
		},
	}
	cmd.Flags().String("path", "", "要清空的接口路径(以 / 开头)")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli openapi reset-api", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli openapi reset-api --path /xxx",
		Summary:        "重置指定接口的每日调用次数(更精准,每账号每月限 50 次)。",
		PreConditions:  []string{"已配置鉴权。", "--path 必填。"},
		PostConditions: []string{"stdout 输出 ok。"},
		Idempotency:    "NO.",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"消耗 1 次月度清零机会。"},
		CommonErrors:   []ErrSpec{{Code: "76021", Desc: "cgi_path 填错", Fix: "检查 --path。"}},
		Scope:          "开发者权限。",
		RateLimit:      "每月 50 次。",
		RaceConditions: []string{"清零立即生效。"},
		Orchestration:  "对单个高耗接口定向清零。",
		ParsingHints:   "文本输出 ok。",
		Examples:       []string{"wechat-gzh-cli openapi reset-api --path /cgi-bin/message/custom/send"},
	}))
	return cmd
}
