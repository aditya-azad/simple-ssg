package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/aditya-azad/simple-ssg/internal/core"
	"github.com/aditya-azad/simple-ssg/pkg/fs"
	"github.com/aditya-azad/simple-ssg/pkg/logging"
	"github.com/aditya-azad/simple-ssg/pkg/set"
	"gopkg.in/yaml.v2"
)

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
		isFileDir := fs.IsDir(&filePath)
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
	_, _ = core.GenerateFileNodes(inputDir)
	// parse files
	// compress files
	// files to public
	// compress and copy public files
}
