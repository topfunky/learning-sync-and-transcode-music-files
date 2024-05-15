package main

import (
	"path/filepath"
	"strings"
)

const maxASCIIIndex = 127

// containsNonASCII returns true if a string contains non-ASCII characters.
func containsNonASCII(str string) bool {
	for _, char := range str {
		if char > maxASCIIIndex {
			return true
		}
	}
	return false
}

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

// isUntranscodedMusicFile checks if the path is a source music file of common types that need to be converted to MP3 (but are not themselves MP3), based on its extension.
func isUntranscodedMusicFile(path string) bool {
	extensions := []string{".aif", ".wav", ".m4a"}
	return stringInSlice(filepath.Ext(path), extensions)
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
