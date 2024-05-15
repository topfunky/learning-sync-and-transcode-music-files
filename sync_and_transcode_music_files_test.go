package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"log"
	"os/exec"

	"time"

	"github.com/stretchr/testify/assert"
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
	if err := os.WriteFile(path, []byte{}, 0644); err != nil {
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

	// Create some test files with .m4a extension
	testFiles := []string{
		"source/file1.m4a",
		"source/file2.m4a",
		"source/file4.m4a",
		"source/a-band/file5.m4a",
		"source/Whitespace Band/file6.m4a",
		"source/the-band/file7.mp3",
		"source/file8.aif",
		"source/file9.wav",
		"source/.DS_Store",
	}
	for _, file := range testFiles[0:numberOfFiles] {
		filePath := filepath.Join(tempDir, file)
		if err := generateM4aFixtureFileAtPath(filePath); err != nil {
			return fmt.Errorf("Failed to create test file: %v", err)
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
	transcodedFiles := []string{
		"destination/file1.mp3",
		"destination/file2.mp3",
		"destination/file4.mp3",
		"destination/a-band/file5.mp3",
		"destination/Whitespace Band/file6.mp3",
		"destination/the-band/file7.mp3",
		"destination/file8.mp3",
		"destination/file9.mp3",
		// NOTE: Do not list .DS_Store or .txt files since they should not be transcoded
	}

	tempDir, err := setup(t, len(transcodedFiles))
	if err != nil {
		t.Fatalf("failed to set up fixture files: %v", err)
	}

	defer os.RemoveAll(tempDir)

	findFiles(filepath.Join(tempDir, "source"), filepath.Join(tempDir, "destination"))

	for _, file := range transcodedFiles {
		t.Run(fmt.Sprintf("File %s should be rendered", file), func(t *testing.T) {
			filePath := filepath.Join(tempDir, file)
			assert.FileExistsf(t, filePath, "Transcoded file not found: %s", file)
		})
	}

	t.Run("Verify that the non-.m4a file was not transcoded", func(t *testing.T) {
		nonTranscodedFile := "file3.txt.transcoded"
		filePath := filepath.Join(tempDir, nonTranscodedFile)
		assert.NoFileExistsf(t, filePath, "unexpected transcoded file found: %s", nonTranscodedFile)
	})

	t.Run("Verify that the .DS_Store file was not transcoded", func(t *testing.T) {
		nonTranscodedFile := ".DS_Store"
		filePath := filepath.Join(tempDir, nonTranscodedFile)
		assert.NoFileExistsf(t, filePath, "unexpected transcoded file found: %s", nonTranscodedFile)
	})
}

func TestFindFiles_EmptyDestinationDirectory(t *testing.T) {
	transcodedFiles := []string{}

	tempDir, err := setup(t, len(transcodedFiles))
	defer os.RemoveAll(tempDir)

	if err != nil {
		t.Fatalf("❗️ Failed to create temporary directory: %v", err)
	}

	sourceDir := filepath.Join(tempDir, "source")
	destinationDir := filepath.Join(tempDir, "destination dir that does not exist")

	err = findFiles(sourceDir, destinationDir)
	assert.NoError(t, err)

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

		info1, _ := os.Stat(destinationPath)
		assert.FileExistsf(t, destinationPath, "Transcoded file not found: %s", file)

		// Wait for a second to ensure the modified time is different
		time.Sleep(time.Second)

		findFiles(sourceDir, destinationDir)

		info2, _ := os.Stat(destinationPath)
		assert.FileExistsf(t, destinationPath, "Transcoded file not found: %s", file)

		assert.Equal(t, info1.ModTime(), info2.ModTime(), fmt.Sprintf("file %s was re-rendered", destinationPath))
	})
}
