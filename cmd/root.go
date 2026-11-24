package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// AppVersion 版本号,与 git tag 保持一致(仅末位递增)。
var AppVersion = "v0.0.3"

// agentDocRegistry 记录每个已注册命令的 --agent-help 状态,供根命令总览。
type agentDocEntry struct {
	Path string
	Doc  string // "" 表示不适用/未提供
}

var registry []agentDocEntry

func registerAgentHelp(path, doc string) {
	registry = append(registry, agentDocEntry{Path: path, Doc: doc})
}

// NewRootCmd 构造根命令。
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:     "wechat-gzh-cli",
		Short:   "微信公众号服务端接口 CLI",
		Long:    "wechat-gzh-cli 封装微信公众号服务端接口,从环境变量读取鉴权凭证,面向 AI Agent 命令行调用。提供 --help(人)与 --agent-help(Agent)双界面。",
		Version: AppVersion,
		Args:    cobra.NoArgs,
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpRequested(c) {
				printRootAgentHelp()
				return nil
			}
			return c.Help()
		},
	}
	addAgentHelpFlag(root)
	setVersionTemplate(root)

	root.AddCommand(newAccessCmd())
	root.AddCommand(newServerIPCmd())
	root.AddCommand(newMaterialCmd())
	root.AddCommand(newMessageCmd())
	root.AddCommand(newUserCmd())
	root.AddCommand(newTagCmd())
	root.AddCommand(newMenuCmd())
	root.AddCommand(newCustomerCmd())
	root.AddCommand(newOpenapiCmd())
	root.AddCommand(newDraftCmd())
	root.AddCommand(newPublishCmd())
	root.AddCommand(newSubscribeCmd())
	root.AddCommand(newAutoreplyCmd())
	root.AddCommand(newStableTokenCmd())
	root.AddCommand(newCallbackCheckCmd())
	root.AddCommand(newPushServerIPCmd())
	root.AddCommand(newDatacubeCmd())
	root.AddCommand(newCommentCmd())
	root.AddCommand(newJSSDKTicketCmd())
	root.AddCommand(newHospitalMsgCmd())

	// 注册根命令覆盖清单(契约)
	registerAgentHelp("wechat-gzh-cli", "")

	return root
}

func setVersionTemplate(c *cobra.Command) {
	c.SetVersionTemplate(fmt.Sprintf("wechat-gzh-cli %s\n", AppVersion))
}

// printRootAgentHelp 输出根命令的 Agent 界面总览(所有命令的 --agent-help 覆盖状态)。
func printRootAgentHelp() {
	fmt.Println("# wechat-gzh-cli — Agent Help Overview")
	fmt.Println("> Root overview of every registered command's --agent-help coverage.")
	for _, e := range registry {
		status := "✅ documented"
		if e.Doc == "" {
			status = "⚠️ not-documented / n/a"
		}
		fmt.Printf("- %-40s %s\n", e.Path, status)
	}
	fmt.Println(`
---
Environment variables
  WX_GZ_AppID      公众号 AppID
  WX_GZ_AppSecret  公众号 AppSecret
Call any documented command with --agent-help to get its 11-dimension guide.`)
}

// Execute 是程序入口。
func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
