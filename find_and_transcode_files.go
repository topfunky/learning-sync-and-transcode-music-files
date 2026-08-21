package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/xfrr/goffmpeg/transcoder"
)

type fileToTranscode struct {
	sourcePath      string
	destinationPath string
}

// findAndTranscodeFiles traverses the specified directory and transcodes music files to .mp3 format.
// MP3 files will be copied to the destination directory as-is.
// When dryRun is true, no files are written; it only prints what would be done.
// fileError pairs a file with the error encountered while processing it.
type fileError struct {
	file fileToTranscode
	err  error
}

// transcodeResult carries the outcome of one worker's transcode job.
type transcodeResult struct {
	file fileToTranscode
	err  error
}

// findAndTranscodeFiles traverses the specified directory and transcodes music files to .mp3 format.
// MP3 files will be copied to the destination directory as-is.
// When dryRun is true, no files are written; it only prints what would be done.
// Transcoding runs in parallel across all CPU cores; MP3 copies stay sequential.
// Per-file errors do not stop processing; a summary is printed to stderr at the end.
func findAndTranscodeFiles(sourceDir, destinationDir string, dryRun bool) error {
	fmt.Printf("🔍 Finding files in source directory %s\n", sourceDir)

	if !dryRun {
		if err := os.MkdirAll(destinationDir, 0755); err != nil {
			return fmt.Errorf("failed to create destination directory: %v", err)
		}
	}

	filesThatNeedToBeTranscoded, err := compareDirectories(sourceDir, destinationDir)
	if err != nil {
		if !(dryRun && os.IsNotExist(err)) {
			return fmt.Errorf("error: %v", err)
		}
		filesThatNeedToBeTranscoded, err = compareDirectories(sourceDir, sourceDir)
		if err != nil {
			return fmt.Errorf("error: %v", err)
		}
	}

	// Partition into files to transcode and MP3s to copy verbatim.
	var toTranscode []fileToTranscode
	var toCopy []fileToTranscode
	for _, file := range filesThatNeedToBeTranscoded {
		if isUntranscodedMusicFile(file.sourcePath) {
			toTranscode = append(toTranscode, file)
		} else {
			toCopy = append(toCopy, file)
		}
	}

	var errs []fileError

	if dryRun {
		for _, file := range toTranscode {
			sourcePath := filepath.Join(sourceDir, file.sourcePath)
			destinationPath := filepath.Join(destinationDir, convertSourceToDestinationFilename(file.sourcePath))
			fmt.Printf("🔍 [dry-run] Would transcode: %s ➡️  %s\n", sourcePath, destinationPath)
		}
		for _, file := range toCopy {
			sourcePath := filepath.Join(sourceDir, file.sourcePath)
			destinationPath := filepath.Join(destinationDir, file.destinationPath)
			fmt.Printf("🔍 [dry-run] Would copy MP3: %s ➡️  %s\n", sourcePath, destinationPath)
		}
		return nil
	}

	// Copy MP3s sequentially; they are fast and I/O-bound.
	for i, file := range toCopy {
		sourcePath := filepath.Join(sourceDir, file.sourcePath)
		destinationPath := filepath.Join(destinationDir, file.destinationPath)
		if err := copyFile(sourcePath, destinationPath); err != nil {
			fmt.Fprintf(os.Stderr, "[%d/%d] ❗️ Error while copying file %s: %v\n", i+1, len(toCopy), sourcePath, err)
			errs = append(errs, fileError{file: file, err: err})
			continue
		}
		fmt.Printf("[%d/%d] 📂 Copied MP3: %s\n", i+1, len(toCopy), destinationPath)
	}

	// Transcode concurrently with a worker pool sized to the CPU count.
	if len(toTranscode) > 0 {
		workers := runtime.NumCPU()
		if workers > len(toTranscode) {
			workers = len(toTranscode)
		}

		jobs := make(chan fileToTranscode)
		results := make(chan transcodeResult, len(toTranscode))

		var wg sync.WaitGroup
		for w := 0; w < workers; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for file := range jobs {
					sourcePath := filepath.Join(sourceDir, file.sourcePath)
					err := transcodeFileAtPath(file.sourcePath, sourcePath, destinationDir)
					results <- transcodeResult{file: file, err: err}
				}
			}()
		}

		go func() {
			for _, file := range toTranscode {
				jobs <- file
			}
			close(jobs)
		}()

		go func() {
			wg.Wait()
			close(results)
		}()

		completed := 0
		for result := range results {
			completed++
			sourcePath := filepath.Join(sourceDir, result.file.sourcePath)
			if result.err != nil {
				fmt.Fprintf(os.Stderr, "[%d/%d] ❗️ Error while transcoding %s: %v\n", completed, len(toTranscode), sourcePath, result.err)
				errs = append(errs, fileError(result))
				continue
			}
			destinationPath := filepath.Join(destinationDir, result.file.destinationPath)
			fmt.Printf("[%d/%d] 🔊 Transcoded: %s ➡️  %s\n", completed, len(toTranscode), sourcePath, destinationPath)
		}
	}

	if len(errs) > 0 {
		fmt.Fprintf(os.Stderr, "❗️ %d of %d files failed:\n", len(errs), len(toTranscode)+len(toCopy))
		for _, fe := range errs {
			fmt.Fprintf(os.Stderr, "  - %s: %v\n", fe.file.sourcePath, fe.err)
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

	return nil
}

// compareDirectories compares the files in two directories and returns a list of the files exclusive to directory A.
// The return value is the files that need to be transcoded (or copied to the destination, if already MP3).
func compareDirectories(a string, b string) ([]fileToTranscode, error) {
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
func getExclusiveFiles(filesA, filesB []string) []fileToTranscode {
	exclusiveFiles := make([]fileToTranscode, 0)

	fileMap := make(map[string]bool)
	for _, file := range filesB {
		fileMap[file] = true
	}

	// Generate list of filenames that need to be transcoded later
	var sourceFileOutputNameList []fileToTranscode
	for _, file := range filesA {
		destinationFilename := ""
		if strings.HasPrefix(filepath.Base(file), "._") {
			// Skip hidden files
			continue
		} else if strings.HasSuffix(file, ".mp3") {
			// Normalize and strip redundant artist/album segments from the .mp3 file name so it can be copied later
			destinationFilename = stripArtistAlbumFromFilename(removeNonASCII(file))
		} else if isUntranscodedMusicFile(file) {
			// Add file to struct so it can be transcoded to .mp3 later
			destinationFilename = convertSourceToDestinationFilename(file)
		} else {
			// Ignore .DS_Store, .txt and other files
			file = ""
		}
		fileToTranscode := fileToTranscode{
			sourcePath:      file,
			destinationPath: destinationFilename,
		}

		sourceFileOutputNameList = append(sourceFileOutputNameList, fileToTranscode)
	}

	for _, file := range sourceFileOutputNameList {
		if !fileMap[file.destinationPath] && file.destinationPath != "" {
			exclusiveFiles = append(exclusiveFiles, file)
		}
	}

	return exclusiveFiles
}

// convertSourceToDestinationFilename converts the filename by replacing the .m4a suffix with .mp3, replacing non-ASCII characters with an ASCII equivalent, and stripping redundant artist/album name segments that repeat a containing directory name.
func convertSourceToDestinationFilename(filename string) string {
	// Replace .m4a suffix with .mp3
	filename = strings.TrimSuffix(filename, filepath.Ext(filename)) + ".mp3"

	// Replace non-ASCII characters with an ASCII equivalent
	filename = removeNonASCII(filename)

	// Remove filename segments that redundantly repeat the artist/album directory names
	filename = stripArtistAlbumFromFilename(filename)

	return filename
}
