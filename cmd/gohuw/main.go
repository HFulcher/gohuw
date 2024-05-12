package main

import (
	"fmt"
	"huwfulcher/gohuw/internal/search"
    "huwfulcher/gohuw/internal/build"
)

func main() {
	contentFiles, err := search.SearchContent("content")
	templateFiles, err := search.SearchTemplates("templates")

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

    build.BuildFiles(contentFiles, templateFiles)
}

func PrintContent(contentFiles search.Directory) {
	fmt.Println("\n\nDirectory: ", contentFiles.Name)

	for _, file := range contentFiles.Files {
		fmt.Println(file.Path)
	}

	for _, subdir := range contentFiles.Subdirectories {
		PrintContent(subdir)
	}
}
