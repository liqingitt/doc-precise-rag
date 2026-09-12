package processimportnodes

import (
	"context"
	"net/url"
	"path"

	"github.com/cloudwego/eino/compose"
)

type EntryNode struct {
}

func (e *EntryNode) Process(fileLink *FileLink) (*DocLink, error) {
	parsedURL, err := url.Parse(fileLink.Url)
	if err != nil {
		return nil, err
	}

	ext := path.Ext(parsedURL.Path)

	return &DocLink{
		ObjKey:   fileLink.ObjKey,
		Url:      fileLink.Url,
		DocTitle: fileLink.DocTitle,
		Ext:      ext,
	}, nil
}

func (e *EntryNode) BuildInvokableLambda() *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, input *FileLink) (*DocLink, error) {
		return e.Process(input)
	})
}

func NewEntryNode() *EntryNode {
	return &EntryNode{}
}
