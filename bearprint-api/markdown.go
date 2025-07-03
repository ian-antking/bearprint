package main

import (
	"bytes"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

func renderMarkdown(mdContent []byte) ([]byte, error) {
    md := goldmark.New(
        goldmark.WithExtensions(
            extension.Table,
            highlighting.NewHighlighting(
                highlighting.WithStyle("dracula"),
                highlighting.WithFormatOptions(),
            ),
        ),
        goldmark.WithRendererOptions(
            html.WithHardWraps(),
            html.WithXHTML(),
        ),
    )
    var buf bytes.Buffer
    if err := md.Convert(mdContent, &buf); err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}
