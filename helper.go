package main

import (
	"path/filepath"
	"strings"
)

const maxASCIIIndex = 127

// removeNonASCII replaces non-ASCII characters in a string with an ASCII equivalent.
func removeNonASCII(str string) string {
	// Create a hashmap to store non-ASCII characters as keys and ASCII characters as values
	nonASCIItoASCII := map[rune]rune{
		'á': 'a',
		'é': 'e',
		'è': 'e',
		'ê': 'e',
		'í': 'i',
		'ó': 'o',
		'ø': 'o',
		'ú': 'u',
		'ñ': 'n',
		'Á': 'A',
		'É': 'E',
		'È': 'E',
		'Ê': 'E',
		'Í': 'I',
		'Ó': 'O',
		'Ø': 'O',
		'Ú': 'U',
		'Ñ': 'N',
		// Add more mappings as needed
	}

	// Replace non-ASCII characters with their ASCII equivalents
	var result strings.Builder
	for _, char := range str {
		asciiChar, isCharMappedToASCII := nonASCIItoASCII[char]
		switch {
		case isCharMappedToASCII:
			result.WriteRune(asciiChar)
		case char > maxASCIIIndex:
			// Don't emit char
		default:
			result.WriteRune(char)
		}
	}

	return result.String()
}

// stripArtistAlbumFromFilename removes filename segments that redundantly repeat a
// containing directory name (typically artist and/or album) from the filename portion
// of relPath. The filename (without extension) is split on "-" into segments; any
// segment matching a directory component case-insensitively is dropped. Matching is
// exact per segment (not substring). If no segment is removed, relPath is returned
// unchanged. If removing segments would leave an empty title (e.g. the track name is
// the same as the album name), the original relPath is returned unchanged.
func stripArtistAlbumFromFilename(relPath string) string {
	dir := filepath.Dir(relPath)
	base := filepath.Base(relPath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	folderNames := make(map[string]bool)
	for _, folder := range strings.Split(filepath.ToSlash(dir), "/") {
		folder = strings.TrimSpace(folder)
		if folder == "" || folder == "." {
			continue
		}
		folderNames[strings.ToLower(folder)] = true
	}

	if len(folderNames) == 0 {
		return relPath
	}

	segments := strings.Split(name, "-")
	var kept []string
	removedAny := false
	for _, segment := range segments {
		trimmed := strings.TrimSpace(segment)
		if folderNames[strings.ToLower(trimmed)] {
			removedAny = true
			continue
		}
		kept = append(kept, trimmed)
	}

	if !removedAny {
		return relPath
	}

	strippedName := strings.TrimSpace(strings.Join(kept, "-"))
	if strippedName == "" {
		return relPath
	}

	return filepath.Join(dir, strippedName+ext)
}

// isUntranscodedMusicFile checks if the path is a source music file of
// common types that need to be converted to MP3 (but are not themselves MP3),
// based on its extension. Extension matching is case-insensitive.
func isUntranscodedMusicFile(path string) bool {
	extensions := []string{".aif", ".aiff", ".aifc", ".wav", ".m4a"}
	return stringInSlice(strings.ToLower(filepath.Ext(path)), extensions)
}

// stringInSlice returns bool if a string is found in any of a list of other strings.
//
// Example usage:
//
//	if stringInSlice("Stevia", []string{"Stevie Nicks", "Stevie Wonder", "Steve Nash", "Steve McQueen"}) {
//
//	}
func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
