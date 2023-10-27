package core

type FileNode struct {
	FilePath      string
	Template      string
	TemplateProps map[string]string
	Expands       []string
	ExpandsProps  []map[string]string
	Blocks        BlockChain
}
