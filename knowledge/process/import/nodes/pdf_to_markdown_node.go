package processimportnodes

import (
	"context"
	chatmodel "doc-precise-rag/knowledge/chat-model"
	"doc-precise-rag/knowledge/clients"
	"doc-precise-rag/knowledge/config"
	"doc-precise-rag/knowledge/utils"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/lithammer/dedent"
	"github.com/tencentyun/cos-go-sdk-v5"
	"golang.org/x/sync/errgroup"
)

type PdfToMarkdownNode struct {
	logger *slog.Logger
}

func (node *PdfToMarkdownNode) Process(ctx context.Context, docLink *DocLink) (*MarkdownData, error) {
	node.logger.Info("开始解析PDF了")
	tempUrl, err := clients.CosClient.Object.GetPresignedURL2(context.Background(), http.MethodGet, docLink.ObjKey, 2*time.Hour, nil)
	if err != nil {
		node.logger.Error("获取临时授权URL失败", "error", err)
		return nil, err
	}
	task, err := clients.MineruClient.CreateParserFileTask(tempUrl.String())
	if err != nil {
		node.logger.Error("创建解析PDF任务失败", "error", err)
		return nil, err
	}

	node.logger.Info("创建解析PDF任务成功，开始获取任务状态,每2秒获取一次", "taskID", *task.Data.TaskID)

	var fullZipURL string

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(2 * time.Second):
		}

		statusResp, err := clients.MineruClient.GetParserFileStatus(*task.Data.TaskID)
		if err != nil {
			node.logger.Error("获取解析PDF任务状态失败", "error", err)
			return nil, err
		}

		if statusResp == nil || statusResp.Data == nil || statusResp.Data.State == nil {
			node.logger.Error("获取解析PDF任务状态失败", "error", errors.New("任务状态为空"))
			return nil, errors.New("任务状态为空")
		}
		if *statusResp.Data.State == "done" {
			if statusResp.Data.FullZipURL == nil {
				node.logger.Error("任务解析成功但下载链接为空", "error", errors.New("fullZipURL为空"))
				return nil, errors.New("任务解析成功但下载链接为空")
			}
			fullZipURL = *statusResp.Data.FullZipURL
			break
		}
		if *statusResp.Data.State == "pending" {
			node.logger.Info("解析PDF任务状态为排队中，继续等待")
			continue
		}

		if *statusResp.Data.State == "running" {
			node.logger.Info("解析PDF任务状态为解析中,继续等待")
			continue
		}

		if *statusResp.Data.State == "converting" {
			node.logger.Info("解析PDF任务状态为转换中,继续等待")
			continue
		}

		if *statusResp.Data.State == "failed" {
			node.logger.Error("解析PDF任务状态为解析失败", "error", *statusResp.Data.ErrorMsg)
			if statusResp.Data.ErrorMsg == nil {
				return nil, errors.New("解析PDF任务状态为解析失败，但错误信息为空")
			}
			return nil, errors.New(*statusResp.Data.ErrorMsg)
		}

		return nil, errors.New("解析PDF任务状态未知：" + *statusResp.Data.State)

	}

	if fullZipURL == "" {
		node.logger.Error("解析PDF任务状态为解析成功，但下载链接为空", "error", errors.New("fullZipURL为空"))
		return nil, errors.New("解析PDF任务状态为解析成功，但下载链接为空")
	}

	node.logger.Info("解析成功，开始下载解析后内容", "fullZipURL", fullZipURL)

	unzipDir, err := node.downloadAndUnzip(fullZipURL, *task.Data.TaskID)
	if err != nil {
		node.logger.Error("下载解析后内容失败", "error", err)
		return nil, err
	}

	node.logger.Info("下载解析后内容成功,开始处理本地图片", "unzipDir", unzipDir)
	mdObjectKey, newMdContent, err := node.dealWithMdImages(ctx, unzipDir, docLink.DocTitle)
	if err != nil {
		node.logger.Error("处理本地图片失败", "error", err)
		return nil, err
	}

	return &MarkdownData{
		Content:  newMdContent,
		ObjKey:   mdObjectKey,
		DocTitle: docLink.DocTitle,
	}, nil
}

