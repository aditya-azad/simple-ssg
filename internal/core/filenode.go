package core

import (
	"github.com/aditya-azad/simple-ssg/pkg/set"
)

type FileNode struct {
	slug          string
	template      string
	templateProps map[string]string
	expands       set.Set[string]
	blocks        []Block
}
