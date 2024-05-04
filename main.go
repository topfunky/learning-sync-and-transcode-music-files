package main

import (
	"flag"
)

func main() {
	sourcePtr := flag.String("source", "source", "Directory in which to find original music files")
	destinationPtr := flag.String("destination", "destination", "Output directory for transcoded files")

	flag.Parse()

	sourceDir := *sourcePtr
	destinationDir := *destinationPtr
	findFiles(sourceDir, destinationDir)
}
