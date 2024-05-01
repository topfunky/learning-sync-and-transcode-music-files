package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xfrr/goffmpeg/transcoder"
)

// TODO: Get rid of this struct and just work with strings
// TODO: Fields should be titlecase if meant to be public (probably not?)
type FileToRender struct {
	sourcePath      string
	destinationPath string
}

// findFiles traverses the specified directory and transcodes all .m4a files to .mp3 format.
func findFiles(sourceDir, destinationDir string) {
	fmt.Printf("🔍 Finding files in source directory %s\n", sourceDir)
	// TODO: Needs to either look for existence of .m4a or compareDirectories() should be rewritten to return source file name
	// TODO: But if .mp3 exists as source, then it should be copied to the destination as-is
	filesThatNeedToBeRendered, err := compareDirectories(sourceDir, destinationDir)
	fmt.Printf("📂 FilesThatNeedToBeRendered %v\n", filesThatNeedToBeRendered)

	if err != nil {
		fmt.Println("Error:", err)
	}
	for _, file := range filesThatNeedToBeRendered {
		if strings.HasSuffix(file.sourcePath, ".m4a") {
			// TODO: Extract to transcodeFileAtPath with sourcePath and destinationDir
			// Make variable from part of path that is subfolders of sourceDir
			trimmedPath := strings.TrimPrefix(file.sourcePath, sourceDir)
			fmt.Printf("🎟️ trimmedPath: %s\n", trimmedPath)
			sourcePath := filepath.Join(sourceDir, file.sourcePath)
			// TODO: Needs to be some part of destination with .mp3 extension (may need to rewrite compareDirectories() and related)
			destinationPath := filepath.Join(destinationDir, strings.TrimSuffix(trimmedPath, ".m4a")+".mp3")

			// destinationPathMP3 := strings.TrimSuffix(destinationPath, filepath.Ext(destinationPath)) + ".mp3"

			fmt.Printf("🔊 Transcoding: source %s destination %s\n", sourcePath, destinationPath)
			err := transcodeFileAtPath(sourcePath, destinationPath)
			if err != nil {
				fmt.Println("Error:", err)
			}
		}

	}

}

// transcodeFileAtPath transcodes the file at the specified path from .m4a to .mp3 format.
func transcodeFileAtPath(sourcePath, destinationPath string) error {
	trans := new(transcoder.Transcoder)

	if err := os.MkdirAll(filepath.Dir(destinationPath), 0755); err != nil {
		return fmt.Errorf("❗️Failed to create directories: %v", err)
	}

	err := trans.Initialize(sourcePath, destinationPath)
	if err != nil {
		return err
	}

	done := trans.Run(false)
	err = <-done
	if err != nil {
		return err
	}

	fmt.Printf("🎶 Transcoded: %s to %s\n", sourcePath, destinationPath)
	return nil
}

// compareDirectories compares the files in two directories and returns a list of the files exclusive to directory A.
func compareDirectories(a string, b string) ([]FileToRender, error) {
	filesA, err := getFilenames(a)
	if err != nil {
		return nil, err
	}

	filesB, err := getFilenames(b)
	if err != nil {
		return nil, err
	}

	exclusiveFiles := getExclusiveFiles(filesA, filesB)

	return exclusiveFiles, nil
}

// getFilenames returns a list of filenames in the specified directory.
func getFilenames(directory string) ([]string, error) {
	var filenames []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relativePath := strings.TrimPrefix(path, directory)
			filenames = append(filenames, relativePath)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return filenames, nil
}

// getExclusiveFiles returns the files exclusive to filesA compared to filesB.
func getExclusiveFiles(filesA, filesB []string) []FileToRender {
	exclusiveFiles := make([]FileToRender, 0)

	fileMap := make(map[string]bool)
	for _, file := range filesB {
		fileMap[file] = true
	}

	// Generate destination filenames so they can be compared to rendered output filenames
	var sourceFileOutputNameList []FileToRender
	for _, file := range filesA {
		// TODO: This needs to be a struct with source and destination filenames (so that they can be rendered properly)
		destinationFilename := ""
		if strings.HasSuffix(file, ".mp3") {
			destinationFilename = file
		} else {
			destinationFilename = strings.TrimSuffix(file, ".m4a") + ".mp3"
		}
		fileToRender := FileToRender{
			sourcePath:      file,
			destinationPath: destinationFilename,
		}

		sourceFileOutputNameList = append(sourceFileOutputNameList, fileToRender)
	}

	for _, file := range sourceFileOutputNameList {
		if !fileMap[file.destinationPath] {
			exclusiveFiles = append(exclusiveFiles, file)
		}
	}

	return exclusiveFiles
}
