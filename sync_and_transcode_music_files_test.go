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
)

func generateM4aFixtureFileAtPath(path string) error {
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

func TestFindFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create some test files with .m4a extension
	testFiles := []string{
		"file1.m4a",
		"file2.m4a",
		"file4.m4a",
	}
	for _, file := range testFiles {
		filePath := filepath.Join(tempDir, file)
		if err := generateM4aFixtureFileAtPath(filePath); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
	}

	// A text file that is not an m4a file
	testTextFileName := "file3.txt" // Not an .m4a file
	textFilePath := filepath.Join(tempDir, testTextFileName)
	if err := generateTextFileFixtureAtPath(textFilePath); err != nil {
		t.Fatalf("failed to create test text file: %v", err)
	}

	// Call the findFiles function
	findFiles(tempDir)

	// Verify that the transcoding was successful for .m4a files
	transcodedFiles := []string{
		"file1.mp3",
		"file2.mp3",
		"file4.mp3",
	}
	for _, file := range transcodedFiles {
		filePath := filepath.Join(tempDir, file)
		_, err = os.Stat(filePath)
		assert.False(t, os.IsNotExist(err), fmt.Sprintf("transcoded file not found: %s", file))
	}

	// Verify that the non-.m4a file was not transcoded
	nonTranscodedFile := "file3.txt.transcoded"
	filePath := filepath.Join(tempDir, nonTranscodedFile)
	_, err = os.Stat(filePath)
	assert.True(t, os.IsNotExist(err), fmt.Sprintf("unexpected transcoded file found: %s", nonTranscodedFile))
}
