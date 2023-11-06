package main

import (
	"bytes"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var md = goldmark.New(
	goldmark.WithExtensions(
		meta.Meta,
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
	goldmark.WithRendererOptions(
		html.WithHardWraps(),
		html.WithXHTML(),
	),
)

type MarkdownFile struct {
	title       string
	path        string
	destination string
	content     string
	metadata    map[string]interface{}
}

func getContent() []string {
	var files []string
	if err := filepath.WalkDir("./content", func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ".md" {
			files = append(files, s)
		}
		return nil
	}); err != nil {
		log.Println("ERROR", err)
	}

	return files
}

func parseMarkdown(f []string) {
	for _, file := range f {
		in, _ := os.Open(file)
		fileStream, _ := io.ReadAll(in)

		var buf bytes.Buffer
		ctx := parser.NewContext()

		if err := md.Convert(fileStream, &buf, parser.WithContext(ctx)); err != nil {
			panic(err)
		}

		destination := strings.Replace(strings.TrimPrefix(file, "content/"), ".md", ".html", -1)

		metaData := meta.Get(ctx)

		page := MarkdownFile{
			title:       "",
			path:        file,
			destination: destination,
			content:     buf.String(),
			metadata:    metaData,
		}

		if title, ok := metaData["title"].(string); ok {
			page.title = title
		}

		markdownFiles = append(markdownFiles, page)
	}
}
