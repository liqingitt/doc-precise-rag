package processimport

import (
	"context"
	processimportnodes "doc-precise-rag/knowledge/process/import/nodes"

	"github.com/cloudwego/eino/compose"
)

var ImportGraphCompiledInstance compose.Runnable[*processimportnodes.FileLink, string]

func init() {
	graph := compose.NewGraph[*processimportnodes.FileLink, string]()

	entryNode := processimportnodes.NewEntryNode()
	fileTypeBranchNode := processimportnodes.NewFileTypeBranchNode()
	pdfToMarkdownNode := processimportnodes.NewPdfToMarkdownNode()
	readMarkdownLinkNode := processimportnodes.NewReadMarkdownLinkNode()
	splitDocNode := processimportnodes.NewSplitDocNode()

	// compose.NewGraphBranch()

	graph.AddLambdaNode("entryNode", entryNode.BuildInvokableLambda())
	graph.AddLambdaNode("pdfToMarkdownNode", pdfToMarkdownNode.BuildInvokableLambda())
	graph.AddLambdaNode("readMarkdownLinkNode", readMarkdownLinkNode.BuildInvokableLambda())
	graph.AddLambdaNode("splitDocNode", splitDocNode.BuildInvokableLambda())

	graph.AddEdge(compose.START, "entryNode")
	graph.AddBranch("entryNode", fileTypeBranchNode.BuildGraphBranch())
	graph.AddEdge("pdfToMarkdownNode", "splitDocNode")
	graph.AddEdge("readMarkdownLinkNode", "splitDocNode")
	graph.AddEdge("splitDocNode", compose.END)

	compiled, err := graph.Compile(context.Background())
	if err != nil {
		panic(err)
	}

	ImportGraphCompiledInstance = compiled
}
