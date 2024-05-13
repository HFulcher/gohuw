package build

import (
	"html/template"
	"huwfulcher/gohuw/internal/search"
	"os"
	"path/filepath"
)

type PageData struct {
	Site  map[string]interface{}
	Page  search.MarkdownFile
	Pages []search.MarkdownFile
    FullUrl string
}

func BuildFiles(files search.Directory, templates map[string]*template.Template, config map[string]interface{}) ([]PageData, error) {
    var pages []PageData

	for _, file := range files.Files {
		var templateName string
		layout, ok := file.Metadata["layout"].(string)

		if !ok {
			templateName = "single.html"
		} else {
			templateName = layout + ".html"
		}

		template, ok := templates[templateName]

		if !ok {
			template = templates["single.html"]
		}

		err := os.MkdirAll(filepath.Dir(file.Path), os.ModePerm)

		if err != nil {
			panic(err)
		}

		writtenFile, err := os.Create(file.Path)

		if err != nil {
			panic(err)
		}

		defer writtenFile.Close()

		pageData := PageData{
			Site:  config,
			Page:  file,
			Pages: files.Files,
            FullUrl: config["Url"].(string) + "/" + file.Slug, 
		}

		err = template.ExecuteTemplate(writtenFile, "baseof", pageData)

		if err != nil {
			panic(err)
		}

        pages = append(pages, pageData)
	}

	for _, subdir := range files.Subdirectories {
        subpages, err := BuildFiles(subdir, templates, config)

        if err != nil {
            panic(err)
        }

        pages = append(pages, subpages...)
	}

	return pages, nil
}

func BuildSitemap(pages []PageData) {
    sitemap, err := os.Create("./public/sitemap.txt")

    if err != nil {
        panic(err)
    }

    for _, page := range pages {
        sitemap.WriteString(page.FullUrl + "\n")
    }

    sitemap.Close()
}
