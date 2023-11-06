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
			if d.Name() != "baseof.html" {
				files = append(files, s)
			}
		}
		return nil
	}); err != nil {
		log.Println("ERROR", err)
	}

	log.Println(files)

	return files
}

func parseTemplates(f []string) {
	templates := make(map[string]*template.Template)

	for _, tmplPath := range f {
		name := filepath.Base(string(tmplPath))
		log.Println(tmplPath)
		tmpl, err := template.ParseFiles("./templates/baseof.html", "./"+tmplPath)
		if err != nil {
			panic(err)
		}
		templates[name] = tmpl
	}

	applyTemplates(templates)
}

func applyTemplates(tmpls map[string]*template.Template) {
	for _, page := range markdownFiles {
		var templateName string

		layout, ok := page.Metadata["layout"].(string)

		if !ok {
			log.Printf("Layout not found for %s, defaulting to single.html", page.Title)
			templateName = "single.html"
		} else {
			log.Printf("Layout found for %s", page.Title)
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
