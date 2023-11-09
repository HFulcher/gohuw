package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

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
		html.WithUnsafe(),
	),
)

type MarkdownFile struct {
	Title       string
	Slug        string
	Path        string
	Destination string
	Content     template.HTML
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
			Content:     template.HTML(buf.String()),
			Metadata:    metaData,
		}

		if title, ok := metaData["title"].(string); ok {
			page.Title = title
		}

		if date, ok := metaData["date"].(string); ok {
			page.Metadata["date"] = getDate(date)
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
				Content:     template.HTML(buf.String()),
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

		sort.Slice(mdFiles[:], func(i, j int) bool {
			iDate, _ := time.Parse(time.RFC3339, mdFiles[i].Metadata["date"].(string))
			jDate, _ := time.Parse(time.RFC3339, mdFiles[j].Metadata["date"].(string))

			// Use Before method to compare the time values.
			return iDate.After(jDate)
		})

		for _, file := range mdFiles {
			if date, ok := file.Metadata["date"].(string); ok {
				file.Metadata["date"] = getDate(date)
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

func getDate(f string) string {
	t, err := time.Parse(time.RFC3339, f)
	if err != nil {
		panic(err)
	}

	// Format the date to "Oct 11, 2022"
	formattedDate := fmt.Sprintf("%v %v, %v", t.Month().String()[:3], t.Day(), t.Year())

	return formattedDate
}
