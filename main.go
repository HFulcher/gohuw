package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	devCmd := flag.NewFlagSet("dev", flag.ExitOnError)
	devPort := devCmd.String("p", "8100", "port to serve on")
	devDirectory := devCmd.String("d", "./public", "directory to serve files from")

	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "dev":
			fmt.Println("dev command")
			devCmd.Parse(os.Args[2:])
			fmt.Println("p: ", *devPort)
			fmt.Println("d: ", *devDirectory)

			http.Handle("/", http.FileServer(http.Dir(*devDirectory)))

			log.Printf("serving")
			log.Fatal(http.ListenAndServe(":"+*devPort, nil))
		default:
			fmt.Println("Command not recognised. Please try again")
			os.Exit(1)
		}
	} else {
		fmt.Println("...")
	}
}
