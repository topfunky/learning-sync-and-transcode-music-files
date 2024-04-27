package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type FixtureParams struct {
	SourceList      []string
	DestinationList []string
	ExpectedOutput  []string
}

func TestGetExclusiveFiles(t *testing.T) {
	// Test case 1: filesA is empty, filesB is empty
	testGetExclusiveFiles(t, FixtureParams{
		SourceList:      []string{},
		DestinationList: []string{},
		ExpectedOutput:  []string{},
	})

	// Test case 2: filesA is empty, filesB has elements
	testGetExclusiveFiles(t, FixtureParams{
		SourceList:      []string{},
		DestinationList: []string{"file1.txt", "file2.txt"},
		ExpectedOutput:  []string{},
	})

	// Test case 3: filesA has elements, filesB is empty
	testGetExclusiveFiles(t, FixtureParams{
		SourceList:      []string{"file1.txt", "file2.txt"},
		DestinationList: []string{},
		ExpectedOutput:  []string{"file1.txt", "file2.txt"},
	})

	// Test case 4: filesA and filesB have common elements
	testGetExclusiveFiles(t, FixtureParams{
		SourceList:      []string{"file1.txt", "file2.txt", "file3.txt"},
		DestinationList: []string{"file2.txt", "file3.txt", "file4.txt"},
		ExpectedOutput:  []string{"file1.txt"},
	})

	// Test case 5: filesA and filesB have no common elements
	testGetExclusiveFiles(t, FixtureParams{
		SourceList:      []string{"file1.txt", "file2.txt", "file3.txt"},
		DestinationList: []string{"file4.txt", "file5.txt", "file6.txt"},
		ExpectedOutput:  []string{"file1.txt", "file2.txt", "file3.txt"},
	})
}

func testGetExclusiveFiles(t *testing.T, params FixtureParams) {
	result := getExclusiveFiles(params.SourceList, params.DestinationList)
	assert.Equal(t, params.ExpectedOutput, result)
}