func (node *PdfToMarkdownNode) BuildInvokableLambda() *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, input *DocLink) (*MarkdownData, error) {
		return node.Process(ctx, input)
	})
}

func NewPdfToMarkdownNode() *PdfToMarkdownNode {
	return &PdfToMarkdownNode{logger: slog.With("流程类型", "文档导入", "节点名称", "pdfToMarkdownNode")}
}

func (node *PdfToMarkdownNode) downloadAndUnzip(fullZipURL string, taskID string) (string, error) {

	baseDir, err := filepath.Abs(
		filepath.Join(*config.AppConfig.ImportProcessConfig.ImportFileTempDir,
			time.Now().Format("2006-01-02")))
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(baseDir, 0755)
	if err != nil {
		return "", err
	}

	zipPath := filepath.Join(baseDir, taskID+".zip")
	err = clients.Request.GetDownload(fullZipURL, zipPath, nil)
	if err != nil {
		return "", err
	}

	unzipDir := filepath.Join(baseDir, taskID)

	err = utils.Unzip(zipPath, unzipDir)

	if err != nil {
		return "", err
	}

	return unzipDir, nil

}

func (node *PdfToMarkdownNode) dealWithMdImages(baseCtx context.Context, markdownDir string, docTitle string) (string, string, error) {
	imageNames := node.getMdImages(markdownDir)
	mdContent, err := node.getMdContent(markdownDir)
	if err != nil {
		return "", "", err
	}

	lines := strings.Split(mdContent, "\n")

	var imageContextMapMu sync.Mutex

	imageContextMap := map[string][]ImageContext{}

	imageContextMapGroup, ctx := errgroup.WithContext(baseCtx)

	node.logger.Info("开始获取本地图片对应上下文")
	for _, imageName := range imageNames {
		imageContextMapGroup.Go(func() error {
			imageContexts, err := node.getImageContext(ctx, imageName, lines)
			if err != nil {
				return err
			}
			imageContextMapMu.Lock()
			imageContextMap[imageName] = imageContexts
			imageContextMapMu.Unlock()
			return nil
		})

	}

	err = imageContextMapGroup.Wait()
	if err != nil {
		return "", "", err
	}

	imageDir := filepath.Join(markdownDir, "images")

	var imageDescriptionMapMu sync.Mutex
	imageDescriptionMap := map[string][]string{}
	imageDescriptionMapGroup, ctx := errgroup.WithContext(baseCtx)
	node.logger.Info("开始生成图片对应描述")
	for imageName, imageContext := range imageContextMap {
		imageDescriptionMapGroup.Go(func() error {
			if len(imageContext) == 0 {
				return nil
			}
			descriptions, err := node.getImageDescription(ctx, docTitle, imageContext, imageName, imageDir)
			if err != nil {
				return err
			}
			imageDescriptionMapMu.Lock()
			imageDescriptionMap[imageName] = descriptions
			imageDescriptionMapMu.Unlock()
			return nil
		})
	}

	err = imageDescriptionMapGroup.Wait()
	if err != nil {
		return "", "", err
	}

	imageOnlineLinkMap := map[string]string{}

	var imageOnlineLinkMapMu sync.Mutex
	imageOnlineLinkMapGroup, ctx := errgroup.WithContext(baseCtx)
	node.logger.Info("开始上传图片到COS")
	for imageName := range imageDescriptionMap {
		imageOnlineLinkMapGroup.Go(func() error {
			onlineLink, err := node.getImageOnlineLink(ctx,
				imageName,
				filepath.Join(imageDir, imageName),
				mime.TypeByExtension(filepath.Ext(imageName)),
				docTitle)
			if err != nil {
				return err
			}
			imageOnlineLinkMapMu.Lock()
			imageOnlineLinkMap[imageName] = onlineLink
			imageOnlineLinkMapMu.Unlock()
			return nil
		})
	}

	err = imageOnlineLinkMapGroup.Wait()
	if err != nil {
		return "", "", err
	}

	node.logger.Info("开始替换图片为在线图片")
	for imageName, descriptions := range imageDescriptionMap {
		onlineLink := imageOnlineLinkMap[imageName]
		imageDescriptionRegex := regexp.MustCompile(`!\[.*?\]\(images/` + regexp.QuoteMeta(imageName) + `(?:\s+.*?)?\)`)
		for _, description := range descriptions {
			onlineImageTag := fmt.Sprintf(`![%s](%s)`, description, onlineLink)
			mdContent = utils.ReplaceStr(imageDescriptionRegex, mdContent, onlineImageTag, 1)
		}
	}

	markdownObjectKey := fmt.Sprintf("markdown/%s.md", docTitle)

	_, err = clients.CosClient.Object.Put(
		baseCtx,
		markdownObjectKey,
		strings.NewReader(mdContent),
		&cos.ObjectPutOptions{
			ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
				ContentType: "text/markdown",
			},
		},
	)
	if err != nil {
		node.logger.Error("cos 上传 markdown 文件失败", "error", err)
		return "", "", err
	}
	return markdownObjectKey, mdContent, nil
}

