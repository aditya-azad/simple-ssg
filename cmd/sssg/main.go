package main

import (
	"flag"
	"github.com/aditya-azad/simple-ssg/pkg/logging"
	"os"
)

func validateDirectoryPath(path *string) bool {
	info, err := os.Stat(*path)
	if err != nil {
		return false
	}
	if !info.IsDir() {
		return false
	}
	return true
}

func main() {
	// parse args
	inputDir := flag.String("in", ".", "Input directory")
	outputDir := flag.String("out", "", "Output directory")
	flag.Parse()
	// validate args
	if !validateDirectoryPath(inputDir) {
		logging.Error("Error reading input directory \"%s\"", *inputDir)
	}
	if !validateDirectoryPath(outputDir) {
		logging.Error("Error reading output directory \"%s\"", *outputDir)
	}
}
