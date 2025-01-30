package main

import (
	"fmt"
	"os"
)

func main() {
	for _, fileName := range os.Args[1:] {
		fmt.Println(fileName)
		for _, subtitles := range GetSubtitles(fileName, []string{}) {
			fmt.Println(subtitles.Path())
		}
	}
}
