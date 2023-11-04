package main

import (
	"flag"
	"os"
	"path/filepath"
	"sync"

	"github.com/aditya-azad/simple-ssg/internal/core"
	"github.com/aditya-azad/simple-ssg/pkg/logging"
	"github.com/aditya-azad/simple-ssg/pkg/set"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"gopkg.in/yaml.v2"
)

func isDir(path *string) bool {
	info, err := os.Stat(*path)
	if err != nil || !info.IsDir() {
		return false
	}
	return true
}

func validateIsDir(path *string) {
	info, err := os.Stat(*path)
	if err != nil || !info.IsDir() {
		logging.Error("\"%s\" is not a valid directory path: %s", *path, err.Error())
	}
}

func validateInputDirectoryStructure(path *string) {
	files, err := os.ReadDir(*path)
	if err != nil {
		logging.Error("Unable to read contents of input directory: %s", err.Error())
	}
	currDir, err := filepath.Abs(*path)
	if err != nil {
		logging.Error("Unable to get absolute path of input directory: %s", err.Error())
	}
	validDirectories := set.New("public", "templates", "pages")
	validFiles := set.New("globals.yml")
	if err != nil {
		logging.Error("Unable to read the input directory: %s", err.Error())
	}
	for _, file := range files {
		filePath := filepath.Join(currDir, file.Name())
		isFileDir := isDir(&filePath)
		if isFileDir && !validDirectories.Has(file.Name()) {
			logging.Error("Illegal directory \"%s\" in input directory", file.Name())
		} else if !isFileDir && !validFiles.Has(file.Name()) {
			logging.Error("Illegal file \"%s\" in input directory", file.Name())
		}
	}
}

func readGlobals(inputDir *string) map[string]string {
	file, err := os.ReadFile(filepath.Join(*inputDir, "globals.yml"))
	if err != nil {
		logging.Error("Error reading globals file: %s", err.Error())
	}
	data := map[string]string{}
	if err := yaml.Unmarshal(file, &data); err != nil {
		logging.Error("Error unmarshalling globals file: %s", err.Error())
	}
	return data
}

func convertMarkdownToHTML(data []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(data)
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	return markdown.Render(doc, renderer)
}

func convertIpynbToHTML(data []byte) []byte {
	// TODO: implement this
	return data
}

func toHTML(data []byte, fileExtension string) []byte {
	htmlData := data
	switch fileExtension {
	case ".md":
		htmlData = convertMarkdownToHTML(data)
	case ".ipynb":
		htmlData = convertIpynbToHTML(data)
	case ".html":
		// do nothing
	default:
		logging.Error("Unknown file type %s received", fileExtension)
	}
	return htmlData
}

// Convert HTML data to chain of blocks
func HTMLToBlocks(data *[]byte) *core.BlockChain {
	bc := core.NewBlockChain()
	isOpen := false
	var start uint64 = 0
	var dataSize uint64 = uint64(len(*data))
	for i := uint64(0); i < dataSize; i += 1 {
		// opening braces
		if i+1 < dataSize && (*data)[i] == byte('{') && (*data)[i+1] == byte('%') {
			if isOpen {
				logging.Error("Invalid syntax, you cannot nest special blocks")
			}
			bc.Append(&core.Block{
				Type: core.BLOCK_HTML,
				Args: []string{string((*data)[start:i])},
			})
			start = i + 2
			isOpen = true
		}
		// closing braces
		if i+1 < dataSize && (*data)[i] == byte('%') && (*data)[i+1] == byte('}') {
			if !isOpen {
				logging.Error("Invalid syntax, you cannot close a unopened block")
			}
			block, err := core.ParseSpecialBlock(data, start, i)
			if err != nil {
				logging.Error(err.Error())
			}
			bc.Append(block)
			start = i + 2
			isOpen = false
		}
	}
	return bc
}

// Generates a map of relative path -> FileNode for all the files
func generateFileNodes(inputDir *string) map[string]core.FileNode {
	nodes := map[string]core.FileNode{}
	var wg sync.WaitGroup
	var mut sync.Mutex

	nodeGenerator := func(path string) {
		defer wg.Done()
		rel, err := filepath.Rel(*inputDir, path)
		if err != nil {
			logging.Error("Error generating relative path: %s", err.Error())
		}
		data, err := os.ReadFile(path)
		if err != nil {
			logging.Error("Error reading file: %s", err.Error())
		}
		data = toHTML(data, filepath.Ext(rel))
		blocks := HTMLToBlocks(&data)
		mut.Lock()
		nodes[rel] = core.FileNode{
			FilePath:  rel,
			Deps:      []string{},
			DepsProps: []map[string]string{},
			Blocks:    blocks,
		}
		mut.Unlock()
	}

	for _, dir := range []string{filepath.Join(*inputDir, "templates/"), filepath.Join(*inputDir, "pages/")} {
		if filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				logging.Error("Error reading file: %s", err.Error())
			}
			if !isDir(&path) {
				wg.Add(1)
				go nodeGenerator(path)
			}
			return nil
		}) != nil {
			logging.Error("Error walking files in %s dir", dir)
		}
	}

	wg.Wait()
	return nodes
}

func main() {
	// parse args
	inputDir := flag.String("in", ".", "Input directory")
	outputDir := flag.String("out", "", "Output directory")
	flag.Parse()

	// validate args
	validateIsDir(inputDir)
	validateIsDir(outputDir)
	validateInputDirectoryStructure(inputDir)

	// read globals file and generate globals
	_ = readGlobals(inputDir)
	// read and convert files
	_ = generateFileNodes(inputDir)
	// parse files
	// compress files
	// files to public
	// compress and copy public files
}
