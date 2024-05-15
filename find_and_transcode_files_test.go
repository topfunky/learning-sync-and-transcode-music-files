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
		{
			Name:            "Destination contains aif and wav files which should be transcoded",
			SourceList:      []string{"file1.m4a", "file2.aif", "file3.wav"},
			DestinationList: []string{},
			ExpectedOutput:  []string{"file1.mp3", "file2.mp3", "file3.mp3"},
		},
		{
			Name:            "Ignore non-music files",
			SourceList:      []string{".DS_Store"},
			DestinationList: []string{},
			ExpectedOutput:  []string(nil),
		},
		{
			Name:            "Correctly compares non-ASCII filenames",
			SourceList:      []string{"Alexandra Stréliski/Néo-Romance (Extended Version) [96kHz · 24bit]/02 - Lumières.m4a"},
			DestinationList: []string{},
			ExpectedOutput:  []string{"Alexandra Streliski/Neo-Romance (Extended Version) [96kHz  24bit]/02 - Lumieres.mp3"},
		},
		{
			Name:            "Correctly compares non-ASCII filenames (alt)",
			SourceList:      []string{"Megan Perry Fisher/Megan Perry Fisher - Pensées/Megan Perry Fisher - Pensées - 12 Pensée xii.m4a", "Stéphane Grappelli, Joe Pass & Niels-Henning Ørsted Pedersen/Tivoli Gardens, Copenhagen, Denmark (Live)/01 It's Only A Paper Moon.m4a"},
			DestinationList: []string{},
			ExpectedOutput:  []string{"Megan Perry Fisher/Megan Perry Fisher - Pensees/Megan Perry Fisher - Pensees - 12 Pensee xii.mp3", "Stephane Grappelli, Joe Pass & Niels-Henning Orsted Pedersen/Tivoli Gardens, Copenhagen, Denmark (Live)/01 It's Only A Paper Moon.mp3"},
		},
		{
			Name:            "Does not re-transcode non-ASCII filenames",
			SourceList:      []string{"Megan Perry Fisher/Megan Perry Fisher - Pensées/Megan Perry Fisher - Pensées - 12 Pensée xii.m4a", "Stéphane Grappelli, Joe Pass & Niels-Henning Ørsted Pedersen/Tivoli Gardens, Copenhagen, Denmark (Live)/01 It's Only A Paper Moon.m4a"},
			DestinationList: []string{"Megan Perry Fisher/Megan Perry Fisher - Pensees/Megan Perry Fisher - Pensees - 12 Pensee xii.mp3", "Stephane Grappelli, Joe Pass & Niels-Henning Orsted Pedersen/Tivoli Gardens, Copenhagen, Denmark (Live)/01 It's Only A Paper Moon.mp3"},
			ExpectedOutput:  []string(nil),
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			result := getExclusiveFiles(c.SourceList, c.DestinationList)
			assert.Equal(t, c.ExpectedOutput, getDestinationPaths(result))
		})
	}
}

// Returns a string array of only the `sourcePath` attribute from an array of `fileToRender` structs.
//
// This makes test assertions cleaner, based on how the fixture data is written.
func getDestinationPaths(files []fileToRender) []string {
	var sources []string
	for _, file := range files {
		sources = append(sources, file.destinationPath)
	}
	return sources
}
func TestConvertSourceToDestinationFilename(t *testing.T) {
	cases := []struct {
		Name           string
		Filename       string
		ExpectedOutput string
	}{
		{
			Name:           "Filename with .m4a extension",
			Filename:       "file1.m4a",
			ExpectedOutput: "file1.mp3",
		},
		{
			Name:           "Filename with .m4a extension and non-ASCII characters",
			Filename:       "Megan Perry Fisher - Pensées.m4a",
			ExpectedOutput: "Megan Perry Fisher - Pensees.mp3",
		},
		{
			Name:           "Filename with .mp3 extension",
			Filename:       "file2.mp3",
			ExpectedOutput: "file2.mp3",
		},
		{
			Name:           "Filename with no extension",
			Filename:       "file3",
			ExpectedOutput: "file3.mp3",
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			result := convertSourceToDestinationFilename(c.Filename)
			assert.Equal(t, c.ExpectedOutput, result)
		})
	}
}
