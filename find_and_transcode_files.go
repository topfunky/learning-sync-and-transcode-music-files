package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xfrr/goffmpeg/transcoder"
)

// findFiles traverses the specified directory and transcodes all .m4a files to .mp3 format.
func findFiles(sourceDir, destinationDir string) {
	fmt.Printf("🔍 Finding files in source directory %s\n", sourceDir)
	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), ".m4a") {
			err := transcodeFileAtPath(path)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path: %v\n", err)
		return
	}
}

func trimFirstFolder(path string) (string, error) {
	// Split the path into a slice of its components
	parts := strings.Split(filepath.ToSlash(path), "/")

	// Check if the path has at least one folder
	if len(parts) < 2 {
		return "", fmt.Errorf("path does not contain a folder")
	}

	// Remove the first folder
	trimmedPath := strings.Join(parts[1:], "/")

	return trimmedPath, nil
}

// transcodeFileAtPath transcodes the file at the specified path from .m4a to .mp3 format.
// TODO: Take `sourcePath` and `destinationPath` as args.
// TODO: Render to destinationPath at subfolder that matches sourcePath
func transcodeFileAtPath(sourcePath string) error {
	trans := new(transcoder.Transcoder)
	// TODO: Call trimFirstFolder(sourcePath)
	output := strings.TrimSuffix(sourcePath, filepath.Ext(sourcePath)) + ".mp3"
	err := trans.Initialize(sourcePath, output)
	if err != nil {
		return err
	}

	done := trans.Run(false)
	err = <-done
	if err != nil {
		return err
	}

	fmt.Printf("🎶 Transcoded: %s to %s\n", sourcePath, output)
	return nil
}

// compareDirectories compares the files in two directories and returns a list of the files exclusive to directory A.
func compareDirectories(a string, b string) ([]string, error) {
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
			filenames = append(filenames, info.Name())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return filenames, nil
}

// getExclusiveFiles returns the files exclusive to filesA compared to filesB.
func getExclusiveFiles(filesA, filesB []string) []string {
	exclusiveFiles := make([]string, 0)

	fileMap := make(map[string]bool)
	for _, file := range filesB {
		fileMap[file] = true
	}

	// Generate destination filenames so they can be compared to rendered output filenames
	var sourceFileOutputNameList []string
	for _, file := range filesA {
		destinationFilename := ""
		if strings.HasSuffix(file, ".mp3") {
			destinationFilename = file
		} else {
			destinationFilename = strings.TrimSuffix(file, ".m4a") + ".mp3"
		}
		sourceFileOutputNameList = append(sourceFileOutputNameList, destinationFilename)
	}

	for _, file := range sourceFileOutputNameList {
		if !fileMap[file] {
			exclusiveFiles = append(exclusiveFiles, file)
		}
	}

	return exclusiveFiles
}
