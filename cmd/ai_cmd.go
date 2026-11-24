package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"tx-gongzhonghao-cli/internal/client"
)

// newAICmd 智能接口命令组。
func newAICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ai",
		Short: "智能接口(AI 能力)",
		Long:  "微信翻译、语音转写、OCR 识别与图像处理。",
	}
	cmd.AddCommand(newAITranslateCmd())
	cmd.AddCommand(newAIVoiceToTextCmd())
	cmd.AddCommand(newAIVoiceResultCmd())
	for _, ep := range client.OcrEndpoints {
		cmd.AddCommand(newAIOcrSubCmd(ep))
	}
	return cmd
}

func newAITranslateCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "translate --content 文本 [--from zh_CN --to en_US]",
		Short: "微信翻译",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			content, _ := c.Flags().GetString("content")
			lfrom, _ := c.Flags().GetString("from")
			lto, _ := c.Flags().GetString("to")
			if content == "" {
				return fmt.Errorf("--content 必填")
			}
			if lfrom == "" {
				lfrom = "zh_CN"
			}
			if lto == "" {
				lto = "en_US"
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				translated, err := cl.TranslateContent(ctx, lfrom, lto, content)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(map[string]any{"lfrom": lfrom, "lto": lto, "tocontent": translated})
				}
				fmt.Fprintln(c.OutOrStdout(), translated)
				return nil
			})
		},
	}
	cmd.Flags().String("content", "", "待翻译文本")
	cmd.Flags().String("from", "zh_CN", "源语言")
	cmd.Flags().String("to", "en_US", "目标语言")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli ai translate", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli ai translate --content 你好 --from zh_CN --to en_US",
		Summary:        "微信翻译:把文本从源语言翻译为目标语言。",
		PreConditions:  []string{"已配置鉴权。", "--content 必填。"},
		PostConditions: []string{"stdout 输出翻译结果。"},
		Idempotency:    "YES. 只读。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限。",
		RateLimit:      "受智能接口配额限制。",
		RaceConditions: []string{"无。"},
		Orchestration:  "多语言内容运营。",
		ParsingHints:   "文本输出 tocontent;--json 输出 {\"tocontent\":..}。",
		Examples:       []string{"wechat-gzh-cli ai translate --content 你好 --to en_US"},
	}))
	return cmd
}

func newAIVoiceToTextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "voice-to-text --format mp3 --voice-id xxx --lang zh_CN --file ./a.mp3",
		Short: "上传语音转文字",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			format, _ := c.Flags().GetString("format")
			voiceID, _ := c.Flags().GetString("voice-id")
			lang, _ := c.Flags().GetString("lang")
			file, _ := c.Flags().GetString("file")
			if format == "" || voiceID == "" || file == "" {
				return fmt.Errorf("--format/--voice-id/--file 必填")
			}
			if lang == "" {
				lang = "zh_CN"
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				if err := cl.VoiceToText(ctx, format, voiceID, lang, file); err != nil {
					return err
				}
				fmt.Fprintln(c.OutOrStdout(), "ok")
				return nil
			})
		},
	}
	cmd.Flags().String("format", "mp3", "语音格式: mp3/...")
	cmd.Flags().String("voice-id", "", "语音 id,用于后续查询结果")
	cmd.Flags().String("lang", "zh_CN", "语言")
	cmd.Flags().String("file", "", "本地语音文件路径")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli ai voice-to-text", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli ai voice-to-text --format mp3 --voice-id v1 --file ./a.mp3",
		Summary:        "上传语音文件供转文字识别,返回 voice_id 用于查询。",
		PreConditions:  []string{"已配置鉴权。", "--format/--voice-id/--file 必填。"},
		PostConditions: []string{"stdout 输出 ok;语音已上传待识别。"},
		Idempotency:    "NO.",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"上传语音文件。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限。",
		RateLimit:      "受智能接口配额限制。",
		RaceConditions: []string{"识别结果异步,需用 voice-result 查询。"},
		Orchestration:  "上传后用 ai voice-result --voice-id 查询识别结果。",
		ParsingHints:   "文本输出 ok。",
		Examples:       []string{"wechat-gzh-cli ai voice-to-text --voice-id v1 --file ./a.mp3"},
	}))
	return cmd
}

func newAIVoiceResultCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "voice-result --voice-id xxx [--lang zh_CN]",
		Short: "查询语音识别结果",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			voiceID, _ := c.Flags().GetString("voice-id")
			lang, _ := c.Flags().GetString("lang")
			if voiceID == "" {
				return fmt.Errorf("--voice-id 必填")
			}
			if lang == "" {
				lang = "zh_CN"
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				res, err := cl.VoiceResult(ctx, voiceID, lang)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(res)
				}
				fmt.Fprintln(c.OutOrStdout(), res)
				return nil
			})
		},
	}
	cmd.Flags().String("voice-id", "", "语音 id")
	cmd.Flags().String("lang", "zh_CN", "语言")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli ai voice-result", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli ai voice-result --voice-id v1",
		Summary:        "查询此前上传语音的转文字识别结果。",
		PreConditions:  []string{"已配置鉴权。", "--voice-id 必填。"},
		PostConditions: []string{"stdout 输出识别结果。"},
		Idempotency:    "YES. 只读。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{},
		Scope:          "开发者权限。",
		RateLimit:      "受智能接口配额限制。",
		RaceConditions: []string{"结果可能要等识别完成。"},
		Orchestration:  "在 voice-to-text 之后轮询本命令。",
		ParsingHints:   "--json 输出完整识别结果。",
		Examples:       []string{"wechat-gzh-cli ai voice-result --voice-id v1 --json"},
	}))
	return cmd
}

// newAIOcrSubCmd 为每个 OCR/图像处理接口生成子命令。
func newAIOcrSubCmd(ep client.OcrEndpoint) *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   ep.Name + " --img-url HTTPS_URL",
		Short: ep.Summary,
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			imgURL, _ := c.Flags().GetString("img-url")
			if imgURL == "" {
				return fmt.Errorf("--img-url 必填(图片 URL,需编码)")
			}
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				res, err := cl.OcrQuery(ctx, ep.Path, imgURL)
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
	cmd.Flags().String("img-url", "", "图片的 URL 地址")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli ai "+ep.Name, RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli ai " + ep.Name + " --img-url https://.../img.jpg",
		Summary:        ep.Summary + "(基于图片 URL)。",
		PreConditions:  []string{"已配置鉴权。", "--img-url 必填且为可访问的图片 URL。"},
		PostConditions: []string{"stdout 输出识别/处理结果。"},
		Idempotency:    "YES. 只读。",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{{Code: "40002", Desc: "img_url 无效", Fix: "检查 URL。"}},
		Scope:          "开发者权限。",
		RateLimit:      "受智能接口配额限制。",
		RaceConditions: []string{"无。"},
		Orchestration:  "图片来自素材或外链。",
		ParsingHints:   "--json 输出识别字段。",
		Examples:       []string{"wechat-gzh-cli ai " + ep.Name + " --img-url https://example.com/a.jpg --json"},
	}))
	return cmd
}
