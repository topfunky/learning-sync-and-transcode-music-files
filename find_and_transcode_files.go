package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/xfrr/goffmpeg/transcoder"
)

type fileToRender struct {
	sourcePath      string
	destinationPath string
}

// findAndTranscodeFiles traverses the specified directory and transcodes music files to .mp3 format.
// MP3 files will be copied to the destination directory as-is.
func findAndTranscodeFiles(sourceDir, destinationDir string) error {
	fmt.Printf("🔍 Finding files in source directory %s\n", sourceDir)

	if err := os.MkdirAll(destinationDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %v", err)
	}

	filesThatNeedToBeRendered, err := compareDirectories(sourceDir, destinationDir)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	for _, file := range filesThatNeedToBeRendered {
		sourcePath := filepath.Join(sourceDir, file.sourcePath)

		if isUntranscodedMusicFile(sourcePath) {
			err := transcodeFileAtPath(file.sourcePath, sourcePath, destinationDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "❗️ Error while transcoding file: %v\n", err)
				// TODO: Maybe return error or queue for return
			}
		} else {
			// Copy mp3 from source to destination
			destinationPath := filepath.Join(destinationDir, file.sourcePath)
			if err := copyFile(sourcePath, destinationPath); err != nil {
				fmt.Fprintf(os.Stderr, "❗️ Error while copying file: %v\n", err)
				// TODO: Maybe return error or queue for return
			}
			fmt.Printf("📂 Copied MP3: %s\n", destinationPath)
		}
	}
	return nil
}

// copyFile copies a file from the source path to the destination path.
// It creates any necessary directories in the destination path.
// If the file cannot be copied for any reason, it returns an error.
//
// Example usage:
//
//	err := copyFile("/path/to/source", "/path/to/destination")
//	if err != nil {
//	    log.Fatal(err)
//	}
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

// transcodeFileAtPath transcodes the music file at the specified path to .mp3 format.
func transcodeFileAtPath(fileSourcePath, sourcePath, destinationDir string) error {
	// TODO: Rename fileSourcePath to a more descriptive name. It's a relative path and is used for source and destination subdirs (with filename)
	destinationPath := filepath.Join(destinationDir, convertSourceToDestinationFilename(fileSourcePath))

	if err := os.MkdirAll(filepath.Dir(destinationPath), 0755); err != nil {
		return fmt.Errorf("❗️Failed to create directories: %v", err)
	}

	trans := new(transcoder.Transcoder)
	if err := trans.Initialize(sourcePath, destinationPath); err != nil {
		return err
	}

	done := trans.Run(false)
	if err := <-done; err != nil {
		return err
	}

	fmt.Printf("🔊 Transcoded: %s ➡️  %s\n", sourcePath, destinationPath)
	return nil
}

// compareDirectories compares the files in two directories and returns a list of the files exclusive to directory A.
// The return value is the files that need to be transcoded (or copied to the destination, if already MP3).
func compareDirectories(a string, b string) ([]fileToRender, error) {
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
func getExclusiveFiles(filesA, filesB []string) []fileToRender {
	exclusiveFiles := make([]fileToRender, 0)

	fileMap := make(map[string]bool)
	for _, file := range filesB {
		fileMap[file] = true
	}

	// Generate list of filenames that need to be transcoded later
	var sourceFileOutputNameList []fileToRender
	for _, file := range filesA {
		destinationFilename := ""
		if strings.HasSuffix(file, ".mp3") {
			// Save .mp3 file name verbatim so it can be copied later
			destinationFilename = file
		} else if isUntranscodedMusicFile(file) {
			// Add file to struct so it can be transcoded to .mp3 later
			destinationFilename = convertSourceToDestinationFilename(file)
		} else {
			// Ignore .DS_Store, .txt and other files
			file = ""
		}
		fileToRender := fileToRender{
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

// convertSourceToDestinationFilename converts the filename by replacing the .m4a suffix with .mp3 and replacing non-ASCII characters with an ASCII equivalent.
func convertSourceToDestinationFilename(filename string) string {
	// Replace .m4a suffix with .mp3
	filename = strings.TrimSuffix(filename, filepath.Ext(filename)) + ".mp3"

	// Replace non-ASCII characters with an ASCII equivalent
	filename = removeNonASCII(filename)

	return filename
}
