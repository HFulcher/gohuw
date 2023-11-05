package main

import (
	"flag"
	"fmt"
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
	// Add more directories you want to exclude
}

func main() {
	devCmd := flag.NewFlagSet("dev", flag.ExitOnError)
	devPort := devCmd.String("p", "8100", "port to serve on")
	devDirectory := devCmd.String("d", ".", "directory to serve files from")

	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "dev":
			devCmd.Parse(os.Args[2:])
			startDev(*devPort, *devDirectory)

		default:
			fmt.Println("Command not recognised. Please try again")
			os.Exit(1)
		}
	} else {
		fmt.Println("...")
	}
}

func startDev(port string, directory string) {
	http.Handle("/", http.FileServer(http.Dir(directory)))

	errChan := make(chan error, 1)

	go func() {
		log.Printf("serving")
		errChan <- http.ListenAndServe(":"+port, nil)
	}()

	go func() {
		watcher, _ = fsnotify.NewWatcher()
		defer watcher.Close()

		if err := filepath.Walk(".", watchDir); err != nil {
			fmt.Println("ERROR", err)
		}

		log.Println("Directory walked")

		//
		done := make(chan bool)

		//
		go func() {
			for {
				select {
				// watch for events
				case event := <-watcher.Events:
					fmt.Printf("EVENT! %#v\n", event)

					// watch for errors
				case err := <-watcher.Errors:
					fmt.Println("ERROR", err)
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
