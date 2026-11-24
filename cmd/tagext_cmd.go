package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"tx-gongzhonghao-cli/internal/client"
)

func newTagGetFansCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "getfans --tag-id N [--next-openid xxx]",
		Short: "获取标签下粉丝列表",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			tagID, _ := c.Flags().GetInt("tag-id")
			next, _ := c.Flags().GetString("next-openid")
			if tagID == 0 {
				return fmt.Errorf("--tag-id 必填")
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				res, err := cl.TagGetFans(ctx, tagID, next)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(res)
				}
				fmt.Fprintf(c.OutOrStdout(), "count=%d\n", res.Count)
				for _, id := range res.Data.OpenID {
					fmt.Fprintln(c.OutOrStdout(), id)
				}
				if res.NextOpenID != "" {
					fmt.Fprintf(c.OutOrStdout(), "next_openid=%s\n", res.NextOpenID)
				}
				return nil
			})
		},
	}
	cmd.Flags().Int("tag-id", 0, "标签 id")
	cmd.Flags().String("next-openid", "", "翻页游标,用上次返回的 next_openid")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli tag getfans", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli tag getfans --tag-id N --next-openid xxx",
		Summary:        "获取指定标签下的关注者 OpenID 列表,支持 next_openid 翻页。",
		PreConditions:  []string{"已配置鉴权。", "--tag-id 必填。"},
		PostConditions: []string{"stdout 输出 count 与该页 openid;若还有下一页输出 next_openid。"},
		Idempotency:    "YES. 只读。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{{Code: "40003", Desc: "invalid openid", Fix: "检查 openid 来源。"}, {Code: "45159", Desc: "non-existent tag id", Fix: "检查 tag-id。"}},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"分页快照。"},
		Orchestration:  "用于按标签筛选用户后精准群发。",
		ParsingHints:   "--json 输出 {\"count\":..,\"data\":{\"openid\":[..]},\"next_openid\":..}。",
		Examples:       []string{"wechat-gzh-cli tag getfans --tag-id 100 --json"},
	}))
	return cmd
}

func newTagUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update --tag-id N --name 新名",
		Short: "编辑标签",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			tagID, _ := c.Flags().GetInt("tag-id")
			name, _ := c.Flags().GetString("name")
			if tagID == 0 || name == "" {
				return fmt.Errorf("--tag-id 与 --name 必填")
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				if err := cl.TagUpdate(ctx, tagID, name); err != nil {
					return err
				}
				fmt.Fprintln(c.OutOrStdout(), "ok")
				return nil
			})
		},
	}
	cmd.Flags().Int("tag-id", 0, "标签 id")
	cmd.Flags().String("name", "", "新标签名")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli tag update", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli tag update --tag-id N --name 新名",
		Summary:        "修改标签名称。",
		PreConditions:  []string{"已配置鉴权。", "--tag-id/--name 必填。"},
		PostConditions: []string{"stdout 输出 ok;标签名已修改。"},
		Idempotency:    "NO.",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"标签名称被修改。"},
		CommonErrors:   []ErrSpec{{Code: "45157", Desc: "tag name too long", Fix: "限 30 字符。"}},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"无。"},
		Orchestration:  "重命名标签后,原标签 id 不变。",
		ParsingHints:   "文本输出 ok。",
		Examples:       []string{"wechat-gzh-cli tag update --tag-id 100 --name vip2"},
	}))
	return cmd
}

func newTagDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete --tag-id N",
		Short: "删除标签",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			tagID, _ := c.Flags().GetInt("tag-id")
			if tagID == 0 {
				return fmt.Errorf("--tag-id 必填")
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				if err := cl.TagDelete(ctx, tagID); err != nil {
					return err
				}
				fmt.Fprintln(c.OutOrStdout(), "deleted")
				return nil
			})
		},
	}
	cmd.Flags().Int("tag-id", 0, "标签 id")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli tag delete", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli tag delete --tag-id N",
		Summary:        "删除指定标签(标签下用户自动回到未打标状态)。",
		PreConditions:  []string{"已配置鉴权。", "--tag-id 必填。"},
		PostConditions: []string{"stdout 输出 deleted;标签及标记被移除。"},
		Idempotency:    "NO. 有副作用。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"删除标签及其用户关联。"},
		CommonErrors:   []ErrSpec{{Code: "40003", Desc: "invalid openid", Fix: "内部错误。"}},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"删除后群发/菜单引用的 tag_id 失效。"},
		Orchestration:  "删除前确认是否被菜单/群发引用。",
		ParsingHints:   "文本输出 deleted。",
		Examples:       []string{"wechat-gzh-cli tag delete --tag-id 100"},
	}))
	return cmd
}

