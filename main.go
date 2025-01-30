package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Options struct {
	Libraries       []string `json:"libraries"`
	Sleep           int      `json:"sleep"`
	Forced          bool     `json:"forced"`
	SkipSrt         bool     `json:"skip_srt"`
	StripFormatting bool     `json:"strip_formatting"`
	Langs           []string `json:"langs"`
	DataDir         string   `json:"data_dir"`
	ForcedTitles    []string `json:"forced_titles"`
}

func getOptions() *Options {
	options := &Options{
		Libraries:       []string{},
		Sleep:           1,
		Forced:          true,
		SkipSrt:         true,
		StripFormatting: true,
		Langs:           []string{"*"},
		DataDir:         ".",
		ForcedTitles:    []string{},
	}

	libraries := os.Getenv("SUBTITLES_EXTRACTOR_LIBRARIES")
	sleep := os.Getenv("SUBTITLES_EXTRACTOR_SLEEP")
	forced := os.Getenv("SUBTITLES_EXTRACTOR_FORCED_ONLY")
	skipSrt := os.Getenv("SUBTITLES_EXTRACTOR_SKIP_SRT")
	stripFormatting := os.Getenv("SUBTITLES_EXTRACTOR_SKIP_SRT")
	langs := os.Getenv("SUBTITLES_EXTRACTOR_LANGUAGES")
	dataDir := os.Getenv("SUBTITLES_EXTRACTOR_DATA_DIR")
	forcedTitles := os.Getenv("SUBTITLES_EXTRACTOR_FORCED_TITLE")

	if libraries != "" {
		options.Libraries = strings.Split(libraries, ";")
	}

	if sleep != "" {
		value, err := strconv.Atoi(sleep)
		if err == nil {
			options.Sleep = value
		}
	}

	if forced != "" {
		value, err := strconv.ParseBool(forced)
		if err == nil {
			options.Forced = value
		}
	}

	if skipSrt != "" {
		value, err := strconv.ParseBool(skipSrt)
		if err == nil {
			options.SkipSrt = value
		}
	}

	if stripFormatting != "" {
		value, err := strconv.ParseBool(stripFormatting)
		if err == nil {
			options.StripFormatting = value
		}
	}

	if langs != "" {
		options.Langs = strings.Split(langs, ";")
	}

	if dataDir != "" {
		options.DataDir = dataDir
	}

	if forcedTitles != "" {
		options.ForcedTitles = strings.Split(forcedTitles, ";")
	}

	return options
}

func main() {
	options := getOptions()

	for _, fileName := range os.Args[1:] {
		fmt.Println(fileName)

		videoFile := VideoFile{FileName: fileName}

		for _, subtitles := range videoFile.Subtitles(options.ForcedTitles) {
			fmt.Println(subtitles.Path())
		}
	}

	cache := ReadCache(options)
	cache.Save()
}
