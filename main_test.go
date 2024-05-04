package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupMainTest(t *testing.T) (string, error) {
	// Create a temporary directory for testing
	tempDir, err := ioutil.TempDir("", "test")
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
