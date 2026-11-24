package tests

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"tx-gongzhonghao-cli/cmd"
)

// collectLeafPaths 深度优先收集所有叶子命令路径(无子命令的命令)。
func collectLeafPaths(c *cobra.Command, out map[string]bool) {
	if len(c.Commands()) == 0 {
		out[c.CommandPath()] = true
		return
	}
	for _, sub := range c.Commands() {
		collectLeafPaths(sub, out)
	}
}

// TestAgentHelpReachable 遍历全部叶子命令,验证 --agent-help 在**未配置环境变量**
// 且不给任何业务参数时仍可输出(满足三条可达性法则之二:文档先于鉴权可见)。
func TestAgentHelpReachable(t *testing.T) {
	// 确保环境变量未设置,以验证"无鉴权也能看文档"。
	t.Setenv("WX_GZ_AppID", "")
	t.Setenv("WX_GZ_AppSecret", "")

	leaves := map[string]bool{}
	collectLeafPaths(cmd.NewRootCmd(), leaves)
	if len(leaves) == 0 {
		t.Fatal("no leaf commands found")
	}

	for pathStr := range leaves {
		pathStr := pathStr
		t.Run(pathStr, func(t *testing.T) {
			// root 命令本身名为 wechat-gzh-cli,SetArgs 需传入去除根命令名后的子命令参数。
			fields := strings.Fields(pathStr)
			if len(fields) == 0 {
				t.Fatal("empty path")
			}
			subArgs := fields[1:] // 去掉 "wechat-gzh-cli"
			buf := &bytes.Buffer{}
			root := cmd.NewRootCmd()
			root.SetOut(buf)
			root.SetErr(buf)
			argv := append(subArgs, "--agent-help")
			root.SetArgs(argv)
			if err := root.Execute(); err != nil {
				t.Fatalf("execute %v: %v", argv, err)
			}
			if !strings.Contains(buf.String(), "Agent Help") {
				t.Fatalf("--agent-help 未见文档(可能被鉴权/参数校验拦截) out=%q", buf.String())
			}
		})
	}
}
