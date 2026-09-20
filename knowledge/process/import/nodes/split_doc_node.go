package processimportnodes

import (
	"context"
	"doc-precise-rag/knowledge/config"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/cloudwego/eino/compose"
)

type SplitDocNode struct {
	logger *slog.Logger
}

func (node *SplitDocNode) Process(markdownData *MarkdownData) (string, error) {
	node.logger.Info("开始分割文档")

	contentLines := strings.Split(markdownData.Content, "\n")
	paragraphs, err := node.splitParagraph(contentLines, markdownData.DocTitle)
	if err != nil {
		node.logger.Error("文档切分失败", "error", err.Error())
		return "", err
	}

	tempJsonDir, err := filepath.Abs(filepath.Join(
		*config.AppConfig.ImportProcessConfig.ImportFileTempDir,
		time.Now().Format("2006-01-02"),
	))

	if err != nil {
		node.logger.Error("目录获取失败", "error", err.Error())
		return "", err
	}

	err = os.MkdirAll(tempJsonDir, 0o755)
	if err != nil {
		node.logger.Error("目录创建失败", "error", err.Error())
		return "", err
	}

	tempJsonPath := filepath.Join(tempJsonDir, "debuger.json")
	node.logger.Info("开始临时写入切割结果", "存储路径", tempJsonPath)

	jsonStr, err := json.MarshalIndent(paragraphs, "", "  ")
	if err != nil {
		node.logger.Error("序列化失败", "error", err.Error())
		return "", err
	}
	err = os.WriteFile(tempJsonPath, jsonStr, 0o644)
	if err != nil {
		node.logger.Error("写入失败", "error", err.Error())
		return "", err
	}
	return markdownData.ObjKey, nil
}

func (node *SplitDocNode) splitParagraph(contentLines []string, docTitle string) ([]*ParagraphData, error) {
	paragraphs := []*ParagraphData{}

	fence := false // 当前是否在代码块中
	fenceRegex := regexp.MustCompile(`^[ \t]{0,3}(\x60{3,}|~{3,})`)
	fenceStart := ""
	contentList := []string{}

	currentLevel := 1
	levelTitles := [7]string{"", "", "", "", "", "", ""}

	titleRegex := regexp.MustCompile(`^(#{1,6})\s+(.*)$`)
	for _, line := range contentLines {

		// 代码围栏检测
		codeMatch := fenceRegex.FindStringSubmatch(line)
		if codeMatch != nil {
			fenceChats := codeMatch[1]

			if !fence {
				fence = true
				fenceStart = fenceChats
			} else {
				if fenceChats[0] == fenceStart[0] && len(fenceChats) >= len(fenceStart) {
					fence = false
					fenceStart = ""
				}
			}
		}
		// 标题检测
		titleMatch := titleRegex.FindStringSubmatch(line)
		if titleMatch != nil && !fence {
			// 检测到当前不是代码围栏但又是标题时，需要对上一段内容进行整合
			if levelTitles[currentLevel] != "" || len(contentList) != 0 {
				data := &ParagraphData{
					DocTitle:         docTitle,
					ParagraphTitle:   levelTitles[currentLevel],
					ParagraphContent: strings.Join(contentList, "\n"),
					ParentTitle:      levelTitles[currentLevel-1],
				}
				if data.ParagraphTitle == "" {
					data.ParagraphTitle = docTitle
				}
				if data.ParentTitle == "" {
					data.ParentTitle = docTitle
				}
				paragraphs = append(paragraphs, data)
			}

			// 收集完成对状态进行一轮重置
			contentList = []string{}
			currentLevel = len(titleMatch[1])
			levelTitles[currentLevel] = titleMatch[2]
			for i := currentLevel + 1; i < len(levelTitles); i++ {
				levelTitles[i] = ""
			}
		} else {
			contentList = append(contentList, line)
		}
	}

	// 对最后部分内容再收集一次
	data := &ParagraphData{
		DocTitle:         docTitle,
		ParagraphTitle:   levelTitles[currentLevel],
		ParagraphContent: strings.Join(contentList, "\n"),
		ParentTitle:      levelTitles[currentLevel-1],
	}
	if data.ParagraphTitle == "" {
		data.ParagraphTitle = docTitle
	}
	if data.ParentTitle == "" {
		data.ParentTitle = docTitle
	}
	paragraphs = append(paragraphs, data)

	return paragraphs, nil
}

func (node *SplitDocNode) BuildInvokableLambda() *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, input *MarkdownData) (string, error) {
		return node.Process(input)
	})
}

func NewSplitDocNode() *SplitDocNode {
	return &SplitDocNode{
		logger: slog.With("流程类型", "文档导入", "节点名称", "splitDocNode"),
	}
}
