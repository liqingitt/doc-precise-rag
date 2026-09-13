package processimportnodes

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/compose"
)

type ReadMarkdownLinkNode struct {
}

func (r *ReadMarkdownLinkNode) Process(docLink *DocLink) (*MarkdownData, error) {
	fmt.Println("暂未实现该功能，读取markdown链接")
	return &MarkdownData{
		ObjKey:   docLink.ObjKey,
		DocTitle: docLink.DocTitle,
	}, nil
}

func (r *ReadMarkdownLinkNode) BuildInvokableLambda() *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, input *DocLink) (*MarkdownData, error) {
		return r.Process(input)
	})
}

func NewReadMarkdownLinkNode() *ReadMarkdownLinkNode {
	return &ReadMarkdownLinkNode{}
}
