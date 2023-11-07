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
		html.WithUnsafe(),
	),
)

type MarkdownFile struct {
	Site        map[string]interface{}
	Title       string
	Slug        string
	Path        string
	Destination string
	Content     string
	Metadata    map[string]interface{}
}

func getSingle() []string {
	var paths []string

	files, err := os.ReadDir(".")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".md" {
			paths = append(paths, file.Name())
		}
	}

	return paths
}

func getContent() []string {
	var paths []string

	if err := filepath.WalkDir("./content", func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ".md" {
			paths = append(paths, s)
		}
		return nil
	}); err != nil {
		log.Println("ERROR", err)
	}

	return paths
}

func parseMarkdown(f []string) []MarkdownFile {
	var mdFiles []MarkdownFile

	for _, file := range f {
		in, _ := os.Open(file)
		fileStream, _ := io.ReadAll(in)

		var buf bytes.Buffer
		ctx := parser.NewContext()

		if err := md.Convert(fileStream, &buf, parser.WithContext(ctx)); err != nil {
			panic(err)
		}

		destination := strings.Replace(strings.TrimPrefix(file, "content/"), ".md", ".html", -1)
		slug := strings.Replace(strings.TrimPrefix(destination, "public/"), ".html", "", -1)

		metaData := meta.Get(ctx)

		page := MarkdownFile{
			Site:        config,
			Title:       "",
			Slug:        slug,
			Path:        file,
			Destination: destination,
			Content:     buf.String(),
			Metadata:    metaData,
		}

		if title, ok := metaData["title"].(string); ok {
			page.Title = title
		}

		mdFiles = append(mdFiles, page)
	}

	return mdFiles
}

func getFilenameWithoutExt(path string) string {
	// Get the filename with extension.
	filenameWithExt := filepath.Base(path)
	// Get the extension.
	ext := filepath.Ext(filenameWithExt)
	// Return the filename without the extension.
	return filenameWithExt[0 : len(filenameWithExt)-len(ext)]
}
