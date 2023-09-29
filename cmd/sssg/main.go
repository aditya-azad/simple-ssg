package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aditya-azad/simple-ssg/internal/core"
	"github.com/aditya-azad/simple-ssg/pkg/logging"
	"github.com/aditya-azad/simple-ssg/pkg/set"
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
	validDirectories := set.NewSet("public", "templates", "pages")
	validFiles := set.NewSet("globals.yml")
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

func generateFileNodes(inputDir *string) {
	templateNodes := map[string]core.FileNode{}
	pagesNodes := map[string]core.FileNode{}
	templatesDir := filepath.Join(*inputDir, "templates/")
	pagesDir := filepath.Join(*inputDir, "pages/")

	templateTraverser := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logging.Error("Error reading template file: %s", err.Error())
		}
		if !isDir(&path) {
			rel, err := filepath.Rel(templatesDir, path)
			if err != nil {
				logging.Error("Error reading template file: %s", err.Error())
			}
			dat, err := os.ReadFile(path)
			if err != nil {
				logging.Error("Error reading template file: %s", err.Error())
			}
			templateNodes[rel] = core.FileNode{
				Slug:          "",
				Template:      "",
				TemplateProps: map[string]string{},
				Expands:       *set.NewSet(""),
				Blocks:        []core.Block{core.RawBlock{Data: dat}},
			}
		}
		return nil
	}

	pagesTraverser := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logging.Error("Error reading page file: %s", err.Error())
		}
		if !isDir(&path) {
			rel, err := filepath.Rel(pagesDir, path)
			if err != nil {
				logging.Error("Error reading page file: %s", err.Error())
			}
			dat, err := os.ReadFile(path)
			if err != nil {
				logging.Error("Error reading page file: %s", err.Error())
			}
			pagesNodes[rel] = core.FileNode{
				Slug:          rel,
				Template:      "",
				TemplateProps: map[string]string{},
				Expands:       *set.NewSet(""),
				Blocks:        []core.Block{core.RawBlock{Data: dat}},
			}
		}
		return nil
	}

	err := filepath.Walk(templatesDir, templateTraverser)
	if err != nil {
		logging.Error("Error walking files in templates dir: %s", err.Error())
	}

	err = filepath.Walk(pagesDir, pagesTraverser)
	if err != nil {
		logging.Error("Error walking files in pages dir: %s", err.Error())
	}
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
	generateFileNodes(inputDir)
	// parse files
	// compress files
	// files to public
	// compress and copy public files
}
