package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xfrr/goffmpeg/transcoder"
)

// findFiles traverses the specified directory and transcodes all .m4a files to .mp3 format.
func findFiles(directory string) {
	fmt.Printf("🔨 Transcoding for directory %s\n", directory)
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
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

// transcodeFileAtPath transcodes the file at the specified path from .m4a to .mp3 format.
func transcodeFileAtPath(path string) error {
	trans := new(transcoder.Transcoder)
	output := strings.TrimSuffix(path, filepath.Ext(path)) + ".mp3"
	err := trans.Initialize(path, output)
	if err != nil {
		return err
	}

	done := trans.Run(false)
	err = <-done
	if err != nil {
		return err
	}

	fmt.Printf("🎶 Transcoded: %s to %s\n", path, output)
	return nil
}

// compareDirectories compares the files in two directories and prints the files exclusive to directory A.
func compareDirectories(a string, b string) error {
	filesA, err := getFilenames(a)
	if err != nil {
		return err
	}

	filesB, err := getFilenames(b)
	if err != nil {
		return err
	}

	exclusiveFiles := getExclusiveFiles(filesA, filesB)

	fmt.Println("Files exclusive to directory A:")
	for _, file := range exclusiveFiles {
		fmt.Println(file)
	}

	return nil
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

	for _, file := range filesA {
		if !fileMap[file] {
			exclusiveFiles = append(exclusiveFiles, file)
		}
	}

	return exclusiveFiles
}
