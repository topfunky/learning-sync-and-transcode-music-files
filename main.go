package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	sourcePtr := flag.String("source", "source", "Directory in which to find original music files")
	destinationPtr := flag.String("destination", "destination", "Output directory for transcoded files")
	dryRunPtr := flag.Bool("dry-run", false, "Show which duplicate files would be deleted without deleting them")

	flag.Parse()

	sourceDir := *sourcePtr
	destinationDir := *destinationPtr
	dryRun := *dryRunPtr

	if err := findAndTranscodeFiles(sourceDir, destinationDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := removeDuplicateFiles(destinationDir, dryRun); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// TODO: Validate after running; display list of files that did not end up in the destination
// TODO: Print list of files that resulted in errors
