package main

import (
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func getTemplates() []string {
	var files []string
	if err := filepath.WalkDir("./templates", func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ".html" {
			files = append(files, s)
		}
		return nil
	}); err != nil {
		log.Println("ERROR", err)
	}

	return files
}

func parseTemplates(f []string) {
	templates := make(map[string]*template.Template)

	for _, tmplPath := range f {
		name := filepath.Base(string(tmplPath))
		tmpl, err := template.ParseFiles(tmplPath)
		if err != nil {
			panic(err)
		}
		templates[name] = tmpl
	}

	log.Println(templates)

	applyTemplates(templates)
}

func applyTemplates(tmpls map[string]*template.Template) {
	for _, page := range markdownFiles {
		var templateName string

		layout, ok := page.metadata["layout"].(string)
		// log.Println(layout)

		if !ok {
			// log.Printf("Layout not found for %s, defaulting to single.html", page.title)
			templateName = "single.html"
		} else {
			templateName = layout + ".html"
		}

		// log.Println(templateName)

		tmpl, ok := tmpls[templateName]

		// log.Println(tmpl)

		if !ok {
			panic("Template not found")
		}

		err := os.MkdirAll(filepath.Join("public", filepath.Dir(page.destination)), os.ModePerm)
		if err != nil {
			panic(err)
		}
		outputFile, err := os.Create(filepath.Join("public", page.destination))
		if err != nil {
			panic(err)
		}
		defer outputFile.Close()

		log.Println(tmpl.Name())
		log.Println(page.title)

		// Execute the template with the page's content and metadata
		if err := tmpl.Execute(outputFile, page); err != nil {
			panic(err)
		}

	}
}
