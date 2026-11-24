package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"tx-gongzhonghao-cli/internal/client"
)

// newCustomerAccountCmds 客服账号管理相关子命令。
func newCustomerAccountCmds() []*cobra.Command {
	var cmds []*cobra.Command

	add := func(cmd *cobra.Command) {
		cmds = append(cmds, cmd)
	}

	// account-add
	{
		cmd := &cobra.Command{
			Use:   "account-add --kf-account xxx --nickname 昵称",
			Short: "添加客服账号",
			Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
			RunE: func(c *cobra.Command, _ []string) error {
				if agentHelpShortcut(c) {
					return nil
				}
				kf, _ := c.Flags().GetString("kf-account")
				nick, _ := c.Flags().GetString("nickname")
				if kf == "" || nick == "" {
					return fmt.Errorf("--kf-account 与 --nickname 必填")
				}
				return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
					if err := cl.AddKfAccount(ctx, kf, nick); err != nil {
						return err
					}
					fmt.Fprintln(c.OutOrStdout(), "ok")
					return nil
				})
			},
		}
		cmd.Flags().String("kf-account", "", "完整客服账号,格式 账号@微信号")
		cmd.Flags().String("nickname", "", "客服昵称")
		addAgentHelpFlag(cmd)
		registerAgentHelp("wechat-gzh-cli customer-msg account-add", RenderAgentHelp(AgentHelpDoc{
			Command:        "wechat-gzh-cli customer-msg account-add --kf-account test1@test --nickname 客服1",
			Summary:        "添加客服账号(最多 100 个)。",
			PreConditions:  []string{"已配置鉴权。", "--kf-account 格式 账号@微信号。"},
			PostConditions: []string{"stdout 输出 ok;客服账号创建。"},
			Idempotency:    "NO.",
			Atomicity:      "Atomic.",
			SideEffects:    []string{"新增一个客服账号。"},
			CommonErrors:   []ErrSpec{{Code: "61451", Desc: "已添加客服", Fix: "检查是否重复。"}, {Code: "61452", Desc: "客服名额不足", Fix: "已达 100 上限。"}},
			Scope:          "客服权限。",
			RateLimit:      "常规。",
			RaceConditions: []string{"无。"},
			Orchestration:  "添加后可 account-list 查看。",
			ParsingHints:   "文本输出 ok。",
			Examples:       []string{"wechat-gzh-cli customer-msg account-add --kf-account test1@test --nickname 客服1"},
		}))
		add(cmd)
	}

	// account-update
	{
		cmd := &cobra.Command{
			Use:   "account-update --kf-account xxx --nickname 昵称",
			Short: "修改客服账号",
			Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
			RunE: func(c *cobra.Command, _ []string) error {
				if agentHelpShortcut(c) {
					return nil
				}
				kf, _ := c.Flags().GetString("kf-account")
				nick, _ := c.Flags().GetString("nickname")
				if kf == "" || nick == "" {
					return fmt.Errorf("--kf-account 与 --nickname 必填")
				}
				return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
					if err := cl.UpdateKfAccount(ctx, kf, nick); err != nil {
						return err
					}
					fmt.Fprintln(c.OutOrStdout(), "ok")
					return nil
				})
			},
		}
		cmd.Flags().String("kf-account", "", "完整客服账号")
		cmd.Flags().String("nickname", "", "新昵称")
		addAgentHelpFlag(cmd)
		registerAgentHelp("wechat-gzh-cli customer-msg account-update", RenderAgentHelp(AgentHelpDoc{
			Command:        "wechat-gzh-cli customer-msg account-update --kf-account test1@test --nickname 新名",
			Summary:        "修改客服账号昵称。",
			PreConditions:  []string{"已配置鉴权。", "--kf-account/--nickname 必填。"},
			PostConditions: []string{"stdout 输出 ok。"},
			Idempotency:    "NO.",
			Atomicity:      "Atomic.",
			SideEffects:    []string{"客服昵称被修改。"},
			CommonErrors:   []ErrSpec{{Code: "61453", Desc: "客服未开通", Fix: "先添加。"}},
			Scope:          "客服权限。",
			RateLimit:      "常规。",
			RaceConditions: []string{"无。"},
			Orchestration:  "用于客服资料维护。",
			ParsingHints:   "文本输出 ok。",
			Examples:       []string{"wechat-gzh-cli customer-msg account-update --kf-account test1@test --nickname 新名"},
		}))
		add(cmd)
	}

	// account-del
	{
		cmd := &cobra.Command{
			Use:   "account-del --kf-account xxx",
			Short: "删除客服账号",
			Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
			RunE: func(c *cobra.Command, _ []string) error {
				if agentHelpShortcut(c) {
					return nil
				}
				kf, _ := c.Flags().GetString("kf-account")
				if kf == "" {
					return fmt.Errorf("--kf-account 必填")
				}
				return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
					if err := cl.DelKfAccount(ctx, kf); err != nil {
						return err
					}
					fmt.Fprintln(c.OutOrStdout(), "deleted")
					return nil
				})
			},
		}
		cmd.Flags().String("kf-account", "", "完整客服账号")
		addAgentHelpFlag(cmd)
		registerAgentHelp("wechat-gzh-cli customer-msg account-del", RenderAgentHelp(AgentHelpDoc{
			Command:        "wechat-gzh-cli customer-msg account-del --kf-account test1@test",
			Summary:        "删除客服账号。",
			PreConditions:  []string{"已配置鉴权。", "--kf-account 必填。"},
			PostConditions: []string{"stdout 输出 deleted。"},
			Idempotency:    "NO.",
			Atomicity:      "Atomic.",
			SideEffects:    []string{"客服账号被删除。"},
			CommonErrors:   []ErrSpec{},
			Scope:          "客服权限。",
			RateLimit:      "常规。",
			RaceConditions: []string{"无。"},
			Orchestration:  "删除前确认,避免会话中断。",
			ParsingHints:   "文本输出 deleted。",
			Examples:       []string{"wechat-gzh-cli customer-msg account-del --kf-account test1@test"},
		}))
		add(cmd)
	}

	// account-list
	{
		var jsonOut bool
		cmd := &cobra.Command{
			Use:   "account-list",
			Short: "获取所有客服账号",
			Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
			RunE: func(c *cobra.Command, _ []string) error {
				return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
					res, err := cl.GetKfList(ctx)
					if err != nil {
						return err
					}
					if jsonOut {
						return writeJSON(res)
					}
					for _, k := range res.KFList {
						fmt.Fprintf(c.OutOrStdout(), "%s\t%s\t%d\n", k.KfAccount, k.KfNick, k.KfID)
					}
					return nil
				})
			},
		}
		cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
		addAgentHelpFlag(cmd)
		registerAgentHelp("wechat-gzh-cli customer-msg account-list", RenderAgentHelp(AgentHelpDoc{
			Command:        "wechat-gzh-cli customer-msg account-list",
			Summary:        "获取公众号全部客服账号及其基本信息。",
			PreConditions:  []string{"已配置鉴权。"},
			PostConditions: []string{"stdout 输出 kf_account\\tnick\\tkf_id。"},
			Idempotency:    "YES. 只读。",
			Atomicity:      "Atomic.",
			SideEffects:    []string{"无。"},
			CommonErrors:   []ErrSpec{},
			Scope:          "客服权限。",
			RateLimit:      "常规。",
			RaceConditions: []string{"无。"},
			Orchestration:  "管理客服前先列出现有账号。",
			ParsingHints:   "--json 输出 kf_list。",
			Examples:       []string{"wechat-gzh-cli customer-msg account-list --json"},
		}))
		add(cmd)
	}

	// online-list
	{
		var jsonOut bool
		cmd := &cobra.Command{
			Use:   "online-list",
			Short: "获取在线客服列表",
			Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
			RunE: func(c *cobra.Command, _ []string) error {
				return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
					res, err := cl.GetOnlineKfList(ctx)
					if err != nil {
						return err
					}
					if jsonOut {
						return writeJSON(res)
					}
					for _, k := range res.KFOnlineList {
						fmt.Fprintf(c.OutOrStdout(), "%s\tstatus=%d\n", k.KfAccount, k.Status)
					}
					return nil
				})
			},
		}
		cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
		addAgentHelpFlag(cmd)
		registerAgentHelp("wechat-gzh-cli customer-msg online-list", RenderAgentHelp(AgentHelpDoc{
			Command:        "wechat-gzh-cli customer-msg online-list",
			Summary:        "获取当前在线客服及其状态。",
			PreConditions:  []string{"已配置鉴权。"},
			PostConditions: []string{"stdout 输出在线客服与状态。"},
			Idempotency:    "YES. 只读。",
			Atomicity:      "Atomic.",
			SideEffects:    []string{"无。"},
			CommonErrors:   []ErrSpec{},
			Scope:          "客服权限。",
			RateLimit:      "常规。",
			RaceConditions: []string{"在线状态实时变化。"},
			Orchestration:  "用于客服调度判断。",
			ParsingHints:   "--json 输出 kf_online_list。",
			Examples:       []string{"wechat-gzh-cli customer-msg online-list --json"},
		}))
		add(cmd)
	}

	// inviteworker
	{
		cmd := &cobra.Command{
			Use:   "inviteworker --kf-account xxx --invite-wx 微信",
			Short: "邀请绑定客服",
			Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
			RunE: func(c *cobra.Command, _ []string) error {
				if agentHelpShortcut(c) {
					return nil
				}
				kf, _ := c.Flags().GetString("kf-account")
				wx, _ := c.Flags().GetString("invite-wx")
				if kf == "" || wx == "" {
					return fmt.Errorf("--kf-account 与 --invite-wx 必填")
				}
				return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
					if err := cl.InviteKfWorker(ctx, kf, wx); err != nil {
						return err
					}
					fmt.Fprintln(c.OutOrStdout(), "ok")
					return nil
				})
			},
		}
		cmd.Flags().String("kf-account", "", "完整客服账号")
		cmd.Flags().String("invite-wx", "", "被邀请的微信账号")
		addAgentHelpFlag(cmd)
		registerAgentHelp("wechat-gzh-cli customer-msg inviteworker", RenderAgentHelp(AgentHelpDoc{
			Command:        "wechat-gzh-cli customer-msg inviteworker --kf-account test1@test --invite-wx wxid",
			Summary:        "邀请微信用户成为客服绑定账号。",
			PreConditions:  []string{"已配置鉴权。", "--kf-account/--invite-wx 必填。"},
			PostConditions: []string{"stdout 输出 ok;被邀请人收到邀请。"},
			Idempotency:    "NO.",
			Atomicity:      "Atomic.",
			SideEffects:    []string{"发出邀请。"},
			CommonErrors:   []ErrSpec{},
			Scope:          "客服权限。",
			RateLimit:      "常规。",
			RaceConditions: []string{"无。"},
			Orchestration:  "邀请后需对方确认绑定。",
			ParsingHints:   "文本输出 ok。",
			Examples:       []string{"wechat-gzh-cli customer-msg inviteworker --kf-account test1@test --invite-wx wxid"},
		}))
		add(cmd)
	}

	return cmds
}

