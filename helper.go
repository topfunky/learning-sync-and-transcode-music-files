package main

// containsNonASCII returns true if a string contains non-ASCII characters.
func containsNonASCII(str string) bool {
	for _, char := range str {
		if char > 127 {
			return true
		}
	}
	return false
}
