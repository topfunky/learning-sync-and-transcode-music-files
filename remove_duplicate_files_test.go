package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupFilesByBasePath(t *testing.T) {
	cases := []struct {
		Name           string
		Input          []string
		ExpectedGroups map[string][]string
	}{
		{
			Name:           "Empty input",
			Input:          []string{},
			ExpectedGroups: map[string][]string{},
		},
		{
			Name:  "No duplicates",
			Input: []string{"/dest/song.mp3", "/dest/other.mp3"},
			ExpectedGroups: map[string][]string{
				"/dest/song":  {"/dest/song.mp3"},
				"/dest/other": {"/dest/other.mp3"},
			},
		},
		{
			Name:  "MP3 and M4A with same base name are grouped together",
			Input: []string{"/dest/song.mp3", "/dest/song.m4a"},
			ExpectedGroups: map[string][]string{
				"/dest/song": {"/dest/song.mp3", "/dest/song.m4a"},
			},
		},
		{
			Name:  "Files in subdirectories are grouped by their full base path",
			Input: []string{"/dest/Artist/song.mp3", "/dest/Artist/song.m4a", "/dest/Other/song.mp3"},
			ExpectedGroups: map[string][]string{
				"/dest/Artist/song": {"/dest/Artist/song.mp3", "/dest/Artist/song.m4a"},
				"/dest/Other/song":  {"/dest/Other/song.mp3"},
			},
		},
		{
			Name:  "MP3 and M4A alongside an unrelated file",
			Input: []string{"/dest/a.mp3", "/dest/a.m4a", "/dest/b.mp3"},
			ExpectedGroups: map[string][]string{
				"/dest/a": {"/dest/a.mp3", "/dest/a.m4a"},
				"/dest/b": {"/dest/b.mp3"},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			result := groupFilesByBasePath(c.Input)
			assert.Equal(t, c.ExpectedGroups, result)
		})
	}
}

func TestSelectPreferredFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-select-preferred")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	smallMP3 := filepath.Join(tempDir, "song_small.mp3")
	largeMP3 := filepath.Join(tempDir, "song_large.mp3")
	m4aFile := filepath.Join(tempDir, "song.m4a")
	aifFile := filepath.Join(tempDir, "song.aif")

	// m4aFile is larger in bytes than smallMP3, but MP3 should still win
	os.WriteFile(smallMP3, make([]byte, 100), 0644)
	os.WriteFile(largeMP3, make([]byte, 500), 0644)
	os.WriteFile(m4aFile, make([]byte, 1000), 0644)
	os.WriteFile(aifFile, make([]byte, 800), 0644)

	cases := []struct {
		Name           string
		Candidates     []string
		ExpectedKeep   string
		ExpectedDelete []string
	}{
		{
			Name:           "Single file: nothing to delete",
			Candidates:     []string{smallMP3},
			ExpectedKeep:   smallMP3,
			ExpectedDelete: nil,
		},
		{
			Name:           "MP3 beats M4A regardless of file size",
			Candidates:     []string{smallMP3, m4aFile},
			ExpectedKeep:   smallMP3,
			ExpectedDelete: []string{m4aFile},
		},
		{
			Name:           "MP3 beats AIF",
			Candidates:     []string{smallMP3, aifFile},
			ExpectedKeep:   smallMP3,
			ExpectedDelete: []string{aifFile},
		},
		{
			Name:           "Larger MP3 is preferred over smaller MP3",
			Candidates:     []string{smallMP3, largeMP3},
			ExpectedKeep:   largeMP3,
			ExpectedDelete: []string{smallMP3},
		},
		{
			Name:           "Largest MP3 wins; non-MP3 and smaller MP3 are both deleted",
			Candidates:     []string{smallMP3, largeMP3, m4aFile},
			ExpectedKeep:   largeMP3,
			ExpectedDelete: []string{smallMP3, m4aFile},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			keep, toDelete, err := selectPreferredFile(c.Candidates)
			assert.NoError(t, err)
			assert.Equal(t, c.ExpectedKeep, keep)
			assert.ElementsMatch(t, c.ExpectedDelete, toDelete)
		})
	}
}

