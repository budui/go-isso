package markdown

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark"
)

// Worker do markdown convert
type Worker struct {
	m goldmark.Markdown
}

// Convert markdown to html
func (w *Worker) Convert(source string) (string, error) {
	var buf bytes.Buffer
	err := w.m.Convert([]byte(source), &buf)
	if err != nil {
		return "", fmt.Errorf("markdown convert failed %w", err)
	}
	return buf.String(), nil
}

// New return a converter worker
func New() *Worker {
	return &Worker{goldmark.New()}
}
