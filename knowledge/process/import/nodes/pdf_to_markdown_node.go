package processimportnodes

import (
	"context"
	"doc-precise-rag/knowledge/clients"
	"doc-precise-rag/knowledge/config"
	"doc-precise-rag/knowledge/utils"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/cloudwego/eino/compose"
)

type PdfToMarkdownNode struct {
	logger *slog.Logger
}

func (node *PdfToMarkdownNode) Process(docLink *DocLink) (*DocLink, error) {
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
		time.Sleep(2 * time.Second)
		statusResp, err := clients.MineruClient.GetParserFileStatus(*task.Data.TaskID)
		if err != nil {
			node.logger.Error("获取解析PDF任务状态失败", "error", err)
			return nil, err
		}
		if *statusResp.Data.State == "done" {
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

		if *statusResp.Data.State == "failed" {
			node.logger.Error("解析PDF任务状态为解析失败", "error", *statusResp.Data.ErrorMsg)
			return nil, errors.New(*statusResp.Data.ErrorMsg)
		}

		if *statusResp.Data.State == "converting" {
			node.logger.Info("解析PDF任务状态为转换中,继续等待")
			return nil, errors.New(*statusResp.Data.ErrorMsg)
		}

	}

	if fullZipURL == "" {
		node.logger.Error("获取解析PDF任务状态失败", "error", errors.New("fullZipURL为空"))
		return nil, errors.New("fullZipURL为空")
	}

	node.logger.Info("解析成功，开始下载解析后内容", "fullZipURL", fullZipURL)

	unzipDir, err := node.downloadAndUnzip(fullZipURL, *task.Data.TaskID)
	if err != nil {
		node.logger.Error("下载解析后内容失败", "error", err)
		return nil, err
	}

	node.logger.Info("下载解析后内容成功", "unzipDir", unzipDir)

	return docLink, nil
}

func (node *PdfToMarkdownNode) BuildInvokableLambda() *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, input *DocLink) (*DocLink, error) {
		return node.Process(input)
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
