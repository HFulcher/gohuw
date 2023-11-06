package main

import (
	"io/fs"
	"log"
	"path/filepath"
)

func getAssets() []string {
	var files []string
	if err := filepath.WalkDir("./assets", func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}

		files = append(files, s)

		return nil
	}); err != nil {
		log.Println("ERROR", err)
	}

	return files
}
