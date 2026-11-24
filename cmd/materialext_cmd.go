package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"tx-gongzhonghao-cli/internal/client"
)

func newMaterialCountCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "count",
		Short: "获取永久素材总数",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				mc, err := cl.GetMaterialCount(ctx)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(mc)
				}
				fmt.Fprintf(c.OutOrStdout(), "img=%d voice=%d video=%d news=%d\n", mc.ImageCount, mc.VoiceCount, mc.VideoCount, mc.NewsCount)
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli material count", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli material count",
		Summary:        "获取永久素材总数(图片/语音/视频/图文)。",
		PreConditions:  []string{"已配置鉴权。"},
		PostConditions: []string{"stdout 输出各类型数量。"},
		Idempotency:    "YES. 只读。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"无。"},
		Orchestration:  "用于监控素材配额使用。",
		ParsingHints:   "--json 输出 {\"image_count\":..,..}。",
		Examples:       []string{"wechat-gzh-cli material count --json"},
	}))
	return cmd
}

func newMaterialGetTempCmd() *cobra.Command {
	var out string
	cmd := &cobra.Command{
		Use:   "get-temp --media-id MEDIA_ID [--out 文件]",
		Short: "获取临时素材",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			mediaID, _ := c.Flags().GetString("media-id")
			if mediaID == "" {
				return fmt.Errorf("--media-id 必填")
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				data, err := cl.GetTempMedia(ctx, mediaID)
				if err != nil {
					return err
				}
				if out != "" {
					if err := os.WriteFile(out, data, 0o644); err != nil {
						return err
					}
					fmt.Fprintf(c.OutOrStdout(), "saved -> %s (%d bytes)\n", out, len(data))
					return nil
				}
				_, err = fmt.Fprintln(c.OutOrStdout(), "bytes=", len(data))
				return err
			})
		},
	}
	cmd.Flags().String("media-id", "", "临时素材 media_id")
	cmd.Flags().StringVar(&out, "out", "", "保存到本地文件路径(否则仅打印大小)")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli material get-temp", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli material get-temp --media-id MEDIA_ID --out ./a.jpg",
		Summary:        "下载临时素材文件(需 3 天内)。",
		PreConditions:  []string{"已配置鉴权。", "--media-id 必填且未过期。"},
		PostConditions: []string{"stdout 输出保存路径/字节数。"},
		Idempotency:    "YES. 只读下载。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"下载文件到本地。"},
		CommonErrors:   []ErrSpec{{Code: "40007", Desc: "invalid media_id", Fix: "media_id 过期或错误。"}},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"临时素材 3 天失效。"},
		Orchestration:  "在 3 天有效期内取用临时素材。",
		ParsingHints:   "文本输出保存信息。",
		Examples:       []string{"wechat-gzh-cli material get-temp --media-id m_xxx --out ./a.jpg"},
	}))
	return cmd
}

func newMaterialGetJSSDKVoiceCmd() *cobra.Command {
	var out string
	cmd := &cobra.Command{
		Use:   "get-jssdk-voice --media-id MEDIA_ID [--out 文件]",
		Short: "获取高清语音素材",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			mediaID, _ := c.Flags().GetString("media-id")
			if mediaID == "" {
				return fmt.Errorf("--media-id 必填")
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				data, err := cl.GetJSSDKRawVoice(ctx, mediaID)
				if err != nil {
					return err
				}
				if out != "" {
					if err := os.WriteFile(out, data, 0o644); err != nil {
						return err
					}
					fmt.Fprintf(c.OutOrStdout(), "saved -> %s (%d bytes)\n", out, len(data))
					return nil
				}
				_, err = fmt.Fprintln(c.OutOrStdout(), "bytes=", len(data))
				return err
			})
		},
	}
	cmd.Flags().String("media-id", "", "高清语音素材 media_id")
	cmd.Flags().StringVar(&out, "out", "", "保存到本地文件路径(否则仅打印大小)")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli material get-jssdk-voice", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli material get-jssdk-voice --media-id m_xxx --out ./a.mp3",
		Summary:        "获取高清语音素材(高质量音源)。",
		PreConditions:  []string{"已配置鉴权。", "--media-id 必填。"},
		PostConditions: []string{"stdout 输出保存路径/字节数。"},
		Idempotency:    "YES. 只读下载。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"下载文件到本地。"},
		CommonErrors:   []ErrSpec{{Code: "40007", Desc: "invalid media_id", Fix: "media_id 错误。"}},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"无。"},
		Orchestration:  "用于需要高质量音频的场景。",
		ParsingHints:   "文本输出保存信息。",
		Examples:       []string{"wechat-gzh-cli material get-jssdk-voice --media-id m_xxx --out ./a.mp3"},
	}))
	return cmd
}
