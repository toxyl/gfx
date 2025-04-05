package markdown

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var (
	parser = goldmark.New(goldmark.WithExtensions(extension.GFM))
)

func ToHTML(source string) (string, error) {
	var buf bytes.Buffer
	if err := parser.Convert([]byte(source), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
