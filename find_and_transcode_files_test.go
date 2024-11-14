package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"log"

	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetExclusiveFiles(t *testing.T) {
	cases := []struct {
		Name            string
		SourceList      []string
		DestinationList []string
		ExpectedOutput  []string
	}{
		{
			Name:            "Both file lists are empty",
			SourceList:      []string{},
			DestinationList: []string{},
			ExpectedOutput:  []string(nil),
		},
		{
			Name:            "Empty source list, destination list has elements",
			SourceList:      []string{},
			DestinationList: []string{"file1.mp3", "file2.mp3"},
			ExpectedOutput:  []string(nil),
		},
		{
			Name:            "Source list has elements, empty destination list",
			SourceList:      []string{"file1.m4a", "file2.m4a"},
			DestinationList: []string{},
			ExpectedOutput:  []string{"file1.mp3", "file2.mp3"},
		},
		{
			Name:            "Source and destination have common elements",
			SourceList:      []string{"file1.m4a", "file2.m4a", "file3.m4a"},
			DestinationList: []string{"file2.mp3", "file3.mp3", "file4.mp3"},
			ExpectedOutput:  []string{"file1.mp3"},
		},
		{
			Name:            "Source and destination have disjoint elements",
			SourceList:      []string{"file1.m4a", "file2.m4a", "file3.m4a"},
			DestinationList: []string{"file4.mp3", "file5.mp3", "file6.mp3"},
			ExpectedOutput:  []string{"file1.mp3", "file2.mp3", "file3.mp3"},
		},
		{
			Name:            "Source contains mp3 files which should be copied verbatim",
			SourceList:      []string{"file1.m4a", "file2.m4a", "file3.m4a", "file103.mp3"},
			DestinationList: []string{"file4.mp3", "file5.mp3", "file6.mp3"},
			ExpectedOutput:  []string{"file1.mp3", "file2.mp3", "file3.mp3", "file103.mp3"},
		},
		{
			Name:            "Destination contains m4a files which should be ignored",
			SourceList:      []string{"file1.m4a", "file2.m4a", "file3.m4a"},
			DestinationList: []string{"file4.m4a", "file5.mp3", "file6.mp3"},
			ExpectedOutput:  []string{"file1.mp3", "file2.mp3", "file3.mp3"},
		},
		{
			Name:            "Destination contains aif and wav files which should be transcoded",
			SourceList:      []string{"file1.m4a", "file2.aif", "file3.wav"},
			DestinationList: []string{},
			ExpectedOutput:  []string{"file1.mp3", "file2.mp3", "file3.mp3"},
		},
		{
			Name:            "Ignore non-music files",
			SourceList:      []string{".DS_Store"},
			DestinationList: []string{},
			ExpectedOutput:  []string(nil),
		},
		{
			Name:            "Ignore dotfiles",
			SourceList:      []string{"._file7.m4a"},
			DestinationList: []string{},
			ExpectedOutput:  []string(nil),
		},
		{
			Name:            "Correctly compares non-ASCII filenames",
			SourceList:      []string{"Alexandra Stréliski/Néo-Romance (Extended Version) [96kHz · 24bit]/02 - Lumières.m4a"},
			DestinationList: []string{},
			ExpectedOutput:  []string{"Alexandra Streliski/Neo-Romance (Extended Version) [96kHz  24bit]/02 - Lumieres.mp3"},
		},
		{
			Name:            "Correctly compares non-ASCII filenames (alt)",
			SourceList:      []string{"Megan Perry Fisher/Megan Perry Fisher - Pensées/Megan Perry Fisher - Pensées - 12 Pensée xii.m4a", "Stéphane Grappelli, Joe Pass & Niels-Henning Ørsted Pedersen/Tivoli Gardens, Copenhagen, Denmark (Live)/01 It's Only A Paper Moon.m4a"},
			DestinationList: []string{},
			ExpectedOutput:  []string{"Megan Perry Fisher/Megan Perry Fisher - Pensees/Megan Perry Fisher - Pensees - 12 Pensee xii.mp3", "Stephane Grappelli, Joe Pass & Niels-Henning Orsted Pedersen/Tivoli Gardens, Copenhagen, Denmark (Live)/01 It's Only A Paper Moon.mp3"},
		},
		{
			Name:            "Does not re-transcode non-ASCII filenames",
			SourceList:      []string{"Megan Perry Fisher/Megan Perry Fisher - Pensées/Megan Perry Fisher - Pensées - 12 Pensée xii.m4a", "Stéphane Grappelli, Joe Pass & Niels-Henning Ørsted Pedersen/Tivoli Gardens, Copenhagen, Denmark (Live)/01 It's Only A Paper Moon.m4a"},
			DestinationList: []string{"Megan Perry Fisher/Megan Perry Fisher - Pensees/Megan Perry Fisher - Pensees - 12 Pensee xii.mp3", "Stephane Grappelli, Joe Pass & Niels-Henning Orsted Pedersen/Tivoli Gardens, Copenhagen, Denmark (Live)/01 It's Only A Paper Moon.mp3"},
			ExpectedOutput:  []string(nil),
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			result := getExclusiveFiles(c.SourceList, c.DestinationList)
			assert.Equal(t, c.ExpectedOutput, getDestinationPaths(result))
		})
	}
}

