package build

import (
    "os"
    "path/filepath"
    "huwfulcher/gohuw/internal/search"
    "html/template"
)

type PageData struct {
    Site map[string]interface{}
    Page search.MarkdownFile
    Pages []search.MarkdownFile
}

func BuildFiles(files search.Directory, templates map[string]*template.Template) error {
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

        pageData := PageData {
            Page: file,
            Pages: files.Files,
        }

        err = template.ExecuteTemplate(writtenFile, "baseof", pageData)

        if err != nil {
            panic(err)
        }
    }
    
    for _, subdir := range files.Subdirectories {
        BuildFiles(subdir, templates)
    }
    
    return nil
}
