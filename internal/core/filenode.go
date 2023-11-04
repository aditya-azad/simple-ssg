package core

type FileNode struct {
	FilePath  string
	Deps      []string
	DepsProps []map[string]string
	Blocks    *BlockChain
}
