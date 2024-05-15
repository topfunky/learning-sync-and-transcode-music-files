package main

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
