package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// groupFilesByBasePath groups file paths by their path without the file extension.
// The key is the full path minus the extension; the value is every file that
// shares that base path (potentially across different extensions).
func groupFilesByBasePath(files []string) map[string][]string {
	groups := make(map[string][]string)
	for _, f := range files {
		basePath := strings.TrimSuffix(f, filepath.Ext(f))
		groups[basePath] = append(groups[basePath], f)
	}
	return groups
}

// selectPreferredFile returns the file to keep and the files to delete from a
// group of candidates that share the same base path.
//
// Selection rules (in order):
//  1. Files ending in .mp3 (case-insensitive) are preferred over all other formats.
//  2. When multiple .mp3 files are present, the largest file (by byte size) is
//     kept as a proxy for higher bit rate.
//
// Returns (fileToKeep, filesToDelete, error).
func selectPreferredFile(candidates []string) (string, []string, error) {
	if len(candidates) == 0 {
		return "", nil, fmt.Errorf("no candidates provided")
	}
	if len(candidates) == 1 {
		return candidates[0], nil, nil
	}

	// Partition candidates into MP3 files and everything else.
	var mp3Files []string
	var otherFiles []string
	for _, f := range candidates {
		if strings.EqualFold(filepath.Ext(f), ".mp3") {
			mp3Files = append(mp3Files, f)
		} else {
			otherFiles = append(otherFiles, f)
		}
	}

	// Decide which pool to pick the best file from.
	pickFrom := candidates
	deleteAll := []string(nil)
	if len(mp3Files) > 0 {
		pickFrom = mp3Files
		deleteAll = otherFiles
	}

	// Among the chosen pool, keep the largest file (highest bit-rate proxy).
	bestFile := ""
	bestSize := int64(-1)
	for _, f := range pickFrom {
		info, err := os.Stat(f)
		if err != nil {
			return "", nil, fmt.Errorf("failed to stat %s: %v", f, err)
		}
		if info.Size() > bestSize {
			bestSize = info.Size()
			bestFile = f
		}
	}

	// Every other file in the chosen pool is also a duplicate to delete.
	for _, f := range pickFrom {
		if f != bestFile {
			deleteAll = append(deleteAll, f)
		}
	}

	return bestFile, deleteAll, nil
}

// findDuplicates scans a directory and returns groups of files that share the
// same base path (full path without extension) and have more than one member.
func findDuplicates(dir string) (map[string][]string, error) {
	relPaths, err := getFilenames(dir)
	if err != nil {
		return nil, err
	}

	var absPaths []string
	for _, rel := range relPaths {
		absPaths = append(absPaths, filepath.Join(dir, rel))
	}

	groups := groupFilesByBasePath(absPaths)
	duplicates := make(map[string][]string)
	for basePath, group := range groups {
		if len(group) > 1 {
			duplicates[basePath] = group
		}
	}
	return duplicates, nil
}

// removeDuplicateFiles scans the destination directory for duplicate files
// (same base path, different extensions or multiple MP3s) and removes the
// lower-quality copies. When dryRun is true it only prints what would be
// deleted without removing anything.
func removeDuplicateFiles(dir string, dryRun bool) error {
	duplicates, err := findDuplicates(dir)
	if err != nil {
		return fmt.Errorf("error finding duplicates: %v", err)
	}

	for _, candidates := range duplicates {
		keep, toDelete, err := selectPreferredFile(candidates)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❗️ Error selecting preferred file: %v\n", err)
			continue
		}

		fmt.Printf("✅ Keeping: %s\n", keep)
		for _, f := range toDelete {
			if dryRun {
				fmt.Printf("🔍 [dry-run] Would delete duplicate: %s\n", f)
			} else {
				if err := os.Remove(f); err != nil {
					fmt.Fprintf(os.Stderr, "❗️ Error deleting duplicate %s: %v\n", f, err)
				} else {
					fmt.Printf("🗑️  Deleted duplicate: %s\n", f)
				}
			}
		}
	}
	return nil
}
