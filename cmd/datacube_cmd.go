package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"tx-gongzhonghao-cli/internal/client"
)

// newDatacubeCmd 数据统计(datacube)命令组:循环为每个 datacube 接口生成子命令。
func newDatacubeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "datacube",
		Short: "数据统计分析(datacube)",
		Long:  "公众号运营数据分析:用户/图文/消息/接口数据,均需申请数据权限。",
	}
	for _, ep := range client.DatacubeEndpoints {
		cmd.AddCommand(newDatacubeSubCmd(ep))
	}
	return cmd
}

func newDatacubeSubCmd(ep client.DatacubePath) *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   ep.Name + " --begin-date YYYY-MM-DD --end-date YYYY-MM-DD",
		Short: ep.Summary,
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			begin, _ := c.Flags().GetString("begin-date")
			end, _ := c.Flags().GetString("end-date")
			if begin == "" || end == "" {
				return fmt.Errorf("--begin-date 与 --end-date 必填(YYYY-MM-DD)")
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				rows, err := cl.DatacubeQuery(ctx, ep.Path, begin, end)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(map[string]any{"list": rows})
				}
				if len(rows) == 0 {
					fmt.Fprintln(c.OutOrStdout(), "(no data)")
					return nil
				}
				for _, row := range rows {
					fmt.Fprintln(c.OutOrStdout(), row)
				}
				return nil
			})
		},
	}
	cmd.Flags().String("begin-date", "", "开始日期 YYYY-MM-DD")
	cmd.Flags().String("end-date", "", "结束日期 YYYY-MM-DD")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli datacube "+ep.Name, RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli datacube " + ep.Name + " --begin-date 2026-11-01 --end-date 2026-11-03",
		Summary:        ep.Summary + "。",
		PreConditions:  []string{"已配置鉴权。", "已申请数据接口权限。", "--begin-date/--end-date 必填,格式 YYYY-MM-DD,多为近 3 天窗口。"},
		PostConditions: []string{"stdout 输出数据行(每行一个对象),或 --json 输出 {\"list\":[...]}。"},
		Idempotency:    "YES. 只读。",
		Atomicity:      "Atomic. 单次 POST。",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{{Code: "41030", Desc: "数据接口权限未开通", Fix: "在公众号后台申请数据权限。"}, {Code: "61501", Desc: "日期非法/超窗口", Fix: "检查日期范围(多限近 3 天)。"}},
		Scope:          "数据统计权限。",
		RateLimit:      "受数据接口调用额度限制。",
		RaceConditions: []string{"数据在次日 T+1 才完整,当天查询可能为空。"},
		Orchestration:  "多接口组合可做运营周报;数据 T+1 更新。",
		ParsingHints:   "--json 输出 list 数组。",
		Examples:       []string{"wechat-gzh-cli datacube " + ep.Name + " --begin-date 2026-11-01 --end-date 2026-11-03 --json"},
	}))
	return cmd
}
