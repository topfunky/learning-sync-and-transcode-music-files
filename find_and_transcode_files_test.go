package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

	findAndTranscodeFiles(filepath.Join(tempDir, "source"), filepath.Join(tempDir, "destination"), false)

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

	err = findAndTranscodeFiles(sourceDir, destinationDir, false)
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
	findAndTranscodeFiles(sourceDir, destinationDir, false)

	// Verify that the destination files were not re-rendered
	file := "source/file1.m4a"
	t.Run(fmt.Sprintf("File %s should not be re-rendered", file), func(t *testing.T) {
		destinationPath := filepath.Join(tempDir, "destination/file1.mp3")

		info1, _ := os.Stat(destinationPath)
		assert.FileExistsf(t, destinationPath, "Transcoded file not found: %s", file)

		// Wait for a second to ensure the modified time is different
		time.Sleep(time.Second)

		findAndTranscodeFiles(sourceDir, destinationDir, false)

		info2, _ := os.Stat(destinationPath)
		assert.FileExistsf(t, destinationPath, "Transcoded file not found: %s", file)

		assert.Equal(t, info1.ModTime(), info2.ModTime(), fmt.Sprintf("file %s was re-rendered", destinationPath))
	})
}

func TestCopyFile_SourceNotFound(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-copyfile")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	err = copyFile(filepath.Join(tempDir, "nonexistent.mp3"), filepath.Join(tempDir, "dest.mp3"))
	assert.Error(t, err)
}

func TestCopyFile_Success(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-copyfile-ok")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	src := filepath.Join(tempDir, "src.mp3")
	dst := filepath.Join(tempDir, "subdir", "dst.mp3")
	os.WriteFile(src, []byte("audio data"), 0644)

	err = copyFile(src, dst)
	assert.NoError(t, err)
	assert.FileExists(t, dst)
}

func TestCompareDirectories_InvalidSource(t *testing.T) {
	_, err := compareDirectories("/nonexistent/source/dir", "/tmp")
	assert.Error(t, err)
}

func TestCompareDirectories_InvalidDestination(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-compare")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	_, err = compareDirectories(tempDir, "/nonexistent/destination/dir")
	assert.Error(t, err)
}

func TestGetFilenames_NonExistentDirectory(t *testing.T) {
	_, err := getFilenames("/nonexistent/dir")
	assert.Error(t, err)
}

func TestGetFilenames_ReturnsRelativePaths(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-getfilenames")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	os.WriteFile(filepath.Join(tempDir, "a.mp3"), []byte{}, 0644)
	os.WriteFile(filepath.Join(tempDir, "b.mp3"), []byte{}, 0644)

	names, err := getFilenames(tempDir)
	assert.NoError(t, err)
	assert.Len(t, names, 2)
	for _, n := range names {
		assert.True(t, !filepath.IsAbs(n) || n[0] == '/', "path should be relative to dir or start with separator")
	}
}

// TestFindFiles_CopiesMP3OnlyFileVerbatim verifies that an mp3 file that exists
// in the source with no higher-quality counterpart is always copied to the
// destination verbatim, regardless of its file size (used as a bitrate proxy).
func TestFindFiles_CopiesMP3OnlyFileVerbatim(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-mp3-only")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sourceDir := filepath.Join(tempDir, "source")
	destinationDir := filepath.Join(tempDir, "destination")
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("failed to create source dir: %v", err)
	}

	// Write a small MP3 file to simulate a low-bitrate track that only exists as MP3.
	mp3Filename := "only-as-mp3.mp3"
	mp3Content := []byte("fake mp3 content representing a low-bitrate file")
	if err := os.WriteFile(filepath.Join(sourceDir, mp3Filename), mp3Content, 0644); err != nil {
		t.Fatalf("failed to create MP3 source file: %v", err)
	}

	err = findAndTranscodeFiles(sourceDir, destinationDir, false)
	assert.NoError(t, err)

	// The MP3 must be present in the destination regardless of its size.
	destFile := filepath.Join(destinationDir, mp3Filename)
	assert.FileExistsf(t, destFile, "mp3-only source file was not copied to destination")

	// The copy must be verbatim (byte-for-byte identical).
	destContent, err := os.ReadFile(destFile)
	assert.NoError(t, err)
	assert.Equal(t, mp3Content, destContent, "destination file content differs from source")
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

// captureStdout redirects os.Stdout for the duration of fn and returns what was printed.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w
	defer func() { os.Stdout = old }()

	done := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		done <- buf.String()
	}()

	fn()

	w.Close()
	os.Stdout = old
	return <-done
}

func TestFindFiles_DryRunDoesNotCreateDestinationOrWriteFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-dryrun")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sourceDir := filepath.Join(tempDir, "source")
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("failed to create source dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "track.wav"), []byte("fake wav"), 0644); err != nil {
		t.Fatalf("failed to write wav fixture: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "track.mp3"), []byte("fake mp3"), 0644); err != nil {
		t.Fatalf("failed to write mp3 fixture: %v", err)
	}

	destinationDir := filepath.Join(tempDir, "destination")

	err = findAndTranscodeFiles(sourceDir, destinationDir, true)
	assert.NoError(t, err)

	assert.NoDirExists(t, destinationDir, "dry-run must not create the destination directory")
}

func TestFindFiles_DryRunPrintsPlannedActions(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-dryrun-output")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sourceDir := filepath.Join(tempDir, "source")
	destinationDir := filepath.Join(tempDir, "destination")
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("failed to create source dir: %v", err)
	}
	if err := os.MkdirAll(destinationDir, 0755); err != nil {
		t.Fatalf("failed to create destination dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "track.wav"), []byte("fake wav"), 0644); err != nil {
		t.Fatalf("failed to write wav fixture: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "song.mp3"), []byte("fake mp3"), 0644); err != nil {
		t.Fatalf("failed to write mp3 fixture: %v", err)
	}

	output := captureStdout(t, func() {
		err := findAndTranscodeFiles(sourceDir, destinationDir, true)
		assert.NoError(t, err)
	})

	wavSource := filepath.Join(sourceDir, "track.wav")
	wavDest := filepath.Join(destinationDir, "track.mp3")
	mp3Source := filepath.Join(sourceDir, "song.mp3")
	mp3Dest := filepath.Join(destinationDir, "song.mp3")

	assert.True(t, strings.Contains(output, fmt.Sprintf("🔍 [dry-run] Would transcode: %s ➡️  %s", wavSource, wavDest)),
		"expected dry-run transcode message, got:\n%s", output)
	assert.True(t, strings.Contains(output, fmt.Sprintf("🔍 [dry-run] Would copy MP3: %s ➡️  %s", mp3Source, mp3Dest)),
		"expected dry-run copy message, got:\n%s", output)

	assert.NoFileExists(t, wavDest, "dry-run must not write transcoded files")
	assert.NoFileExists(t, mp3Dest, "dry-run must not copy mp3 files")
}

func TestFindFiles_DryRunFalseStillWritesFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-dryrun-off")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sourceDir := filepath.Join(tempDir, "source")
	destinationDir := filepath.Join(tempDir, "destination")
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("failed to create source dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "song.mp3"), []byte("fake mp3 content"), 0644); err != nil {
		t.Fatalf("failed to write mp3 fixture: %v", err)
	}

	err = findAndTranscodeFiles(sourceDir, destinationDir, false)
	assert.NoError(t, err)

	assert.FileExists(t, filepath.Join(destinationDir, "song.mp3"), "non-dry-run must copy mp3 files")
}
