package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/aditya-azad/simple-ssg/pkg/fs"
)

type FileNode struct {
	FilePath  string
	Deps      []string
	DepsProps []map[string]string
	Blocks    *BlockChain
}

// Generates a map of relative path -> FileNode for all the files.
// Runs concurrently internally (it is blocking)
func GenerateFileNodes(inputDir *string) (map[string]FileNode, error) {
	nodes := map[string]FileNode{}
	var wg sync.WaitGroup
	var mut sync.Mutex

	nodeGenerator := func(path string) error {
		defer wg.Done()
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
		mut.Lock()
		nodes[rel] = FileNode{
			FilePath:  rel,
			Deps:      []string{},
			DepsProps: []map[string]string{},
			Blocks:    blocks,
		}
		mut.Unlock()
		return nil
	}

	for _, dir := range []string{filepath.Join(*inputDir, "templates/"), filepath.Join(*inputDir, "pages/")} {
		if filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !fs.IsDir(&path) {
				wg.Add(1)
				go nodeGenerator(path)
			}
			return nil
		}) != nil {
			return nil, errors.New(fmt.Sprintf("Error walking files in %s dir", dir))
		}
	}

	wg.Wait()
	return nodes, nil
}
