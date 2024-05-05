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
			ExpectedOutput:  []string(nil),
		},
		{
			Name:            "Empty source list, destination list has elements",
			SourceList:      []string{},
			DestinationList: []string{"file1.mp3", "file2.mp3"},
			ExpectedOutput:  []string(nil),
		},
		{
			Name:            "Source list has elements, empty destination list",
			SourceList:      []string{"file1.m4a", "file2.m4a"},
			DestinationList: []string{},
			ExpectedOutput:  []string{"file1.m4a", "file2.m4a"},
		},
		{
			Name:            "Source and destination have common elements",
			SourceList:      []string{"file1.m4a", "file2.m4a", "file3.m4a"},
			DestinationList: []string{"file2.mp3", "file3.mp3", "file4.mp3"},
			ExpectedOutput:  []string{"file1.m4a"},
		},
		{
			Name:            "Source and destination have disjoint elements",
			SourceList:      []string{"file1.m4a", "file2.m4a", "file3.m4a"},
			DestinationList: []string{"file4.mp3", "file5.mp3", "file6.mp3"},
			ExpectedOutput:  []string{"file1.m4a", "file2.m4a", "file3.m4a"},
		},
		{
			Name:            "Source contains mp3 files which should be copied verbatim",
			SourceList:      []string{"file1.m4a", "file2.m4a", "file3.m4a", "file103.mp3"},
			DestinationList: []string{"file4.mp3", "file5.mp3", "file6.mp3"},
			ExpectedOutput:  []string{"file1.m4a", "file2.m4a", "file3.m4a", "file103.mp3"},
		},
		{
			Name:            "Destination contains m4a files which should be ignored",
			SourceList:      []string{"file1.m4a", "file2.m4a", "file3.m4a"},
			DestinationList: []string{"file4.m4a", "file5.mp3", "file6.mp3"},
			ExpectedOutput:  []string{"file1.m4a", "file2.m4a", "file3.m4a"},
		},
		{
			Name:            "Destination contains aif and wav files which should be transcoded",
			SourceList:      []string{"file1.m4a", "file2.aif", "file3.wav"},
			DestinationList: []string{},
			ExpectedOutput:  []string{"file1.m4a", "file2.aif", "file3.wav"},
		},
		{
			Name:            "Ignore non-music files",
			SourceList:      []string{".DS_Store"},
			DestinationList: []string{},
			ExpectedOutput:  []string(nil),
		},
		{
			Name:            "Correctly compares non-ASCII filenames",
			SourceList:      []string{"Megan Perry Fisher/Megan Perry Fisher - Pensées/Megan Perry Fisher - Pensées - 12 Pensée xii.m4a", "Stéphane Grappelli, Joe Pass & Niels-Henning Ørsted Pedersen/Tivoli Gardens, Copenhagen, Denmark (Live)/01 It's Only A Paper Moon.m4a"},
			DestinationList: []string{},
			ExpectedOutput:  []string{"Megan Perry Fisher/Megan Perry Fisher - Pensées/Megan Perry Fisher - Pensées - 12 Pensée xii.m4a", "Stéphane Grappelli, Joe Pass & Niels-Henning Ørsted Pedersen/Tivoli Gardens, Copenhagen, Denmark (Live)/01 It's Only A Paper Moon.m4a"},
		},
		{
			Name:            "Does not re-transcode non-ASCII filenames",
			SourceList:      []string{"Megan Perry Fisher/Megan Perry Fisher - Pensées/Megan Perry Fisher - Pensées - 12 Pensée xii.m4a", "Stéphane Grappelli, Joe Pass & Niels-Henning Ørsted Pedersen/Tivoli Gardens, Copenhagen, Denmark (Live)/01 It's Only A Paper Moon.m4a"},
			DestinationList: []string{"Megan Perry Fisher/Megan Perry Fisher - Pensées/Megan Perry Fisher - Pensées - 12 Pensée xii.mp3", "Stéphane Grappelli, Joe Pass & Niels-Henning Ørsted Pedersen/Tivoli Gardens, Copenhagen, Denmark (Live)/01 It's Only A Paper Moon.mp3"},
			ExpectedOutput:  []string(nil),
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			result := getExclusiveFiles(c.SourceList, c.DestinationList)
			assert.Equal(t, c.ExpectedOutput, getSourcePaths(result))
		})
	}
}

// Returns a string array of only the `sourcePath` attribute from an array of `fileToRender` structs.
//
// This makes test assertions cleaner, based on how the fixture data is written.
func getSourcePaths(files []fileToRender) []string {
	var sources []string
	for _, file := range files {
		sources = append(sources, file.sourcePath)
	}
	return sources
}