func newTagBatchCmd(tagging bool) *cobra.Command {
	name := "batch-untag"
	verb := "取消"
	if tagging {
		name = "batch-tag"
		verb = "打上"
	}
	use := name + " --tag-id N --openids o1,o2,..."
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   use,
		Short: "批量" + verb + "标签",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			tagID, _ := c.Flags().GetInt("tag-id")
			ids, _ := c.Flags().GetString("openids")
			if tagID == 0 || ids == "" {
				return fmt.Errorf("--tag-id 与 --openids 必填")
			}
			parts := strings.Split(ids, ",")
			var openids []string
			for _, p := range parts {
				if p = strings.TrimSpace(p); p != "" {
					openids = append(openids, p)
				}
			}
			if len(openids) == 0 {
				return fmt.Errorf("--openids 无有效 openid")
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				var err error
				if tagging {
					err = cl.BatchTagging(ctx, openids, tagID)
				} else {
					err = cl.BatchUntagging(ctx, openids, tagID)
				}
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(map[string]any{"tag_id": tagID, "openid_list": openids, "action": name})
				}
				fmt.Fprintln(c.OutOrStdout(), "ok")
				return nil
			})
		},
	}
	cmd.Flags().Int("tag-id", 0, "标签 id")
	cmd.Flags().String("openids", "", "逗号分隔的 OpenID 列表(最多 50 个)")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli tag "+name, RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli tag " + name + " --tag-id N --openids o1,o2",
		Summary:        "批量为用户" + verb + "指定标签(每次最多 50 个 openid)。",
		PreConditions:  []string{"已配置鉴权。", "--tag-id 与 --openids 必填, openid 以逗号分隔。"},
		PostConditions: []string{"stdout 输出 ok;用户标签已更新。"},
		Idempotency:    "YES. 重复设置相同标签结果一致。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"批量修改用户标签。"},
		CommonErrors:   []ErrSpec{{Code: "40003", Desc: "invalid openid", Fix: "检查 openid 列表。"}},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"无。"},
		Orchestration:  "配合 tag getfans 按标签管理用户。",
		ParsingHints:   "文本输出 ok;--json 输出回显。",
		Examples:       []string{"wechat-gzh-cli tag " + name + " --tag-id 100 --openids o1,o2,o3"},
	}))
	return cmd
}

func newTagOfUserCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "of-user --openid xxx",
		Short: "获取用户标签列表",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			id, _ := c.Flags().GetString("openid")
			if id == "" {
				return fmt.Errorf("--openid 必填")
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				ids, err := cl.GetUserTagList(ctx, id)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(map[string]any{"openid": id, "tagid_list": ids})
				}
				for _, t := range ids {
					fmt.Fprintln(c.OutOrStdout(), t)
				}
				return nil
			})
		},
	}
	cmd.Flags().String("openid", "", "用户 OpenID")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli tag of-user", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli tag of-user --openid xxx",
		Summary:        "查询某用户身上的全部标签 id。",
		PreConditions:  []string{"已配置鉴权。", "--openid 必填。"},
		PostConditions: []string{"stdout 输出该用户的 tagid_list。"},
		Idempotency:    "YES. 只读。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{{Code: "40003", Desc: "invalid openid", Fix: "检查 openid。"}},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"无。"},
		Orchestration:  "查看用户分组以决定打/取标签。",
		ParsingHints:   "--json 输出 {\"openid\":..,\"tagid_list\":[..]}。",
		Examples:       []string{"wechat-gzh-cli tag of-user --openid o_xxx --json"},
	}))
	return cmd
}
