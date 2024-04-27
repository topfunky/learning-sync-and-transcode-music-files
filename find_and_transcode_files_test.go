package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetExclusiveFiles(t *testing.T) {
	// Test case 1: filesA is empty, filesB is empty
	filesA := []string{}
	filesB := []string{}
	expected := []string{}
	result := getExclusiveFiles(filesA, filesB)
	assert.Equal(t, expected, result)

	// Test case 2: filesA is empty, filesB has elements
	filesA = []string{}
	filesB = []string{"file1.txt", "file2.txt"}
	expected = []string{}
	result = getExclusiveFiles(filesA, filesB)
	assert.Equal(t, expected, result)

	// Test case 3: filesA has elements, filesB is empty
	filesA = []string{"file1.txt", "file2.txt"}
	filesB = []string{}
	expected = []string{"file1.txt", "file2.txt"}
	result = getExclusiveFiles(filesA, filesB)
	assert.Equal(t, expected, result)

	// Test case 4: filesA and filesB have common elements
	filesA = []string{"file1.txt", "file2.txt", "file3.txt"}
	filesB = []string{"file2.txt", "file3.txt", "file4.txt"}
	expected = []string{"file1.txt"}
	result = getExclusiveFiles(filesA, filesB)
	assert.Equal(t, expected, result)

	// Test case 5: filesA and filesB have no common elements
	filesA = []string{"file1.txt", "file2.txt", "file3.txt"}
	filesB = []string{"file4.txt", "file5.txt", "file6.txt"}
	expected = []string{"file1.txt", "file2.txt", "file3.txt"}
	result = getExclusiveFiles(filesA, filesB)
	assert.Equal(t, expected, result)
}
