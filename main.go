package main

import (
	"fmt"
	"os"
)

// TODO: Use named options as arg to command line for source and destination
func main() {
	if len(os.Args) > 2 {
		sourceDir := os.Args[1]
		destinationDir := os.Args[2]
		findFiles(sourceDir, destinationDir)
	} else {
		fmt.Println("Usage: sync_and_transcode_music_files <source-directory> <destination-directory>")
		os.Exit(1)
	}
}