// newCustomerSessionCmds 会话控制相关子命令。
func newCustomerSessionCmds() []*cobra.Command {
	var cmds []*cobra.Command

	// session-create
	{
		cmd := &cobra.Command{
			Use:   "session-create --kf-account xxx --openid o_xxx",
			Short: "创建会话",
			Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
			RunE: func(c *cobra.Command, _ []string) error {
				if agentHelpShortcut(c) {
					return nil
				}
				kf, _ := c.Flags().GetString("kf-account")
				id, _ := c.Flags().GetString("openid")
				if kf == "" || id == "" {
					return fmt.Errorf("--kf-account 与 --openid 必填")
				}
				return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
					if err := cl.CreateKfSession(ctx, kf, id); err != nil {
						return err
					}
					fmt.Fprintln(c.OutOrStdout(), "ok")
					return nil
				})
			},
		}
		cmd.Flags().String("kf-account", "", "客服账号")
		cmd.Flags().String("openid", "", "用户 openid")
		addAgentHelpFlag(cmd)
		registerAgentHelp("wechat-gzh-cli customer-msg session-create", RenderAgentHelp(AgentHelpDoc{
			Command:        "wechat-gzh-cli customer-msg session-create --kf-account test1@test --openid o_xxx",
			Summary:        "创建客服与用户之间的会话(需用户先发过消息)。",
			PreConditions:  []string{"已配置鉴权。", "用户与该客服账号发送过消息。", "--kf-account/--openid 必填。"},
			PostConditions: []string{"stdout 输出 ok;会话建立。"},
			Idempotency:    "NO.",
			Atomicity:      "Atomic.",
			SideEffects:    []string{"创建一段客服会话。"},
			CommonErrors:   []ErrSpec{{Code: "65401", Desc: "openid 错误", Fix: "检查 openid。"}},
			Scope:          "客服权限。",
			RateLimit:      "常规。",
			RaceConditions: []string{"无。"},
			Orchestration:  "会话建立后可发送客服消息。",
			ParsingHints:   "文本输出 ok。",
			Examples:       []string{"wechat-gzh-cli customer-msg session-create --kf-account test1@test --openid o_xxx"},
		}))
		cmds = append(cmds, cmd)
	}

	// session-close
	{
		cmd := &cobra.Command{
			Use:   "session-close --kf-account xxx --openid o_xxx",
			Short: "关闭会话",
			Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
			RunE: func(c *cobra.Command, _ []string) error {
				if agentHelpShortcut(c) {
					return nil
				}
				kf, _ := c.Flags().GetString("kf-account")
				id, _ := c.Flags().GetString("openid")
				if kf == "" || id == "" {
					return fmt.Errorf("--kf-account 与 --openid 必填")
				}
				return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
					if err := cl.CloseKfSession(ctx, kf, id); err != nil {
						return err
					}
					fmt.Fprintln(c.OutOrStdout(), "ok")
					return nil
				})
			},
		}
		cmd.Flags().String("kf-account", "", "客服账号")
		cmd.Flags().String("openid", "", "用户 openid")
		addAgentHelpFlag(cmd)
		registerAgentHelp("wechat-gzh-cli customer-msg session-close", RenderAgentHelp(AgentHelpDoc{
			Command:        "wechat-gzh-cli customer-msg session-close --kf-account test1@test --openid o_xxx",
			Summary:        "关闭客服与用户之间的会话。",
			PreConditions:  []string{"已配置鉴权。", "--kf-account/--openid 必填。"},
			PostConditions: []string{"stdout 输出 ok;会话结束。"},
			Idempotency:    "NO.",
			Atomicity:      "Atomic.",
			SideEffects:    []string{"结束会话。"},
			CommonErrors:   []ErrSpec{},
			Scope:          "客服权限。",
			RateLimit:      "常规。",
			RaceConditions: []string{"无。"},
			Orchestration:  "服务完成后可关闭会话。",
			ParsingHints:   "文本输出 ok。",
			Examples:       []string{"wechat-gzh-cli customer-msg session-close --kf-account test1@test --openid o_xxx"},
		}))
		cmds = append(cmds, cmd)
	}

	// session-get
	{
		var jsonOut bool
		cmd := &cobra.Command{
			Use:   "session-get --openid o_xxx",
			Short: "获取客户会话状态",
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
					res, err := cl.GetKfSession(ctx, id)
					if err != nil {
						return err
					}
					if jsonOut {
						return writeJSON(res)
					}
					fmt.Fprintf(c.OutOrStdout(), "%s\t%d\n", res.KfAccount, res.CreateTime)
					return nil
				})
			},
		}
		cmd.Flags().String("openid", "", "用户 openid")
		cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
		addAgentHelpFlag(cmd)
		registerAgentHelp("wechat-gzh-cli customer-msg session-get", RenderAgentHelp(AgentHelpDoc{
			Command:        "wechat-gzh-cli customer-msg session-get --openid o_xxx",
			Summary:        "获取用户当前会话所属客服与创建时间。",
			PreConditions:  []string{"已配置鉴权。", "--openid 必填。"},
			PostConditions: []string{"stdout 输出 kf_account\\tcreatetime。"},
			Idempotency:    "YES. 只读。",
			Atomicity:      "Atomic.",
			SideEffects:    []string{"无。"},
			CommonErrors:   []ErrSpec{},
			Scope:          "客服权限。",
			RateLimit:      "常规。",
			RaceConditions: []string{"无。"},
			Orchestration:  "判断用户是否需要转接客服。",
			ParsingHints:   "--json 输出 kf_account/createtime。",
			Examples:       []string{"wechat-gzh-cli customer-msg session-get --openid o_xxx --json"},
		}))
		cmds = append(cmds, cmd)
	}

	// session-list
	{
		var jsonOut bool
		cmd := &cobra.Command{
			Use:   "session-list --kf-account xxx",
			Short: "获取客服会话列表",
			Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
			RunE: func(c *cobra.Command, _ []string) error {
				if agentHelpShortcut(c) {
					return nil
				}
				kf, _ := c.Flags().GetString("kf-account")
				if kf == "" {
					return fmt.Errorf("--kf-account 必填")
				}
				return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
					list, err := cl.GetKfSessionList(ctx, kf)
					if err != nil {
						return err
					}
					if jsonOut {
						return writeJSON(map[string]any{"sessionlist": list})
					}
					for _, s := range list {
						fmt.Fprintf(c.OutOrStdout(), "%s\t%d\n", s.OpenID, s.CreateTime)
					}
					return nil
				})
			},
		}
		cmd.Flags().String("kf-account", "", "客服账号")
		cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
		addAgentHelpFlag(cmd)
		registerAgentHelp("wechat-gzh-cli customer-msg session-list", RenderAgentHelp(AgentHelpDoc{
			Command:        "wechat-gzh-cli customer-msg session-list --kf-account test1@test",
			Summary:        "获取指定客服正在接待的会话列表。",
			PreConditions:  []string{"已配置鉴权。", "--kf-account 必填。"},
			PostConditions: []string{"stdout 输出 openid\\tcreatetime。"},
			Idempotency:    "YES. 只读。",
			Atomicity:      "Atomic.",
			SideEffects:    []string{"无。"},
			CommonErrors:   []ErrSpec{},
			Scope:          "客服权限。",
			RateLimit:      "常规。",
			RaceConditions: []string{"实时快照。"},
			Orchestration:  "客服工作台负载查询。",
			ParsingHints:   "--json 输出 sessionlist。",
			Examples:       []string{"wechat-gzh-cli customer-msg session-list --kf-account test1@test"},
		}))
		cmds = append(cmds, cmd)
	}

	// session-waiting
	{
		var jsonOut bool
		cmd := &cobra.Command{
			Use:   "session-waiting",
			Short: "获取未接入会话列表",
			Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
			RunE: func(c *cobra.Command, _ []string) error {
				return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
					list, err := cl.GetWaitCase(ctx)
					if err != nil {
						return err
					}
					if jsonOut {
						return writeJSON(map[string]any{"waitcaselist": list})
					}
					for _, s := range list {
						fmt.Fprintf(c.OutOrStdout(), "%s\t%d\n", s.OpenID, s.CreateTime)
					}
					return nil
				})
			},
		}
		cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
		addAgentHelpFlag(cmd)
		registerAgentHelp("wechat-gzh-cli customer-msg session-waiting", RenderAgentHelp(AgentHelpDoc{
			Command:        "wechat-gzh-cli customer-msg session-waiting",
			Summary:        "获取未接入人工客服的等待会话列表。",
			PreConditions:  []string{"已配置鉴权。"},
			PostConditions: []string{"stdout 输出等待会话 openid。"},
			Idempotency:    "YES. 只读。",
			Atomicity:      "Atomic.",
			SideEffects:    []string{"无。"},
			CommonErrors:   []ErrSpec{},
			Scope:          "客服权限。",
			RateLimit:      "常规。",
			RaceConditions: []string{"实时快照。"},
			Orchestration:  "客服抢单/分配依据。",
			ParsingHints:   "--json 输出 waitcaselist。",
			Examples:       []string{"wechat-gzh-cli customer-msg session-waiting --json"},
		}))
		cmds = append(cmds, cmd)
	}

	return cmds
}

