package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"tx-gongzhonghao-cli/internal/client"
)

// newStoreCmd 微信门店命令组(门店小程序 + 门店 POI)。
func newStoreCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "store",
		Short: "微信门店",
		Long:  "门店小程序入驻与门店(POI)管理。",
	}
	for _, ep := range client.WxaEndpoints {
		cmd.AddCommand(newStoreSubCmd(ep))
	}
	return cmd
}

// newStoreSubCmd 为每个门店接口生成子命令。
func newStoreSubCmd(ep client.WxaEndpoint) *cobra.Command {
	var jsonOut bool
	var bodyJSON string
	use := ep.Name
	if ep.NeedBody {
		use += " --body '{...}'"
	}
	cmd := &cobra.Command{
		Use:   use,
		Short: ep.Summary,
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			var body map[string]any
			if ep.NeedBody {
				if bodyJSON == "" {
					return fmt.Errorf("--body 必填(%s)", ep.BodyHint)
				}
				var err error
				body, err = parseJSONMap(bodyJSON)
				if err != nil {
					return fmt.Errorf("--body %s", err)
				}
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				res, err := cl.WxaQuery(ctx, ep.Path, body)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(res)
				}
				printMapFlat(c, res)
				return nil
			})
		},
	}
	if ep.NeedBody {
		cmd.Flags().String("body", "", "请求体 JSON("+ep.BodyHint+")")
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)

	example := "wechat-gzh-cli store " + ep.Name
	if ep.NeedBody {
		example += " --body '{\"poi_id\":\"P001\"}'"
	}
	example += " --json"

	preCond := "无需请求体。"
	if ep.NeedBody {
		preCond = "--body 必填(" + ep.BodyHint + ")。"
	}
	sideEffect := "无。"
	if ep.NeedBody {
		sideEffect = "可能创建/修改/删除门店或小程序数据。"
	}
	registerAgentHelp("wechat-gzh-cli store "+ep.Name, RenderAgentHelp(AgentHelpDoc{
		Command:        example,
		Summary:        ep.Summary + "。",
		PreConditions:  []string{"已配置鉴权。", preCond},
		PostConditions: []string{"stdout 输出接口返回字段,或 --json 输出完整结构。"},
		Idempotency:    "见接口语义:查询类只读 YES;写操作类 NO。",
		Atomicity:      "Atomic. 单次 POST。",
		SideEffects:    []string{sideEffect},
		CommonErrors:   []ErrSpec{},
		Scope:          "微信门店/小程序权限。",
		RateLimit:      "受门店接口配额限制。",
		RaceConditions: []string{"门店数据变更需以实时查询为准。"},
		Orchestration:  "门店管理按接口语义组合使用。",
		ParsingHints:   "--json 输出完整响应体。",
		Examples:       []string{example},
	}))
	return cmd
}