func (node *PdfToMarkdownNode) getMdImages(markdownDir string) []string {
	files, err := os.ReadDir(filepath.Join(markdownDir, "images"))
	if err != nil {
		node.logger.Debug("本地图片目录不存在或读取失败", "error", err)
		return []string{}
	}

	var imageNames []string

	for _, file := range files {
		imageNames = append(imageNames, file.Name())
	}

	return imageNames
}

func (node *PdfToMarkdownNode) getMdContent(markdownDir string) (string, error) {
	markdownPath := filepath.Join(markdownDir, "full.md")
	content, err := os.ReadFile(markdownPath)
	if err != nil {
		node.logger.Debug("markdown文件不存在或读取失败", "error", err)
		return "", err
	}
	return string(content), nil
}

func (node *PdfToMarkdownNode) getImageContext(ctx context.Context, imageName string, contentLines []string) ([]ImageContext, error) {

	re := regexp.MustCompile(`!\[.*?\]\(images/` + regexp.QuoteMeta(imageName) + `(?:\s+.*?)?\)`)
	var imageContexts []ImageContext
	for index, line := range contentLines {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		locations := re.FindStringIndex(line)
		if locations == nil {
			continue
		}
		startText := line[:locations[0]]
		endText := line[locations[1]:]

		currentLineContext := startText + `<<CURRENT_IMAGE>>` + endText

		paragraphTitle, preContext := node.findImagePreContext(index, contentLines, 1000)

		postContext := node.findImagePostContext(index, contentLines, 1000)

		imageContexts = append(imageContexts, ImageContext{
			ParagraphTitle:     paragraphTitle,
			preContext:         preContext,
			postContext:        postContext,
			currentLineContext: currentLineContext,
		})
	}
	return []ImageContext{}, nil
}

var titleRegex = regexp.MustCompile(`^#{1,6}\s+`)

func (node *PdfToMarkdownNode) findImagePreContext(preIndex int, contentLines []string, maxChars int) (string, string) {
	var preContext []string
	var paragraphTitle string
	totalChars := 0
	for preIndex := preIndex - 1; preIndex >= 0; preIndex-- {
		line := contentLines[preIndex]

		if titleRegex.MatchString(line) {
			paragraphTitle = line
			break
		}
		preContext = append(preContext, line)
		totalChars += len(line)
		if totalChars >= maxChars {
			break
		}
	}

	slices.Reverse(preContext)
	return paragraphTitle, strings.Join(preContext, "\n")
}

