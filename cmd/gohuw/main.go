package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"huwfulcher/gohuw/internal/assets"
	"huwfulcher/gohuw/internal/build"
	"huwfulcher/gohuw/internal/search"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var watcher *fsnotify.Watcher

var excludedDirs = []string{
	"node_modules",
	".git",
	"public",
	// Add more directories you want to exclude
}

var config map[string]interface{}

func main() {
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	json.Unmarshal([]byte(configFile), &config)

	devCmd := flag.NewFlagSet("dev", flag.ExitOnError)

	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "dev":
			config["IsProduction"] = false
            config["Url"] = "http://localhost"
			devCmd.Parse(os.Args[2:])
			startDev()

		default:
			fmt.Println("Command not recognised. Please try again")
			os.Exit(1)
		}
	} else {
		config["IsProduction"] = true
		buildSite()
	}

}

func buildSite() {
	contentFiles, err := search.SearchContent("content")
	templateFiles, err := search.SearchTemplates("templates")
	assets.GetAssets()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

    pages, err := build.BuildFiles(contentFiles, templateFiles, config)

    if err != nil {
        panic(err)
    }

    build.BuildSitemap(pages)
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

func startDev() {
	config["Url"] = "http://localhost:8100"

	log.Println("Starting dev mode")
	buildSite()
	log.Printf("Serving on port %s", "8100")
	log.Printf("Serving files from %s", "public")
	log.Printf("Watching for changes from %s", ".")

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", removeHTMLExtension(fs))

	errChan := make(chan error, 1)

	go func() {
		errChan <- http.ListenAndServe(":8100", nil)
	}()

	go func() {
		watcher, _ = fsnotify.NewWatcher()
		defer watcher.Close()

		if err := filepath.Walk(".", watchDir); err != nil {
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
					buildSite()

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
