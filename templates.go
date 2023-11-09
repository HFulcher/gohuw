package main

import (
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type SingleData struct {
	Site map[string]interface{}
	Page MarkdownFile
}

type ContentData struct {
	Site  map[string]interface{}
	Page  MarkdownFile
	Pages []MarkdownFile
}

func getTemplates() []string {
	var paths []string

	if err := filepath.WalkDir("./templates", func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ".html" {
			if d.Name() != "baseof.html" {
				paths = append(paths, s)
			}
		}
		return nil
	}); err != nil {
		log.Println("ERROR", err)
	}

	return paths
}

func parseTemplates(f []string) map[string]*template.Template {
	tmpls := make(map[string]*template.Template)
	for _, tmplPath := range f {
		name := filepath.Base(string(tmplPath))
		tmpl, err := template.ParseFiles("./templates/baseof.html", "./"+tmplPath)
		if err != nil {
			panic(err)
		}
		tmpls[name] = tmpl
	}

	return tmpls
}

func buildFiles(files []MarkdownFile, tmpls map[string]*template.Template) {
	for _, page := range files {
		var templateName string

		layout, ok := page.Metadata["layout"].(string)

		if !ok {
			templateName = "single.html"
		} else {
			templateName = layout + ".html"
		}

		tmpl, ok := tmpls[templateName]

		if !ok {
			tmpl = tmpls["single.html"]
		}

		topPath := "public"
		filename := page.Destination

		if !strings.HasPrefix(filename, "index") {
			topPath = filepath.Join(topPath, strings.Split(filename, ".")[0])
			filename = "index.html"
		}

		err := os.MkdirAll(filepath.Join(topPath, filepath.Dir(filename)), os.ModePerm)
		if err != nil {
			panic(err)
		}

		outputFile, err := os.Create(filepath.Join(topPath, filename))
		if err != nil {
			panic(err)
		}
		defer outputFile.Close()

		pageData := SingleData{
			Site: config,
			Page: page,
		}

		// Execute the template with the page's content and metadata
		if err := tmpl.ExecuteTemplate(outputFile, "baseof", pageData); err != nil {
			log.Println("Throwing error")
			panic(err)
		}

	}
}

func buildContent(files map[string]ContentType, tmpls map[string]*template.Template) {
	for _, content := range files {
		var templateName string

		indexFile := content.IndexFile

		layout, ok := indexFile.Metadata["layout"].(string)

		if !ok {
			templateName = "list.html"
		} else {
			templateName = layout + ".html"
		}

		tmpl, ok := tmpls[templateName]

		if !ok {
			tmpl = tmpls["list.html"]
		}

		err := os.MkdirAll(filepath.Join("public", filepath.Dir(indexFile.Destination)), os.ModePerm)
		if err != nil {
			panic(err)
		}

		outputFile, err := os.Create(filepath.Join("public", indexFile.Destination))
		if err != nil {
			panic(err)
		}
		defer outputFile.Close()

		contentData := ContentData{
			Site:  config,
			Page:  indexFile,
			Pages: content.ContentFiles,
		}

		// Execute the template with the page's content and metadata
		if err := tmpl.ExecuteTemplate(outputFile, "baseof", contentData); err != nil {
			log.Println("Throwing error")
			panic(err)
		}

		for _, page := range content.ContentFiles {
			var templateName string

			layout, ok := page.Metadata["layout"].(string)

			if !ok {
				templateName = "single.html"
			} else {
				templateName = layout + ".html"
			}

			tmpl, ok := tmpls[templateName]

			if !ok {
				tmpl = tmpls["single.html"]
			}

			err := os.MkdirAll(filepath.Join("public", filepath.Dir(page.Destination)), os.ModePerm)
			if err != nil {
				panic(err)
			}

			outputFile, err := os.Create(filepath.Join("public", page.Destination))
			if err != nil {
				panic(err)
			}
			defer outputFile.Close()

			singleData := SingleData{
				Site: config,
				Page: page,
			}

			// Execute the template with the page's content and metadata
			if err := tmpl.ExecuteTemplate(outputFile, "baseof", singleData); err != nil {
				log.Println("Throwing error")
				panic(err)
			}
		}
	}

}
