package processimportnodes

import (
	"context"

	"github.com/cloudwego/eino/compose"
)

type SplitDocNode struct{}

func (s *SplitDocNode) Process(docLink *DocLink) (string, error) {
	return "123", nil
}

func (s *SplitDocNode) BuildInvokableLambda() *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, input *DocLink) (string, error) {
		return s.Process(input)
	})
}

func NewSplitDocNode() *SplitDocNode {
	return &SplitDocNode{}
}
