package main

import (
	"bytes"
	"io"
	"io/fs"
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
	Title       string
	Slug        string
	Path        string
	Destination string
	Content     string
	Metadata    map[string]interface{}
}

type ContentType struct {
	IndexFile    MarkdownFile
	ContentFiles []MarkdownFile
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

func getContent() map[string][]string {
	paths := make(map[string][]string)

	items, err := os.ReadDir("./content")
	if err != nil {
		panic(err)
	}

	for _, item := range items {
		if item.IsDir() {
			var subItems []string

			if err := filepath.WalkDir("./content/"+item.Name(), func(s string, d fs.DirEntry, e error) error {
				if e != nil {
					return e
				}

				if filepath.Ext(d.Name()) == ".md" {
					subItems = append(subItems, s)
				}
				return nil
			}); err != nil {
				panic(err)
			}

			paths[item.Name()] = subItems
		}
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
		slug := strings.TrimSuffix(strings.Replace(strings.TrimPrefix(destination, "public/"), ".html", "", -1), "index")

		metaData := meta.Get(ctx)

		page := MarkdownFile{
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

func parseContent(f map[string][]string) map[string]ContentType {
	parsedFiles := make(map[string]ContentType)

	for key, files := range f {
		var indexFile MarkdownFile
		var mdFiles []MarkdownFile

		for _, file := range files {
			in, _ := os.Open(file)
			fileStream, _ := io.ReadAll(in)

			var buf bytes.Buffer
			ctx := parser.NewContext()

			if err := md.Convert(fileStream, &buf, parser.WithContext(ctx)); err != nil {
				panic(err)
			}

			destination := strings.Replace(strings.TrimPrefix(file, "content/"), ".md", ".html", -1)
			slug := strings.TrimSuffix(strings.Replace(strings.TrimPrefix(destination, "public/"), ".html", "", -1), "index")

			metaData := meta.Get(ctx)

			page := MarkdownFile{
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

			if getFilenameWithoutExt(file) == "index" {
				indexFile = page
			} else {
				mdFiles = append(mdFiles, page)
			}

		}

		parsedFiles[key] = ContentType{
			IndexFile:    indexFile,
			ContentFiles: mdFiles,
		}
	}

	return parsedFiles
}

func getFilenameWithoutExt(path string) string {
	// Get the filename with extension.
	filenameWithExt := filepath.Base(path)
	// Get the extension.
	ext := filepath.Ext(filenameWithExt)
	// Return the filename without the extension.
	return filenameWithExt[0 : len(filenameWithExt)-len(ext)]
}
