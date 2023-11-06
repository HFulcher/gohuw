package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

var watcher *fsnotify.Watcher

var excludedDirs = []string{
	"node_modules",
	".git",
	"public",
	// Add more directories you want to exclude
}

func startDev(port string, directory string, watchDirectory string) {
	log.Println("Starting dev mode")
	log.Printf("Serving on port %s", port)
	log.Printf("Serving files from %s", directory)
	log.Printf("Watching for changes from %s", watchDirectory)

	http.Handle("/", http.FileServer(http.Dir(directory)))

	errChan := make(chan error, 1)

	go func() {
		errChan <- http.ListenAndServe(":"+port, nil)
	}()

	go func() {
		watcher, _ = fsnotify.NewWatcher()
		defer watcher.Close()

		if err := filepath.Walk(watchDirectory, watchDir); err != nil {
			log.Println("ERROR", err)
		}

		//
		done := make(chan bool)

		//
		go func() {
			for {
				select {
				// watch for events
				case event := <-watcher.Events:
					log.Printf("%s changed, rebuilding", event.Name)

					// watch for errors
				case err := <-watcher.Errors:
					log.Println("ERROR", err)
				}
			}
		}()

		<-done
	}()

	select {
	case err := <-errChan:
		if err != nil {
			log.Fatalf("Errpr: %v", err)
		}
	}
}

func watchDir(path string, fi os.FileInfo, err error) error {

	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		for _, excludedDir := range excludedDirs {
			if strings.HasPrefix(filepath.Base(path), excludedDir) {
				return filepath.SkipDir
			}
		}

		return watcher.Add(path)
	}

	return nil
}
