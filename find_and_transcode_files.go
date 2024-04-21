package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xfrr/goffmpeg/transcoder"
)

func findFiles(directory string) {
	fmt.Printf("🔨 Transcoding for directory %s\n", directory)
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), ".m4a") {
			err := transcodeFileAtPath(path)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path: %v\n", err)
		return
	}
}

func transcodeFileAtPath(path string) error {
	trans := new(transcoder.Transcoder)
	output := strings.TrimSuffix(path, filepath.Ext(path)) + ".mp3"
	err := trans.Initialize(path, output)
	if err != nil {
		return err
	}

	done := trans.Run(false)
	err = <-done
	if err != nil {
		return err
	}

	fmt.Printf("🎶 Transcoded: %s to %s\n", path, output)
	return nil
}