// newCustomerMsgExtCmds 客服消息扩展(聊天记录/输入状态)。
func newCustomerMsgExtCmds() []*cobra.Command {
	var cmds []*cobra.Command

	// msg-record
	{
		var jsonOut bool
		cmd := &cobra.Command{
			Use:   "msg-record --starttime TS --endtime TS --msgid MSGID --number N",
			Short: "获取聊天记录",
			Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
			RunE: func(c *cobra.Command, _ []string) error {
				if agentHelpShortcut(c) {
					return nil
				}
				start, _ := c.Flags().GetInt64("starttime")
				end, _ := c.Flags().GetInt64("endtime")
				msgid, _ := c.Flags().GetInt64("msgid")
				number, _ := c.Flags().GetInt("number")
				if start == 0 || end == 0 || number == 0 {
					return fmt.Errorf("--starttime/--endtime/--number 必填")
				}
				return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
					list, err := cl.GetMsgList(ctx, start, end, msgid, number)
					if err != nil {
						return err
					}
					if jsonOut {
						return writeJSON(map[string]any{"recordlist": list})
					}
					for _, r := range list {
						fmt.Fprintln(c.OutOrStdout(), r)
					}
					return nil
				})
			},
		}
		cmd.Flags().Int64("starttime", 0, "起始时间(Unix 秒)")
		cmd.Flags().Int64("endtime", 0, "结束时间(Unix 秒)")
		cmd.Flags().Int64("msgid", 0, "上一次返回的 msgid,用于翻页;首页为 0")
		cmd.Flags().Int("number", 10000, "数量(1-10000)")
		cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
		addAgentHelpFlag(cmd)
		registerAgentHelp("wechat-gzh-cli customer-msg msg-record", RenderAgentHelp(AgentHelpDoc{
			Command:        "wechat-gzh-cli customer-msg msg-record --starttime 1667 --endtime 1668 --msgid 0 --number 10000",
			Summary:        "获取指定时间段内的客服聊天记录。",
			PreConditions:  []string{"已配置鉴权。", "--starttime/--endtime/--number 必填,msgid=0 起翻页。"},
			PostConditions: []string{"stdout 输出聊天记录。"},
			Idempotency:    "YES. 只读。",
			Atomicity:      "Atomic.",
			SideEffects:    []string{"无。"},
			CommonErrors:   []ErrSpec{{Code: "60009", Desc: "时间范围超限", Fix: "单次查询时间跨度有限。"}},
			Scope:          "客服权限。",
			RateLimit:      "常规。",
			RaceConditions: []string{"翻页用 msgid 游标。"},
			Orchestration:  "用于质检/合规留档。",
			ParsingHints:   "--json 输出 recordlist。",
			Examples:       []string{"wechat-gzh-cli customer-msg msg-record --starttime 1667600000 --endtime 1667686400 --number 10000 --json"},
		}))
		cmds = append(cmds, cmd)
	}

	// typing
	{
		cmd := &cobra.Command{
			Use:   "typing --openid o_xxx --command Typing",
			Short: "显示客服输入状态",
			Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
			RunE: func(c *cobra.Command, _ []string) error {
				if agentHelpShortcut(c) {
					return nil
				}
				id, _ := c.Flags().GetString("openid")
				command, _ := c.Flags().GetString("command")
				if id == "" || command == "" {
					return fmt.Errorf("--openid 与 --command 必填")
				}
				return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
					if err := cl.SendTyping(ctx, id, command); err != nil {
						return err
					}
					fmt.Fprintln(c.OutOrStdout(), "ok")
					return nil
				})
			},
		}
		cmd.Flags().String("openid", "", "用户 openid")
		cmd.Flags().String("command", "Typing", "Typing 或 CancelTyping")
		addAgentHelpFlag(cmd)
		registerAgentHelp("wechat-gzh-cli customer-msg typing", RenderAgentHelp(AgentHelpDoc{
			Command:        "wechat-gzh-cli customer-msg typing --openid o_xxx --command Typing",
			Summary:        "向用户显示客服正在输入的状态(或取消)。",
			PreConditions:  []string{"已配置鉴权。", "--openid/--command 必填,command 为 Typing/CancelTyping。"},
			PostConditions: []string{"stdout 输出 ok;用户端显示输入状态。"},
			Idempotency:    "NO.",
			Atomicity:      "Atomic.",
			SideEffects:    []string{"用户会话显示输入指示。"},
			CommonErrors:   []ErrSpec{},
			Scope:          "客服权限。",
			RateLimit:      "常规。",
			RaceConditions: []string{"48 小时内有效。"},
			Orchestration:  "发送客服消息前先提示正在输入。",
			ParsingHints:   "文本输出 ok。",
			Examples:       []string{"wechat-gzh-cli customer-msg typing --openid o_xxx --command Typing"},
		}))
		cmds = append(cmds, cmd)
	}

	return cmds
}
