package main

import (
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

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
			panic("Template not found")
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

		// Execute the template with the page's content and metadata
		if err := tmpl.ExecuteTemplate(outputFile, "baseof", page); err != nil {
			log.Println("Throwing error")
			panic(err)
		}

	}
}
