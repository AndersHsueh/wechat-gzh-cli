package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"tx-gongzhonghao-cli/internal/client"
)

func newMaterialCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "material",
		Short: "素材管理(永久/临时/正文图)",
		Long:  "素材管理:上传/获取/列表/删除永久素材,以及正文图与临时素材上传。",
	}
	cmd.AddCommand(newMaterialUploadCmd())
	cmd.AddCommand(newMaterialListCmd())
	cmd.AddCommand(newMaterialGetCmd())
	cmd.AddCommand(newMaterialDeleteCmd())
	cmd.AddCommand(newMaterialUploadTempCmd())
	cmd.AddCommand(newMaterialUploadImgCmd())
	return cmd
}

func newMaterialListCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "list --type image [--offset 0 --count 20]",
		Short: "分类型获取永久素材列表",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			mediaType, _ := c.Flags().GetString("type")
			offset, _ := c.Flags().GetInt("offset")
			count, _ := c.Flags().GetInt("count")
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				res, err := cl.BatchGetMaterial(ctx, mediaType, offset, count)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(res)
				}
				for _, it := range res.Item {
					fmt.Fprintf(c.OutOrStdout(), "%s\t%s\t%d\n", it.MediaID, it.Name, it.UpdateTime)
				}
				return nil
			})
		},
	}
	cmd.Flags().String("type", "image", "素材类型: image/video/voice/news")
	cmd.Flags().Int("offset", 0, "偏移位置")
	cmd.Flags().Int("count", 20, "数量(1-20)")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli material list", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli material list --type image",
		Summary:        "分类型分页获取公众号永久素材列表。",
		PreConditions:  []string{"已配置 WX_GZ_AppID/WX_GZ_AppSecret。", "--type 必填(image/video/voice/news)。"},
		PostConditions: []string{"stdout 输出素材列表,每行为一个 media_id + name + update_time,或 --json 输出完整结构。"},
		Idempotency:    "YES. 只读查询。",
		Atomicity:      "Atomic. 单次 POST。",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{{Code: "40004", Desc: "invalid media type", Fix: "检查 --type 枚举。"}, {Code: "40007", Desc: "invalid media_id", Fix: "internal"}},
		Scope:          "开发者权限。",
		RateLimit:      "受素材接口调用额度限制。",
		RaceConditions: []string{"列表是分页快照,翻页时素材可能变化。"},
		Orchestration:  "拿到 media_id 后可进一步用 material get 查看详情或用 material delete 删除。",
		ParsingHints:   "--json 输出 {\"total_count\":..,\"item\":[{\"media_id\":..}]}。",
		Examples:       []string{"wechat-gzh-cli material list --type image --json"},
	}))
	return cmd
}

func newMaterialUploadCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "upload --type image --file ./img.jpg",
		Short: "上传永久素材",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			mediaType, _ := c.Flags().GetString("type")
			file, _ := c.Flags().GetString("file")
			if file == "" {
				return fmt.Errorf("--file 必填")
			}
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				res, err := cl.AddMaterial(ctx, mediaType, file, f)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(res)
				}
				fmt.Fprintf(c.OutOrStdout(), "%s\n", res.MediaID)
				return nil
			})
		},
	}
	cmd.Flags().String("type", "image", "媒体类型: image/voice/video/thumb")
	cmd.Flags().String("file", "", "本地文件路径")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli material upload", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli material upload --type image --file ./img.jpg",
		Summary:        "上传永久素材(图片/语音/视频/缩略图),返回 media_id。",
		PreConditions:  []string{"已配置鉴权。", "--file 指向存在的文件。"},
		PostConditions: []string{"stdout 输出 media_id;素材持久保存在素材库。"},
		Idempotency:    "NO. 每次调用都会新建一个素材条目(不同 media_id)。",
		Atomicity:      "Atomic. 单次 multipart POST。",
		SideEffects:    []string{"公众号素材库新增一条素材。"},
		CommonErrors:   []ErrSpec{{Code: "40005", Desc: "invalid file type", Fix: "检查文件格式/大小。"}},
		Scope:          "开发者权限。",
		RateLimit:      "受素材接口调用额度限制。",
		RaceConditions: []string{"并发上传同名文件会各自生成 media_id,互不影响。"},
		Orchestration:  "上传成功后,把 media_id 用于群发、菜单或草稿封面。",
		ParsingHints:   "--json 输出 {\"media_id\":\"...\"}。",
		Examples:       []string{`wechat-gzh-cli material upload --type image --file ./a.jpg --json`},
	}))
	return cmd
}

func newMaterialGetCmd() *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "get --media-id MEDIA_ID",
		Short: "获取永久素材详情",
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
				items, err := cl.GetMaterial(ctx, mediaID)
				if err != nil {
					return err
				}
				if jsonOut {
					return writeJSON(items)
				}
				for _, it := range items {
					fmt.Fprintf(c.OutOrStdout(), "%s\t%s\n", it.Title, it.URL)
				}
				return nil
			})
		},
	}
	cmd.Flags().String("media-id", "", "素材 media_id")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "以 JSON 输出")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli material get", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli material get --media-id MEDIA_ID",
		Summary:        "按 media_id 获取永久素材详细信息(图文素材返回 news_item)。",
		PreConditions:  []string{"已配置鉴权。", "--media-id 必填。"},
		PostConditions: []string{"stdout 输出素材内容(标题/URL 等)。"},
		Idempotency:    "YES. 只读。",
		Atomicity:      "Atomic. 单次 POST。",
		SideEffects:    []string{"无。"},
		CommonErrors:   []ErrSpec{{Code: "40007", Desc: "invalid media_id", Fix: "检查 media_id 是否正确。"}},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"无。"},
		Orchestration:  "用于查看素材正文/封面后再决定是否群发。",
		ParsingHints:   "--json 输出 news_item 数组。",
		Examples:       []string{"wechat-gzh-cli material get --media-id xxxx --json"},
	}))
	return cmd
}

func newMaterialDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete --media-id MEDIA_ID",
		Short: "删除永久素材",
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
				if err := cl.DelMaterial(ctx, mediaID); err != nil {
					return err
				}
				fmt.Fprintln(c.OutOrStdout(), "deleted")
				return nil
			})
		},
	}
	cmd.Flags().String("media-id", "", "素材 media_id")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli material delete", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli material delete --media-id MEDIA_ID",
		Summary:        "删除指定的永久素材。",
		PreConditions:  []string{"已配置鉴权。", "--media-id 必填。"},
		PostConditions: []string{"stdout 输出 deleted;素材从素材库移除。"},
		Idempotency:    "NO. 删除有副作用;重复删除同一 media_id 会报错。",
		Atomicity:      "Atomic. 单次 POST。",
		SideEffects:    []string{"公众号素材库中的对应素材被永久删除。"},
		CommonErrors:   []ErrSpec{{Code: "40007", Desc: "invalid media_id", Fix: "检查 media_id。"}},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"若素材正被群发引用,删除可能导致引用失效。"},
		Orchestration:  "删除前建议先用 material list 确认 media_id。",
		ParsingHints:   "文本输出 deleted 表示成功。",
		Examples:       []string{"wechat-gzh-cli material delete --media-id xxxx"},
	}))
	return cmd
}

func newMaterialUploadTempCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload-temp --type image --file ./a.jpg",
		Short: "新增临时素材",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			mediaType, _ := c.Flags().GetString("type")
			file, _ := c.Flags().GetString("file")
			if file == "" {
				return fmt.Errorf("--file 必填")
			}
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				res, err := cl.UploadTempMedia(ctx, mediaType, file, f)
				if err != nil {
					return err
				}
				fmt.Fprintf(c.OutOrStdout(), "%s\n", res.MediaID)
				return nil
			})
		},
	}
	cmd.Flags().String("type", "image", "媒体类型: image/voice/video/thumb")
	cmd.Flags().String("file", "", "本地文件路径")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli material upload-temp", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli material upload-temp --type image --file ./a.jpg",
		Summary:        "上传临时素材,media_id 3 天内有效。",
		PreConditions:  []string{"已配置鉴权。", "--file 存在。"},
		PostConditions: []string{"stdout 输出临时 media_id;3 天后失效。"},
		Idempotency:    "NO.",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"微信后台临时保存素材(3 天)。"},
		CommonErrors:   []ErrSpec{{Code: "40004", Desc: "invalid media type", Fix: "检查 type。"}},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"临时素材 3 天后失效需提前取用。"},
		Orchestration:  "临时 media_id 常用于群发前的资源准备。",
		ParsingHints:   "文本输出 media_id。",
		Examples:       []string{"wechat-gzh-cli material upload-temp --type image --file ./a.jpg"},
	}))
	return cmd
}

func newMaterialUploadImgCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload-img --file ./img.jpg",
		Short: "上传发表内容中的图片(正文图)",
		Args:  AgentHelpAwareArgs(cobra.NoArgs, nil),
		RunE: func(c *cobra.Command, _ []string) error {
			if agentHelpShortcut(c) {
				return nil
			}
			file, _ := c.Flags().GetString("file")
			if file == "" {
				return fmt.Errorf("--file 必填")
			}
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()
			return runWithClient(c, func(ctx context.Context, cl *client.Client) error {
				u, err := cl.UploadImage(ctx, file, f)
				if err != nil {
					return err
				}
				fmt.Fprintln(c.OutOrStdout(), u)
				return nil
			})
		},
	}
	cmd.Flags().String("file", "", "本地图片路径(jpg/png,<1MB)")
	addAgentHelpFlag(cmd)
	registerAgentHelp("wechat-gzh-cli material upload-img", RenderAgentHelp(AgentHelpDoc{
		Command:        "wechat-gzh-cli material upload-img --file ./img.jpg",
		Summary:        "上传发表内容(文章/贴图)所需的图片,返回可直接用于图文正文的 URL。不占用素材库 10 万上限。",
		PreConditions:  []string{"已配置鉴权。", "--file 存在, jpg/png, <1MB。"},
		PostConditions: []string{"stdout 输出图片 URL,可插入图文 content。"},
		Idempotency:    "NO.",
		Atomicity:      "Atomic.",
		SideEffects:    []string{"生成一个图片 URL(不计入素材配额)。"},
		CommonErrors:   []ErrSpec{{Code: "40005", Desc: "invalid file type", Fix: "检查格式。"}, {Code: "40009", Desc: "invalid image size", Fix: "压缩图片到 1MB 内。"}},
		Scope:          "开发者权限。",
		RateLimit:      "常规。",
		RaceConditions: []string{"URL 长期有效,可复用。"},
		Orchestration:  "上传后把 URL(base) 用于草稿/群发图文正文。",
		ParsingHints:   "文本输出 URL。",
		Examples:       []string{"wechat-gzh-cli material upload-img --file ./a.jpg"},
	}))
	return cmd
}
