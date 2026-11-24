package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"tx-gongzhonghao-cli/internal/client"
)

func newStableTokenCmd() *cobra.Command {
	var jsonOut bool
	var forceRefresh bool
	cmd := &cobra.Command{
		Use:   "stable-token",
		Short: "获取稳定版接口调用凭据",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			var fr *bool
			if c.Flags().Changed("force-refresh") {
				fr = &forceRefresh
			}
			cl, err := client.FromEnv()
			if err != nil {
				return err
			}
			token, expires, err := cl.StableAccessToken(context.Background(), fr)
			if err != nil {
				return err
			}
			if jsonOut {
				return writeJSON(map[string]any{"access_token": token, "expires_in": expires})
			}
			fmt.Fprintf(c.OutOrStdout(), "%s\n", token)
			return nil
		},
	}
	cmd.Flags().BoolVar(&forceRefresh, "force-refresh", false, "强制刷新(使之前 token 失效,每日限 20 次)")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli stable-token", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli stable-token [--force-refresh]",
		Summary:        "获取稳定版接口调用凭据(POST /cgi-bin/stable_token),与 getAccessToken 隔离且更稳定。",
		PreConditions:  []string{"已配置 WX_GZ_AppID/WX_GZ_AppSecret。"},
		PostConditions: []string{"stdout 输出 access_token(默认文本/--json 带 expires_in)。"},
		Idempotency:    "YES(普通模式)。--force-refresh 会强制刷新并使旧 token 失效。",
		Atomicity:      "Atomic. 单次 POST。",
		SideEffects:    []string{"普通模式无副作用;force-refresh 会使旧 token 失效。"},
		CommonErrors:   []ErrSpec{{Code: "40013", Desc: "invalid appid", Fix: "检查 AppID。"}, {Code: "40125", Desc: "invalid secret", Fix: "检查 AppSecret。"}},
		Scope:          "开发者权限。",
		RateLimit:      "1 万次/分钟、每日 50 万次;force-refresh 每日 20 次且间隔 30 秒。",
		RaceConditions: []string{"强制刷新有日限,勿频繁调用。"},
		Orchestration:  "与 access-token 二选一用,推荐本命令。",
		ParsingHints:   "--json 输出 {\"access_token\":..,\"expires_in\":7200}。",
		Examples:       []string{"wechat-gzh-cli stable-token --json"},
	}))
	return cmd
}

func newCallbackCheckCmd() *cobra.Command {
	var jsonOut bool
	var operator string
	cmd := &cobra.Command{
		Use:   "callback-check [--action all] [--operator DEFAULT]",
		Short: "网络通信检测(回调连通性)",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			action, _ := c.Flags().GetString("action")
			if action == "" {
				action = "all"
			}
			if operator == "" {
				operator = "DEFAULT"
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				res, err := cl.CallbackCheck(ctx, action, operator)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(res)
				}
				for _, d := range res.DNS {
					fmt.Fprintf(c.OutOrStdout(), "dns: %s %s\n", d.IP, d.RealOperator)
				}
				for _, p := range res.Ping {
					fmt.Fprintf(c.OutOrStdout(), "ping: %s loss=%s time=%s\n", p.IP, p.PackageLoss, p.Time)
				}
				return nil
			})
		},
	}
	cmd.Flags().String("action", "all", "检测动作: dns/ping/all")
	cmd.Flags().StringVar(&operator, "operator", "DEFAULT", "运营商: CHINANET/UNICOM/CAP/DEFAULT")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli callback-check", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli callback-check --action all",
		Summary:        "对开发者 URL 做域名解析与 ping 检测,排查回调连接失败。",
		PreConditions:  []string{"已配置鉴权。", "已在后台配置回调 URL(否则 40201)。"},
		PostConditions: []string{"stdout 输出 dns 与 ping 检测结果。"},
		Idempotency:    "YES. 只读诊断。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{{Code: "40201", Desc: "未设置回调 URL", Fix: "先在后台配置。"}, {Code: "40202", Desc: "action 不合法", Fix: "检查 --action。"}},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"检测为瞬时快照。"},
		Orchestration:  "回调异常时先跑本命令定位连通性问题。",
		ParsingHints:   "--json 输出 {\"dns\":[..],\"ping\":[..]}。",
		Examples:       []string{"wechat-gzh-cli callback-check --json"},
	}))
	return cmd
}

func newPushServerIPCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "push-server-ip",
		Short: "获取微信推送服务器 IP",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				ips, err := cl.GetCallbackIP(ctx)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(map[string]any{"ip_list": ips})
				}
				for _, ip := range ips {
					fmt.Fprintln(c.OutOrStdout(), ip)
				}
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli push-server-ip", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli push-server-ip",
		Summary:        "获取微信推送服务器(向开发者服务器推送消息)的 IP 列表,用于放行与比对。",
		PreConditions:  []string{"已配置鉴权。"},
		PostConditions: []string{"stdout 输出 IP 列表。"},
		Idempotency:    "YES. 只读。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限。",
		RateLimit:      "建议每日请求一次。",
		RaceConditions: []string{"IP 可能变动,勿长期缓存。"},
		Orchestration:  "用于回调服务器防火墙放行。",
		ParsingHints:   "--json 输出 {\"ip_list\":[..]}。",
		Examples:       []string{"wechat-gzh-cli push-server-ip --json"},
	}))
	return cmd
}
