package main

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupMainTest(t *testing.T) (string, error) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	return tempDir, nil
}

func TestMainFunction(t *testing.T) {
	// Save command-line arguments
	oldArgs := os.Args
	tempDir, _ := setupMainTest(t)
	defer func() {
		// Restore command-line arguments
		os.Args = oldArgs
		os.RemoveAll(tempDir)
	}()

	sourceDir := filepath.Join(tempDir, "source")
	_ = os.MkdirAll(sourceDir, 0755)
	destinationDir := filepath.Join(tempDir, "destination")

	// Set up command-line arguments for testing
	os.Args = []string{"cmd", "-source=" + sourceDir, "-destination=" + destinationDir}

	// Reset flag.CommandLine to clear flags defined in other tests
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Call main() function
	main()

	// TODO: Write more meaningful assertions to verify that the application works at a high level
	assert.DirExists(t, sourceDir)
	assert.DirExists(t, destinationDir)
}

func TestMainFunction_DryRunFlag(t *testing.T) {
	oldArgs := os.Args
	tempDir, _ := setupMainTest(t)
	defer func() {
		os.Args = oldArgs
		os.RemoveAll(tempDir)
	}()

	sourceDir := filepath.Join(tempDir, "source")
	_ = os.MkdirAll(sourceDir, 0755)
	destinationDir := filepath.Join(tempDir, "destination")
	_ = os.MkdirAll(destinationDir, 0755)

	// Place a duplicate pair directly in the destination directory.
	mp3File := filepath.Join(destinationDir, "song.mp3")
	m4aFile := filepath.Join(destinationDir, "song.m4a")
	os.WriteFile(mp3File, make([]byte, 100), 0644)
	os.WriteFile(m4aFile, make([]byte, 200), 0644)

	os.Args = []string{"cmd", "-source=" + sourceDir, "-destination=" + destinationDir, "-dry-run"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	main()

	// Dry run must not delete anything
	assert.FileExists(t, mp3File)
	assert.FileExists(t, m4aFile)
}
