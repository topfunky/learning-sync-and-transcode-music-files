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
			ExpectedOutput:  []string{},
		},
		{
			Name:            "Empty source list, destination list has elements",
			SourceList:      []string{},
			DestinationList: []string{"file1.txt", "file2.txt"},
			ExpectedOutput:  []string{},
		},
		{
			Name:            "Source list has elements, empty destination list",
			SourceList:      []string{"file1.txt", "file2.txt"},
			DestinationList: []string{},
			ExpectedOutput:  []string{"file1.txt", "file2.txt"},
		},
		{
			Name:            "Source and destination have common elements",
			SourceList:      []string{"file1.txt", "file2.txt", "file3.txt"},
			DestinationList: []string{"file2.txt", "file3.txt", "file4.txt"},
			ExpectedOutput:  []string{"file1.txt"},
		},
		{
			Name:            "Source and destination have disjoint elements",
			SourceList:      []string{"file1.txt", "file2.txt", "file3.txt"},
			DestinationList: []string{"file4.txt", "file5.txt", "file6.txt"},
			ExpectedOutput:  []string{"file1.txt", "file2.txt", "file3.txt"},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			result := getExclusiveFiles(c.SourceList, c.DestinationList)
			assert.Equal(t, c.ExpectedOutput, result)
		})
	}
}
