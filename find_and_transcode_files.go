package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xfrr/goffmpeg/transcoder"
)

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
