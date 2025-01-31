package main

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func SaveSubtitles(fileName string, options *Options) {
	videoFile := VideoFile{FileName: fileName}
	for _, subtitles := range videoFile.Subtitles(options.ForcedTitles) {
		dst := subtitles.Path()

		if !subtitles.Check(options) {
			err := os.Remove(dst)
			if err == nil {
				log.Println("Unlinked:", dst)
			}
			continue
		}

		stdErr := strings.Builder{}
		command := exec.Command(
			"ffmpeg",
			"-v",
			"error",
			"-y",
			"-i",
			fileName,
			"-map",
			"0:"+strconv.Itoa(subtitles.Index),
			"-f",
			"srt",
			"pipe:1",
		)
		command.Stderr = &stdErr

		output, err := command.Output()
		if err != nil {
			log.Print("FFMPEG error:", stdErr.String())
			continue
		}

		if options.StripFormatting {
			output = []byte(Strip(string(output)))
		}

		os.WriteFile(dst, output, 0o644)
		log.Println("Extracted:", dst)
	}
}
