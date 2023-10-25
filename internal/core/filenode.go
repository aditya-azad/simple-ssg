package core

import (
	"github.com/aditya-azad/simple-ssg/pkg/set"
)

type FileNode struct {
	FilePath      string
	Template      string
	TemplateProps map[string]string
	Expands       set.Set[string]
	Blocks        BlockChain
}
