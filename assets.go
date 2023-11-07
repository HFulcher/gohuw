package main

import (
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func getAssets() []string {
	var paths []string

	if err := filepath.WalkDir("./static", func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}

		if !d.IsDir() {
			paths = append(paths, s)
		}

		return nil
	}); err != nil {
		log.Println("ERROR", err)
	}

	return paths
}

func parseAssets(f []string) {
	for _, assetPath := range f {
		err := os.MkdirAll(filepath.Join("public", filepath.Dir(assetPath)), os.ModePerm)
		if err != nil {
			panic(err)
		}

		source, err := os.Open(assetPath)
		if err != nil {
			panic(err)
		}
		defer source.Close()

		destination, err := os.Create(filepath.Join("public", assetPath))
		if err != nil {
			panic(err)
		}
		defer destination.Close()

		_, err = io.Copy(destination, source)

		if err != nil {
			panic(err)
		}
	}
}
