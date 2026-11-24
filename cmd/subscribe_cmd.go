package cmd

import (
	"github.com/spf13/cobra"
)

func newSubscribeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscribe",
		Short: "一次性订阅消息",
		Long:  "向授权用户发送一次性订阅模板消息。",
	}
	cmd.AddCommand(newSubscribeSendCmd())
	return cmd
}
