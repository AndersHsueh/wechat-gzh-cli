// Package cmd 定义 wechat-gzh-cli 的所有命令(基于 Cobra)。
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// ErrSpec 描述一个常见错误及其处理建议,用于 "Common errors" 章节。
type ErrSpec struct {
	Code string // 错误码或 HTTP 状态,如 "40001"、"401"
	Desc string // 错误说明
	Fix  string // Agent 应如何判断/恢复
}

// AgentHelpDoc 是 --agent-help 输出的结构化内容(遵循双界面 CLI 设计的 11 维度模板)。
type AgentHelpDoc struct {
	Command         string   // 命令路径,如 "wechat-gzh-cli material list"
	Summary         string   // 一句话功能
	PreConditions   []string // 1. 前置条件
	PostConditions  []string // 2. 后置条件
	Idempotency     string   // 3. 幂等性
	Atomicity       string   // 4. 原子性
	SideEffects     []string // 5. 副作用
	CommonErrors    []ErrSpec // 6. 常见错误与恢复
	Scope           string   // 7. 权限范围(OAuth scope)
	RateLimit       string   // 8. 频率限制
	RaceConditions  []string // 9. 竞态条件
	Orchestration   string   // 10. 编排提示(下一步做什么)
	ParsingHints    string   // 11. 解析提示(如何解析输出,含 jq 示例)
	Examples        []string // 可直接复制运行的示例
}

func bullet(items []string) string {
	if len(items) == 0 {
		return "- (none)"
	}
	var sb strings.Builder
	for _, it := range items {
		sb.WriteString("- " + it + "\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}

// RenderAgentHelp 按固定模板渲染 --agent-help 文档,保证结构可被程序稳定解析。
func RenderAgentHelp(doc AgentHelpDoc) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("# %s — Agent Help\n\n", doc.Command))
	b.WriteString("> A structured help document for AI agents (dual-surface CLI design).\n\n")
	b.WriteString("## Summary\n" + doc.Summary + "\n\n")

	b.WriteString("## Pre-conditions\n" + bullet(doc.PreConditions) + "\n\n")
	b.WriteString("## Post-conditions (on success)\n" + bullet(doc.PostConditions) + "\n\n")
	b.WriteString("## Idempotency\n" + doc.Idempotency + "\n\n")
	b.WriteString("## Atomicity\n" + doc.Atomicity + "\n\n")
	b.WriteString("## Side effects\n" + bullet(doc.SideEffects) + "\n\n")

	b.WriteString("## Common errors\n")
	if len(doc.CommonErrors) == 0 {
		b.WriteString("- (none)\n")
	} else {
		for _, e := range doc.CommonErrors {
			b.WriteString(fmt.Sprintf("- **%s**: %s. Fix: %s\n", e.Code, e.Desc, e.Fix))
		}
	}
	b.WriteString("\n")

	b.WriteString("## Scope (OAuth permission)\n" + doc.Scope + "\n\n")
	b.WriteString("## Rate limit\n" + doc.RateLimit + "\n\n")
	b.WriteString("## Race conditions\n" + bullet(doc.RaceConditions) + "\n\n")
	b.WriteString("## Orchestration hints\n" + doc.Orchestration + "\n\n")
	b.WriteString("## Parsing hints\n" + doc.ParsingHints + "\n\n")

	if len(doc.Examples) > 0 {
		b.WriteString("## Examples\n")
		for _, ex := range doc.Examples {
			b.WriteString("```\n" + ex + "\n```\n")
		}
	}
	return b.String()
}

// addAgentHelpFlag 为命令添加 --agent-help 布尔开关。
func addAgentHelpFlag(c *cobra.Command) {
	c.Flags().Bool("agent-help", false, "面向 AI Agent 的结构化帮助文档(双界面 CLI 设计)")
}

// agentHelpRequested 判断当前命令是否请求 --agent-help(满足三条可达性法则:不触发真实执行)。
func agentHelpRequested(c *cobra.Command) bool {
	if c == nil {
		return false
	}
	if c.Flags().Changed("agent-help") {
		v, err := c.Flags().GetBool("agent-help")
		return err == nil && v
	}
	return false
}

// agentHelpShortcut 在命令 RunE 最开头调用:若请求了 --agent-help,打印注册的文档并返回 true。
// 因 Args 已被 AgentHelpAwareArgs 跳过,故 --agent-help 可在任何业务校验之前可见(法则一),
// 且只打印、不触达 client(法则二/三)。
func agentHelpShortcut(c *cobra.Command) bool {
	if !agentHelpRequested(c) {
		return false
	}
	for _, e := range registry {
		if e.Path == c.CommandPath() {
			fmt.Fprint(c.OutOrStdout(), e.Doc)
		}
	}
	return true
}

// AgentHelpAwareArgs 按 double-surface 设计包装位置参数校验:
// 当 --agent-help 在场时跳过参数校验(法则一:文档必须先于参数校验可见)。
func AgentHelpAwareArgs(wrapped cobra.PositionalArgs, c *cobra.Command) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if agentHelpRequested(c) {
			return nil
		}
		return wrapped(cmd, args)
	}
}
