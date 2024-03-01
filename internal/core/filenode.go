package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aditya-azad/simple-ssg/pkg/fs"
	"github.com/aditya-azad/simple-ssg/pkg/logging"
)

type FileNode struct {
	FilePath  string
	Deps      []string
	Vars      []string
	DepsProps []map[string]string
	Blocks    *blockChain
}

// Generates a map of relative path -> FileNode for all the files.
// TODO: Runs concurrently internally (it is blocking)
func GenerateFileNodes(inputDir *string) (map[string]FileNode, error) {
	nodes := map[string]FileNode{}

	for _, dir := range []string{filepath.Join(*inputDir, "templates/"), filepath.Join(*inputDir, "pages/")} {
		if filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !fs.IsDir(&path) {
				rel, err := filepath.Rel(*inputDir, path)
				if err != nil {
					return err
				}
				data, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				data, err = ToHTML(data, filepath.Ext(rel))
				if err != nil {
					return err
				}
				blocks, err := HTMLToBlocks(&data)
				if err != nil {
					return err
				}
				nodes[rel] = FileNode{
					FilePath:  rel,
					Deps:      []string{},
					DepsProps: []map[string]string{},
					Blocks:    blocks,
				}
			}
			return nil
		}) != nil {
			return nil, errors.New(fmt.Sprintf("Error walking files in %s dir", dir))
		}
	}

	for name, data := range nodes {
		logging.Info("%s\n%s", name, data.Blocks.toString(true))
		logging.Info("")
	}

	return nodes, nil
}
