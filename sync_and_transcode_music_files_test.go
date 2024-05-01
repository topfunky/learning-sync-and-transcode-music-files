package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"log"
	"os/exec"

	"github.com/stretchr/testify/assert"
	"time"
)

func generateM4aFixtureFileAtPath(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create directories: %v", err)
	}
	cmd := exec.Command("ffmpeg", "-f", "lavfi", "-i", "sine=frequency=1000:duration=5", path)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func generateTextFileFixtureAtPath(path string) error {
	if err := ioutil.WriteFile(path, []byte{}, 0644); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func setupFixtureFilesInDirectory(tempDir string, numberOfFiles int) error {
	// Create a directory within tempDir named "source"
	sourceDir := filepath.Join(tempDir, "source")
	if err := os.Mkdir(sourceDir, 0755); err != nil {
		return fmt.Errorf("failed to create source directory: %v", err)
	}

	destinationDir := filepath.Join(tempDir, "destination")
	if err := os.Mkdir(destinationDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %v", err)
	}

	// Create some test files with .m4a extension
	testFiles := []string{
		"source/file1.m4a",
		"source/file2.m4a",
		"source/file4.m4a",
		"source/a-band/file5.m4a",
		"source/Whitespace Band/file6.m4a",
	}
	for _, file := range testFiles[0:numberOfFiles] {
		filePath := filepath.Join(tempDir, file)
		if err := generateM4aFixtureFileAtPath(filePath); err != nil {
			return fmt.Errorf("failed to create test file: %v", err)
		}
	}

	// A text file that is not an m4a file
	testTextFileName := "file3.txt" // Not an .m4a file
	textFilePath := filepath.Join(tempDir, testTextFileName)
	if err := generateTextFileFixtureAtPath(textFilePath); err != nil {
		return fmt.Errorf("failed to create test text file: %v", err)
	}

	return nil
}

func setup(t *testing.T, numberOfFiles int) (string, error) {
	// Create a temporary directory for testing
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}

	// Set up fixture files in the temporary directory
	if err := setupFixtureFilesInDirectory(tempDir, numberOfFiles); err != nil {
		t.Fatalf("failed to set up fixture files: %v", err)
	}
	return tempDir, nil
}

func TestFindFiles(t *testing.T) {
	tempDir, err := setup(t, 5)
	defer os.RemoveAll(tempDir)

	findFiles(filepath.Join(tempDir, "source"), filepath.Join(tempDir, "destination"))

	transcodedFiles := []string{
		"destination/file1.mp3",
		"destination/file2.mp3",
		"destination/file4.mp3",
		"destination/a-band/file5.mp3",
		"destination/Whitespace Band/file6.mp3",
	}
	for _, file := range transcodedFiles {
		t.Run(fmt.Sprintf("File %s should be rendered", file), func(t *testing.T) {
			filePath := filepath.Join(tempDir, file)
			_, err = os.Stat(filePath)
			assert.False(t, os.IsNotExist(err), fmt.Sprintf("transcoded file not found: %s", file))
		})
	}

	t.Run("Verify that the non-.m4a file was not transcoded", func(t *testing.T) {
		nonTranscodedFile := "file3.txt.transcoded"
		filePath := filepath.Join(tempDir, nonTranscodedFile)
		_, err = os.Stat(filePath)
		assert.True(t, os.IsNotExist(err), fmt.Sprintf("unexpected transcoded file found: %s", nonTranscodedFile))
	})

}

// Destination files should not be re-rendered (check file modified time from first render and compare to second render)
func TestFindFiles_NoReRender(t *testing.T) {
	// Generate limited test fixtures with one media file.
	tempDir, _ := setup(t, 1)
	defer os.RemoveAll(tempDir)

	sourceDir := filepath.Join(tempDir, "source")
	destinationDir := filepath.Join(tempDir, "destination")

	// Run the function for the first time
	findFiles(sourceDir, destinationDir)

	// Verify that the destination files were not re-rendered
	file := "source/file1.m4a"
	t.Run(fmt.Sprintf("File %s should not be re-rendered", file), func(t *testing.T) {
		destinationPath := filepath.Join(tempDir, "destination/file1.mp3")

		info1, err := os.Stat(destinationPath)
		assert.False(t, os.IsNotExist(err), fmt.Sprintf("transcoded file not found: %s", file))

		// Wait for a second to ensure the modified time is different
		time.Sleep(time.Second)

		findFiles(sourceDir, destinationDir)

		info2, err := os.Stat(destinationPath)
		assert.False(t, os.IsNotExist(err), fmt.Sprintf("transcoded file not found: %s", file))

		assert.Equal(t, info1.ModTime(), info2.ModTime(), fmt.Sprintf("file %s was re-rendered", destinationPath))
	})
}