func (node *PdfToMarkdownNode) findImagePostContext(postIndex int, contentLines []string, maxChars int) string {
	var postContext []string
	totalChars := 0
	for postIndex := postIndex + 1; postIndex < len(contentLines); postIndex++ {
		line := contentLines[postIndex]

		if titleRegex.MatchString(line) {
			continue
		}
		postContext = append(postContext, line)
		totalChars += len(line)
		if totalChars >= maxChars {
			break
		}
	}
	return strings.Join(postContext, "\n")
}

var vlmSem = make(chan struct{}, 5)

func generateImageDescription(ctx context.Context, docTitle string, imageContext ImageContext, imageBase64 string, mimeType string) (string, error) {

	select {
	case vlmSem <- struct{}{}:
		defer func() {
			<-vlmSem
		}()
	default:
		return "", ctx.Err()
	}
	response, err := chatmodel.VlmModel.Generate(ctx, []*schema.Message{
		{
			Role: schema.System,
			Content: dedent.Dedent(`
			# 你是一个图片内容总结专家
			## 你的任务是根据上下文来为当前图片生成一段检索用的摘要。
			## 规则
			1. 文档名称段落标题均有可能为空
			2. "<<CURRENT_IMAGE>>" 是当前图片在上下文中的位置。上下文里若出现其他本地图片路径（如 "![](images/xxx.jpg)" 或类似写法），那些是别的图
			3. 一句话概括此图片是干嘛的，但不要长篇大论
			4. 不要标题、列表、markdown 或解释你的推理过程
			5. 输出内容不要带「图片展示了 xxx、图片 xxx」这种类似前缀或者说明
			`),
		},
		{
			Role: "user",
			UserInputMultiContent: []schema.MessageInputPart{
				{
					Type: schema.ChatMessagePartTypeImageURL,
					Image: &schema.MessageInputImage{
						MessagePartCommon: schema.MessagePartCommon{
							Base64Data: &imageBase64,
							MIMEType:   mimeType,
						},
					},
				},
				{
					Type: schema.ChatMessagePartTypeText,
					Text: dedent.Dedent(fmt.Sprintf(`
					# 内容
					【文档名称】
					%s
					【段落标题】
					%s
					【上下文】
					%s
					`, docTitle, imageContext.ParagraphTitle, strings.Join(
						[]string{imageContext.preContext, imageContext.currentLineContext, imageContext.postContext},
						"\n"),
					)),
				},
			},
		},
	})
	if err != nil {
		return "", err
	}

	return response.Content, nil

}

func (node *PdfToMarkdownNode) getImageDescription(ctx context.Context, docTitle string, imageContexts []ImageContext, imageName string, imageDir string) ([]string, error) {
	var descriptions []string
	for _, imageContext := range imageContexts {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		contentType := mime.TypeByExtension(filepath.Ext(imageName))

		if contentType == "" {
			return nil, errors.New("unknown image type")
		}
		imageData, err := os.ReadFile(filepath.Join(imageDir, imageName))
		if err != nil {
			return nil, err
		}

		imageBase64 := base64.StdEncoding.EncodeToString(imageData)

		response, err := generateImageDescription(ctx, docTitle, imageContext, imageBase64, contentType)
		if err != nil {
			return nil, err
		}
		descriptions = append(descriptions, response)
	}
	return descriptions, nil
}

func (node *PdfToMarkdownNode) getImageOnlineLink(ctx context.Context, imageName string, imagePath string, mimeType string, docTitle string) (string, error) {

	if mimeType == "" {
		return "", errors.New("unknown image type")
	}
	objectKey := fmt.Sprintf("markdown-images/%s/%s", docTitle, imageName)

	imageFile, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer imageFile.Close()
	_, err = clients.CosClient.Object.Put(
		ctx,
		objectKey,
		imageFile,
		&cos.ObjectPutOptions{
			ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
				ContentType: mimeType,
			},
		},
	)
	if err != nil {
		return "", err
	}
	return *config.AppConfig.CosConfig.Domain + "/" + objectKey, nil

}