func TestFindDuplicates(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-find-duplicates")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create files: one duplicate pair and one unique file
	os.WriteFile(filepath.Join(tempDir, "song.mp3"), make([]byte, 200), 0644)
	os.WriteFile(filepath.Join(tempDir, "song.m4a"), make([]byte, 400), 0644)
	os.WriteFile(filepath.Join(tempDir, "other.mp3"), make([]byte, 300), 0644)

	duplicates, err := findDuplicates(tempDir)
	assert.NoError(t, err)

	// Only "song" should appear as a duplicate group; "other" has no duplicate
	assert.Len(t, duplicates, 1)
	basePath := filepath.Join(tempDir, "song")
	assert.Contains(t, duplicates, basePath)
	assert.ElementsMatch(t, duplicates[basePath], []string{
		filepath.Join(tempDir, "song.mp3"),
		filepath.Join(tempDir, "song.m4a"),
	})
}

func TestRemoveDuplicateFiles_DryRun(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-remove-duplicates-dry")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	mp3File := filepath.Join(tempDir, "song.mp3")
	m4aFile := filepath.Join(tempDir, "song.m4a")
	os.WriteFile(mp3File, make([]byte, 200), 0644)
	os.WriteFile(m4aFile, make([]byte, 300), 0644)

	err = removeDuplicateFiles(tempDir, true)
	assert.NoError(t, err)

	// Dry run must not delete anything
	assert.FileExists(t, mp3File)
	assert.FileExists(t, m4aFile)
}

func TestRemoveDuplicateFiles_DeletesNonMP3Duplicate(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-remove-duplicates")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	mp3File := filepath.Join(tempDir, "song.mp3")
	m4aFile := filepath.Join(tempDir, "song.m4a")
	os.WriteFile(mp3File, make([]byte, 200), 0644)
	os.WriteFile(m4aFile, make([]byte, 300), 0644)

	err = removeDuplicateFiles(tempDir, false)
	assert.NoError(t, err)

	// MP3 should be kept, M4A should be deleted
	assert.FileExists(t, mp3File)
	assert.NoFileExists(t, m4aFile)
}

func TestSelectPreferredFile_KeepsLargerMP3(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-select-preferred-mp3")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	subDir := filepath.Join(tempDir, "Artist", "Album")
	os.MkdirAll(subDir, 0755)

	// selectPreferredFile receives a pre-built list of candidates (assembled by
	// findDuplicates). Test it directly with two real MP3 files of different
	// sizes to verify it picks the larger one.
	smallMP3 := filepath.Join(subDir, "01 - Song (128kbps).mp3")
	largeMP3 := filepath.Join(subDir, "01 - Song (320kbps).mp3")
	os.WriteFile(smallMP3, make([]byte, 500), 0644)
	os.WriteFile(largeMP3, make([]byte, 2000), 0644)

	keep, toDelete, err := selectPreferredFile([]string{smallMP3, largeMP3})
	assert.NoError(t, err)
	assert.Equal(t, largeMP3, keep)
	assert.ElementsMatch(t, []string{smallMP3}, toDelete)
}

func TestSelectPreferredFile_EmptyCandidates(t *testing.T) {
	_, _, err := selectPreferredFile([]string{})
	assert.Error(t, err)
}

func TestFindDuplicates_NonExistentDirectory(t *testing.T) {
	_, err := findDuplicates("/nonexistent/dir")
	assert.Error(t, err)
}

func TestRemoveDuplicateFiles_NonExistentDirectory(t *testing.T) {
	err := removeDuplicateFiles("/nonexistent/dir", false)
	assert.Error(t, err)
}

func TestRemoveDuplicateFiles_NoFilesInDirectory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-remove-duplicates-empty")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	err = removeDuplicateFiles(tempDir, false)
	assert.NoError(t, err)
}

func TestRemoveDuplicateFiles_InSubdirectory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-remove-duplicates-sub")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	subDir := filepath.Join(tempDir, "Artist", "Album")
	os.MkdirAll(subDir, 0755)

	mp3File := filepath.Join(subDir, "01 - Song.mp3")
	m4aFile := filepath.Join(subDir, "01 - Song.m4a")
	os.WriteFile(mp3File, make([]byte, 300), 0644)
	os.WriteFile(m4aFile, make([]byte, 600), 0644)

	err = removeDuplicateFiles(tempDir, false)
	assert.NoError(t, err)

	assert.FileExists(t, mp3File)
	assert.NoFileExists(t, m4aFile)
}
