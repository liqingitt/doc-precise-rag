package processimportnodes

import (
	"context"
	"log/slog"

	"github.com/cloudwego/eino/compose"
)

type PdfToMarkdownNode struct{}

func (f *PdfToMarkdownNode) Process(docLink *DocLink) (*DocLink, error) {
	slog.Info("我开始解析PDF了")
	return nil, nil
}

func (f *PdfToMarkdownNode) BuildInvokableLambda() *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, input *DocLink) (*DocLink, error) {
		return f.Process(input)
	})
}

func NewPdfToMarkdownNode() *PdfToMarkdownNode {
	return &PdfToMarkdownNode{}
}
