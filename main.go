package main

import (
	"fmt"
	"os"
)

func main() {
	for _, fileName := range os.Args[1:] {
		fmt.Println(fileName)

		videoFile := VideoFile{FileName: fileName}

		for _, subtitles := range videoFile.Subtitles([]string{}) {
			fmt.Println(subtitles.Path())
		}
	}
}
