package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/xfrr/goffmpeg/transcoder"
)

// TODO: Fields should be titlecase if meant to be public (probably not?)
type FileToRender struct {
	sourcePath      string
	destinationPath string
}

// findFiles traverses the specified directory and transcodes all .m4a files to .mp3 format.
func findFiles(sourceDir, destinationDir string) error {
	fmt.Printf("🔍 Finding files in source directory %s\n", sourceDir)

	if err := os.MkdirAll(destinationDir, 0755); err != nil {
		return fmt.Errorf("Failed to create destination directory: %v", err)
	}

	filesThatNeedToBeRendered, err := compareDirectories(sourceDir, destinationDir)
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	for _, file := range filesThatNeedToBeRendered {
		sourcePath := filepath.Join(sourceDir, file.sourcePath)

		if strings.HasSuffix(sourcePath, ".m4a") || strings.HasSuffix(sourcePath, ".aif") {
			// TODO: Extract to transcodeFileAtPath with sourcePath and destinationDir
			destinationPath := filepath.Join(destinationDir, strings.TrimSuffix(file.sourcePath, filepath.Ext(file.sourcePath))+".mp3")
			err := transcodeFileAtPath(sourcePath, destinationPath)
			if err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			// Copy mp3 from source to destination
			destinationPath := filepath.Join(destinationDir, file.sourcePath)
			fmt.Printf("📂 Copy MP3: %s\n", destinationPath)
			if err := copyFile(sourcePath, destinationPath); err != nil {
				fmt.Println("Error:", err)
			}
		}
	}
	return nil
}

func copyFile(source, destination string) error {
	if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
		return fmt.Errorf("❗️Failed to create directories: %v", err)
	}

	// Open the source file for reading
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create the destination file
	destinationFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// Copy the contents of the source file into the destination file
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	// Call Sync to flush writes to stable storage
	destinationFile.Sync()

	return nil
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

	fmt.Printf("🔊 Transcoded: %s to %s\n", sourcePath, destinationPath)
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

		destinationFilename := ""
		if strings.HasSuffix(file, ".mp3") {
			// Copy .mp3 files over verbatim
			destinationFilename = file
		} else if strings.HasSuffix(file, ".m4a") || strings.HasSuffix(file, ".aif") {
			// Other files need to be transcoded to .mp3
			destinationFilename = strings.TrimSuffix(file, filepath.Ext(file)) + ".mp3"
		} else {
			// .DS_Store and other files should be ignored
			file = ""
		}
		fileToRender := FileToRender{
			sourcePath:      file,
			destinationPath: destinationFilename,
		}

		sourceFileOutputNameList = append(sourceFileOutputNameList, fileToRender)
	}

	for _, file := range sourceFileOutputNameList {
		if !fileMap[file.destinationPath] && file.destinationPath != "" {
			exclusiveFiles = append(exclusiveFiles, file)
		}
	}

	return exclusiveFiles
}
