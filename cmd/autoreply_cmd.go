package cmd

import (
	"github.com/spf13/cobra"
)

func newAutoreplyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "autoreply",
		Short: "自动回复规则",
		Long:  "获取公众号自动回复规则(关注/消息/关键词)。",
	}
	cmd.AddCommand(newAutoreplyGetCmd())
	return cmd
}
