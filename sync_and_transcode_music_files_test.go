package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		"file3.txt", // Not an .m4a file
		"file4.m4a",
	}
	for _, file := range testFiles {
		filePath := filepath.Join(tempDir, file)
		if err := ioutil.WriteFile(filePath, []byte{}, 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
	}

	// Call the findFiles function
	findFiles(tempDir)

	// Verify that the transcoding was successful for .m4a files
	transcodedFiles := []string{
		"file1.m4a.transcoded",
		"file2.m4a.transcoded",
		"file4.m4a.transcoded",
	}
	for _, file := range transcodedFiles {
		filePath := filepath.Join(tempDir, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			assert.Fail(t, fmt.Sprintf("transcoded file not found: %s", file))
		}
	}

	// Verify that the non-.m4a file was not transcoded
	nonTranscodedFile := "file3.txt.transcoded"
	filePath := filepath.Join(tempDir, nonTranscodedFile)
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Errorf("unexpected transcoded file found: %s", nonTranscodedFile)
	}
}
