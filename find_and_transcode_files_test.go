package main

import (
	"testing"

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
			ExpectedOutput:  []string{},
		},
		{
			Name:            "Empty source list, destination list has elements",
			SourceList:      []string{},
			DestinationList: []string{"file1.mp3", "file2.mp3"},
			ExpectedOutput:  []string{},
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
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			result := getExclusiveFiles(c.SourceList, c.DestinationList)
			assert.Equal(t, c.ExpectedOutput, result)
		})
	}
}
