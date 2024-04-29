package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xfrr/goffmpeg/transcoder"
)

type FileToRender struct {
	sourcePath      string
	destinationPath string
}

// findFiles traverses the specified directory and transcodes all .m4a files to .mp3 format.
func findFiles(sourceDir, destinationDir string) {
	fmt.Printf("🔍 Finding files in source directory %s\n", sourceDir)
	// TODO: compareDirectories and transcodeFileAtPath with resulting list

	// TODO: Current output of compareDirectories() is destination name (.mp3) not source name (.m4a).
	// TODO: Needs to either look for existence of .m4a or compareDirectories() should be rewritten to return source file name
	// TODO: But if .mp3 exists as source, then it should be copied to the destination as-is
	// filesThatNeedToBeRendered := compareDirectories(sourceDir, destinationDir)
	// for _, file := range filesThatNeedToBeRendered {
	// 	// TODO: Transcode
	// }

	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), ".m4a") {
			// Make variable from part of path that is subfolders of sourceDir
			trimmedPath := strings.TrimPrefix(path, sourceDir)

			err := transcodeFileAtPath(path, filepath.Join(destinationDir, trimmedPath))
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

// transcodeFileAtPath transcodes the file at the specified path from .m4a to .mp3 format.
func transcodeFileAtPath(sourcePath, destinationPath string) error {
	trans := new(transcoder.Transcoder)

	destinationPathMP3 := strings.TrimSuffix(destinationPath, filepath.Ext(destinationPath)) + ".mp3"

	if err := os.MkdirAll(filepath.Dir(destinationPath), 0755); err != nil {
		return fmt.Errorf("failed to create directories: %v", err)
	}

	err := trans.Initialize(sourcePath, destinationPathMP3)
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
