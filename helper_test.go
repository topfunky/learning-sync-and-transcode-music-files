package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsUntranscodedMusicFile(t *testing.T) {
	cases := []struct {
		Name           string
		Path           string
		ExpectedOutput bool
	}{
		{
			Name:           "AIFF file with long extension",
			Path:           "01_Slow_Southern_Skies.aiff",
			ExpectedOutput: true,
		},
		{
			Name:           "AIFF-C file",
			Path:           "track.aifc",
			ExpectedOutput: true,
		},
		{
			Name:           "Uppercase .AIF extension",
			Path:           "TRACK.AIF",
			ExpectedOutput: true,
		},
		{
			Name:           "Uppercase .AIFF extension",
			Path:           "TRACK.AIFF",
			ExpectedOutput: true,
		},
		{
			Name:           "Uppercase .WAV extension",
			Path:           "TRACK.WAV",
			ExpectedOutput: true,
		},
		{
			Name:           "Uppercase .M4A extension",
			Path:           "TRACK.M4A",
			ExpectedOutput: true,
		},
		{
			Name:           "Lowercase .aif extension",
			Path:           "track.aif",
			ExpectedOutput: true,
		},
		{
			Name:           "Lowercase .wav extension",
			Path:           "track.wav",
			ExpectedOutput: true,
		},
		{
			Name:           "Lowercase .m4a extension",
			Path:           "track.m4a",
			ExpectedOutput: true,
		},
		{
			Name:           "MP3 file is not transcoded",
			Path:           "track.mp3",
			ExpectedOutput: false,
		},
		{
			Name:           ".DS_Store is not a music file",
			Path:           ".DS_Store",
			ExpectedOutput: false,
		},
		{
			Name:           "Text file is not a music file",
			Path:           "notes.txt",
			ExpectedOutput: false,
		},
		{
			Name:           "Hidden file with music extension still matches by extension",
			Path:           "._hidden.aif",
			ExpectedOutput: true,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			result := isUntranscodedMusicFile(c.Path)
			assert.Equal(t, c.ExpectedOutput, result)
		})
	}
}

func TestStripArtistAlbumFromFilename(t *testing.T) {
	cases := []struct {
		Name           string
		Path           string
		ExpectedOutput string
	}{
		{
			Name:           "Repeats artist and album name",
			Path:           "Ulrich Schnauss/Far Away Trains Passing By/Far Away Trains Passing By-01-001-Ulrich Schnauss-Knuddelmaus (2019 remaster).mp3",
			ExpectedOutput: "Ulrich Schnauss/Far Away Trains Passing By/01-001-Knuddelmaus (2019 remaster).mp3",
		},
		{
			Name:           "Case-insensitive match",
			Path:           "Ulrich Schnauss/Far Away Trains Passing By/far away trains passing by-01-001-ULRICH SCHNAUSS-Knuddelmaus.mp3",
			ExpectedOutput: "Ulrich Schnauss/Far Away Trains Passing By/01-001-Knuddelmaus.mp3",
		},
		{
			Name:           "Song title matches album name is restored to original filename",
			Path:           "Artist/Album/Artist-Album.mp3",
			ExpectedOutput: "Artist/Album/Artist-Album.mp3",
		},
		{
			Name:           "No matching segments passes through unchanged",
			Path:           "Artist/Album/01 - Song Title.mp3",
			ExpectedOutput: "Artist/Album/01 - Song Title.mp3",
		},
		{
			Name:           "Flat path with no directories passes through unchanged",
			Path:           "01 - Song Title.mp3",
			ExpectedOutput: "01 - Song Title.mp3",
		},
		{
			Name:           "Substring match is not stripped",
			Path:           "Artist/Album/01-Artiste-Song Title.mp3",
			ExpectedOutput: "Artist/Album/01-Artiste-Song Title.mp3",
		},
		{
			Name:           "Nested directories can all be stripped",
			Path:           "Artist/Album/Disc 1/Artist-Disc 1-01-Song Title.mp3",
			ExpectedOutput: "Artist/Album/Disc 1/01-Song Title.mp3",
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			result := stripArtistAlbumFromFilename(c.Path)
			assert.Equal(t, c.ExpectedOutput, result)
		})
	}
}
