package processimportnodes

import (
	"context"
	"strings"

	"github.com/cloudwego/eino/compose"
)

type FileTypeBranchNode struct {
}

func (n *FileTypeBranchNode) BuildGraphBranch() *compose.GraphBranch {

	return compose.NewGraphBranch(
		func(ctx context.Context, input *DocLink) (string, error) {
			if strings.EqualFold(input.Ext, "pdf") {
				return "pdfToMarkdownNode", nil
			} else {
				return "splitDocNode", nil
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
