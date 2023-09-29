package core

import (
	"github.com/aditya-azad/simple-ssg/pkg/set"
)

type FileNode struct {
	Slug          string
	Template      string
	TemplateProps map[string]string
	Expands       set.Set[string]
	Blocks        []Block
}
