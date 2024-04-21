package main

import (
	"os"
	"fmt"

)

func main() {
	directory := ""
	if len(os.Args) > 1 {
		directory = os.Args[1]

		findFiles(directory)
	} else {
		fmt.Println("Usage: sync_and_transcode_music_files <directory>")
		os.Exit(1)
	}
}
