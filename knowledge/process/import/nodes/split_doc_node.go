package processimportnodes

import (
	"context"

	"github.com/cloudwego/eino/compose"
)

type SplitDocNode struct{}

func (s *SplitDocNode) Process(markdownData *MarkdownData) (string, error) {
	return markdownData.ObjKey, nil
}

func (s *SplitDocNode) BuildInvokableLambda() *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, input *MarkdownData) (string, error) {
		return s.Process(input)
	})
}

func NewSplitDocNode() *SplitDocNode {
	return &SplitDocNode{}
}
