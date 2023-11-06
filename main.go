// TODO: Top level search for pages, all other content that is "blog friendly" then stored in /content

package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	devCmd := flag.NewFlagSet("dev", flag.ExitOnError)
	devPort := devCmd.String("p", "8100", "port to serve on")
	devDirectory := devCmd.String("s", "public", "directory to serve files from")
	devWatch := devCmd.String("w", ".", "directory to watch from")

	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "dev":
			devCmd.Parse(os.Args[2:])
			startDev(*devPort, *devDirectory, *devWatch)

		default:
			fmt.Println("Command not recognised. Please try again")
			os.Exit(1)
		}
	} else {
		build()
	}
}
