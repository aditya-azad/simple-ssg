package core

import (
	"errors"
	"fmt"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

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

// Converts markdown, jupyter notebooks to HTML
// If html file is passed, it is ignored
// All other files return an error
func ToHTML(data []byte, fileExtension string) ([]byte, error) {
	htmlData := data
	switch fileExtension {
	case ".md":
		htmlData = convertMarkdownToHTML(data)
	case ".ipynb":
		htmlData = convertIpynbToHTML(data)
	case ".html":
		// do nothing
	default:
		return nil, errors.New(fmt.Sprintf("Unknown file type %s received", fileExtension))
	}
	return htmlData, nil
}