func TestConvertSourceToDestinationFilename(t *testing.T) {
	cases := []struct {
		Name           string
		Filename       string
		ExpectedOutput string
	}{
		{
			Name:           "Filename with .m4a extension",
			Filename:       "file1.m4a",
			ExpectedOutput: "file1.mp3",
		},
		{
			Name:           "Filename with .m4a extension and non-ASCII characters",
			Filename:       "Megan Perry Fisher - Pensées.m4a",
			ExpectedOutput: "Megan Perry Fisher - Pensees.mp3",
		},
		{
			Name:           "Filename with .mp3 extension",
			Filename:       "file2.mp3",
			ExpectedOutput: "file2.mp3",
		},
		{
			Name:           "Filename with no extension",
			Filename:       "file3",
			ExpectedOutput: "file3.mp3",
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			result := convertSourceToDestinationFilename(c.Filename)
			assert.Equal(t, c.ExpectedOutput, result)
		})
	}
}

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

	// Create test files
	testFiles := []string{
		"source/file1.m4a",
		"source/file2.m4a",
		"source/Alexandra Stréliski/Néo-Romance (Extended Version) [96kHz · 24bit]/02 - Lumières.m4a",
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
	tempDir, err := os.MkdirTemp("", "test")
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
		"destination/Alexandra Streliski/Neo-Romance (Extended Version) [96kHz  24bit]/02 - Lumieres.mp3",
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

	findAndTranscodeFiles(filepath.Join(tempDir, "source"), filepath.Join(tempDir, "destination"))

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

	err = findAndTranscodeFiles(sourceDir, destinationDir)
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
	findAndTranscodeFiles(sourceDir, destinationDir)

	// Verify that the destination files were not re-rendered
	file := "source/file1.m4a"
	t.Run(fmt.Sprintf("File %s should not be re-rendered", file), func(t *testing.T) {
		destinationPath := filepath.Join(tempDir, "destination/file1.mp3")

		info1, _ := os.Stat(destinationPath)
		assert.FileExistsf(t, destinationPath, "Transcoded file not found: %s", file)

		// Wait for a second to ensure the modified time is different
		time.Sleep(time.Second)

		findAndTranscodeFiles(sourceDir, destinationDir)

		info2, _ := os.Stat(destinationPath)
		assert.FileExistsf(t, destinationPath, "Transcoded file not found: %s", file)

		assert.Equal(t, info1.ModTime(), info2.ModTime(), fmt.Sprintf("file %s was re-rendered", destinationPath))
	})
}

// Returns a string array of only the `sourcePath` attribute from an array of `fileToTranscode` structs.
//
// This makes test assertions cleaner, based on how the fixture data is written.
func getDestinationPaths(files []fileToTranscode) []string {
	var sources []string
	for _, file := range files {
		sources = append(sources, file.destinationPath)
	}
	return sources
}
