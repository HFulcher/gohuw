package main

import (
	"log"
	"net/http"
	"os"
	"path"
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
	config["Url"] = "http://localhost:" + port

	log.Println("Starting dev mode")
	build()
	log.Printf("Serving on port %s", port)
	log.Printf("Serving files from %s", directory)
	log.Printf("Watching for changes from %s", watchDirectory)

	fs := http.FileServer(http.Dir(directory))
	http.Handle("/", removeHTMLExtension(fs))

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
					build()

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

func removeHTMLExtension(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// If the request is for a root or a directory, serve as is
		if r.URL.Path == "/" || strings.HasSuffix(r.URL.Path, "/") {
			next.ServeHTTP(w, r)
			return
		}

		// Clean the path to prevent directory traversal
		urlPath := path.Clean(r.URL.Path)

		// Check if the requested file exists with .html extension
		if _, err := http.Dir("./public").Open(urlPath + ".html"); err == nil {
			// Rewrite the URL path to include the .html extension
			r.URL.Path += ".html"
		} else if _, err := http.Dir("./public").Open(urlPath + "/index.html"); err == nil {
			// Handle directory URLs by trying to serve index.html
			r.URL.Path += "/index.html"
		}

		// Call the next handler, which in this case is our file server
		next.ServeHTTP(w, r)
	}
}
