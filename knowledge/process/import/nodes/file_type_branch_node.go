package processimportnodes

import (
	"context"
	"errors"
	"strings"

	"github.com/cloudwego/eino/compose"
)

type FileTypeBranchNode struct {
}

func (n *FileTypeBranchNode) BuildGraphBranch() *compose.GraphBranch {

	return compose.NewGraphBranch(
		func(ctx context.Context, input *DocLink) (string, error) {
			if strings.EqualFold(input.Ext, ".PDF") {
				return "pdfToMarkdownNode", nil
			} else {
				return "", errors.New("暂未支持的文件类型")
			}
		},
		map[string]bool{
			"pdfToMarkdownNode": true,
			"splitDocNode":      true,
		},
	)

}

func NewFileTypeBranchNode() *FileTypeBranchNode {
	return &FileTypeBranchNode{}
}
