package main

import (
	"flag"
	"os"
	"path/filepath"

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
		logging.Error("\"%s\" is not a valid directory path", *path)
	}
}

func validateInputDirectoryStructure(path *string) {
	files, err := os.ReadDir(*path)
	if err != nil {
		logging.Error("Unable to read contents of input directory")
	}
	currDir, err := filepath.Abs(*path)
	if err != nil {
		logging.Error("Unable to get absolute path of input directory")
	}
	validDirectories := set.NewSet("public", "templates", "pages")
	validFiles := set.NewSet("globals.yml")
	if err != nil {
		logging.Error("Unable to read the input directory")
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
		logging.Error("Error reading globals file")
	}
	data := map[string]string{}
	if err := yaml.Unmarshal(file, &data); err != nil {
		logging.Error("Error unmarshalling globals file")
	}
	return data
}

func generateHTMLFiles() {
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
	//generatedHTMLFiles := generateHTMLFiles()
	// parse files
	// compress files
	// files to public
	// compress and copy public files
}
