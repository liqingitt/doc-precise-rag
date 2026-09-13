package processimportnodes

type FileLink struct {
	ObjKey   string
	Url      string
	DocTitle string
}

type DocLink struct {
	ObjKey   string
	Url      string
	DocTitle string
	Ext      string
}

type MarkdownData struct {
	DocTitle string
	ObjKey   string
	Content  string
}

type ImageContext struct {
	ParagraphTitle     string
	preContext         string
	postContext        string
	currentLineContext string
}
