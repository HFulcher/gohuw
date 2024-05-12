package search

import (
	"bytes"
	"fmt"
	"html/template"
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

type MarkdownFile struct {
	Metadata map[string]interface{}
	Slug     string
	Path     string
	Content  template.HTML
}

type Directory struct {
	Name           string
	Path           string
	Files          []MarkdownFile
	Subdirectories []Directory
}

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

func SearchContent(dir string) (Directory, error) {
	direc := Directory{
		Name: filepath.Base(dir),
		Path: dir,
	}

	var files []MarkdownFile

	entries, err := os.ReadDir(dir)
	if err != nil {
		return Directory{}, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			subdir, err := SearchContent(filepath.Join(dir, entry.Name()))
			if err != nil {
				return Directory{}, err
			}

			direc.Subdirectories = append(direc.Subdirectories, subdir)
		} else if strings.HasSuffix(entry.Name(), ".md") {
			path := filepath.Join(dir, entry.Name())

			in, _ := os.Open(path)
			fileStream, _ := io.ReadAll(in)

			var buf bytes.Buffer
			ctx := parser.NewContext()

			err := md.Convert(fileStream, &buf, parser.WithContext(ctx))
			if err != nil {
				return Directory{}, err
			}

			metaData := meta.Get(ctx)
			destination := "./public/" + strings.Replace(strings.TrimPrefix(path, "content/"), ".md", ".html", -1)
			slug := strings.TrimSuffix(strings.Replace(entry.Name(), ".md", "", -1), "index")

			fmt.Println(destination)
			fmt.Println(slug)

			file := MarkdownFile{
				Metadata: metaData,
				Slug:     slug,
				Path:     destination,
				Content:  template.HTML(buf.String()),
			}

			files = append(files, file)
		}
	}

	direc.Files = files

	return direc, nil
}

func SearchTemplates(dir string) (map[string]*template.Template, error) {
    var rawTemplates []string
    templates := make(map[string]*template.Template)

	err := filepath.WalkDir(dir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(info.Name()) == ".html" {
			rawTemplates = append(rawTemplates, path)
		}

		return nil
	})

	if err != nil {
		return map[string]*template.Template{}, err
	}

    for _, templatePath := range rawTemplates {
        name := filepath.Base(templatePath)
        tmpl, err := template.ParseFiles("./templates/baseof.html", "./"+templatePath)

        if err != nil {
            return map[string]*template.Template{}, err
        }

        templates[name] = tmpl
    }

	return templates, err

}
