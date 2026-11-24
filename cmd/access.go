package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"tx-gongzhonghao-cli/internal/client"
)

func newAccessCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "access-token",
		Short: "获取接口调用凭据 access_token",
		Long:  "调用 GET /cgi-bin/token 获取全局唯一接口调用凭据 access_token(有效期 7200 秒)。",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				token, err := cl.AccessToken(ctx)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(map[string]any{"access_token": token, "expires_in": 7200})
				}
				fmt.Fprintln(c.OutOrStdout(), token)
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli access-token", RenderAgentHelp(AgentHelpDoc{
		Command: "wechat-gzh-cli access-token",
		Summary: "获取全局接口调用凭据 access_token(有效期 7200 秒),供后续所有服务端接口鉴权。",
		PreConditions: []string{
			"环境变量 WX_GZ_AppID 与 WX_GZ_AppSecret 已设置且有效。",
			"IP 在公众号后台的 IP 白名单内(若后台开启了 IP 限制)。",
		},
		PostConditions: []string{
			"stdout 输出 access_token(默认纯文本,--json 输出 JSON)。",
			"获得的 access_token 在 7200 秒内有效,内部会缓存,重复调用不重复请求。",
		},
		Idempotency: "YES. 同一凭证有效期内重复调用返回同一/等效 token,不产生副作用(普通模式不刷新)。",
		Atomicity:   "Atomic. 单次 GET 请求,失败不产生任何状态变更。",
		SideEffects: []string{"无外部可见副作用;获取的是只读调用凭据。"},
		CommonErrors: []ErrSpec{
			{Code: "40013", Desc: "invalid appid", Fix: "检查 WX_GZ_AppID。"},
			{Code: "40125", Desc: "invalid secret", Fix: "检查 WX_GZ_AppSecret。"},
			{Code: "40164", Desc: "IP not in whitelist", Fix: "把当前出口 IP 加入后台白名单。"},
		},
		Scope:     "仅需公众号后台的开发者权限(无需额外 OAuth 网页授权 scope)。",
		RateLimit: "受公众号 access_token 获取频率限制(普通模式 2000 次/日,稳定版 1 万次/分钟);本工具做了缓存,通常无需担心。",
		RaceConditions: []string{
			"多进程同时首次获取时可能各取一个 token,两者均有效,不冲突。",
		},
		Orchestration: "获取后应把 token 传给后续 API 调用,或直接使用其它命令(它们内部会自动取 token)。",
		ParsingHints:  "--json 输出形如 {\"access_token\":\"...\",\"expires_in\":7200};纯文本模式 stdout 第一行即 token。",
		Examples: []string{
			"wechat-gzh-cli access-token",
			"wechat-gzh-cli access-token --json",
		},
	}))
	return cmd
}

func newServerIPCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "server-ip",
		Short: "获取微信 API 服务器 IP 列表",
		Long:  "调用 GET /cgi-bin/get_api_domain_ip 获取开发者服务器访问 api.weixin.qq.com 的远端 IP。",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				ips, err := cl.GetApiDomainIP(ctx)
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
	registerAgentHelp("wechat-gzh-cli server-ip", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli server-ip",
		Summary:        "获取微信 API 服务器(api.weixin.qq.com)IP 地址列表,用于防火墙放行或 IP 比对。",
		PreConditions:  []string{"已配置 WX_GZ_AppID/WX_GZ_AppSecret。"},
		PostConditions: []string{"stdout 输出 IP 列表(每行一个,或 --json 输出 {\"ip_list\":[...]})。"},
		Idempotency:    "YES. 只读查询,不产生副作用。",
		Atomicity:      "Atomic. 单次 GET。",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{{Code: "44002", Desc: "empty post data", Fix: "重试。"}},
		Scope:          "开发者权限即可。",
		RateLimit:      "常规;建议每日调用一次以获得最新 IP 列表。",
		RaceConditions: []string{"IP 列表可能随时间变化,勿长期缓存。"},
		Orchestration:  "可据此配置防火墙/白名单后,再调用其它业务接口。",
		ParsingHints:   "--json 输出 {\"ip_list\":[\"1.2.3.4\",...]}。",
		Examples:       []string{"wechat-gzh-cli server-ip --json"},
	}))
	return cmd
}

// runWithClient 从环境变量构造 client,把真实鉴权放在命令逻辑内(而非全局钩子),
// 从而保证 --agent-help 在未配置凭证时也能可见(法则二)。
func runWithClient(c *cobra.Command, fn func(ctx context.Context, cl *client.Client) error) error {
	if agentHelpShortcut(c) {
		return nil
	}
	cl, err := client.FromEnv()
	if err != nil {
		return err
	}
	return fn(context.Background(), cl)
}

func writeJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// parseJSONMap 把字符串反序列化为 map;非法 JSON 返回带说明的错误。
func parseJSONMap(s string) (map[string]any, error) {
	var m map[string]any
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return nil, fmt.Errorf("不是合法 JSON: %w", err)
	}
	return m, nil
}

// printMapFlat 以 "key=value" 形式扁平打印 map 顶层字段,用于结构化结果的文本展示。
// 嵌套值保留为格式化 JSON。
func printMapFlat(c *cobra.Command, m map[string]any) {
	for k, v := range m {
		switch val := v.(type) {
		case string:
			fmt.Fprintf(c.OutOrStdout(), "%s=%s\n", k, val)
		case float64:
			fmt.Fprintf(c.OutOrStdout(), "%s=%v\n", k, val)
		case bool:
			fmt.Fprintf(c.OutOrStdout(), "%s=%v\n", k, val)
		default:
			b, err := json.Marshal(v)
			if err != nil {
				fmt.Fprintf(c.OutOrStdout(), "%s=%v\n", k, v)
				continue
			}
			fmt.Fprintf(c.OutOrStdout(), "%s=%s\n", k, string(b))
		}
	}
}
