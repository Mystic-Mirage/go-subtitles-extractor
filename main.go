package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const EXT = ".srt"

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
	stripFormatting := os.Getenv("SUBTITLES_EXTRACTOR_STRIP_FORMATTING")
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

func endsWith(s string, suffixes []string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}
	return false
}

func run(options *Options) {
	cache := ReadCache(options)
	extExclude := []string{EXT, ".nfo", ".txt"}

	for {
		files := Files{}
		overwriteCache := false

		for _, lib := range options.Libraries {
			filepath.WalkDir(lib, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return nil
				}

				info, err := d.Info()
				if err != nil {
					return nil
				}

				if info.Mode().IsRegular() {
					modTime := info.ModTime()

					cachedModTime, ok := cache.Files[path]
					if !ok || !modTime.Equal(cachedModTime) {
						if time.Since(modTime) > 5*time.Minute {
							if endsWith(strings.ToLower(path), extExclude) || strings.Contains(path, "-TdarrCacheFile-") {
								log.Println("Skipping:", path)
							} else {
								log.Println("Processing:", path)
								SaveSubtitles(path, options)
							}
						} else if ok {
							// keep existing cached time for a while
							files[path] = cachedModTime
							return nil
						} else {
							// don't cache a newly created file
							return nil
						}

						overwriteCache = true
					}

					files[path] = modTime
				}
				return nil
			})
		}

		if overwriteCache || len(files) != len(cache.Files) {
			cache.Save(files)
		}

		time.Sleep(time.Duration(options.Sleep) * time.Minute)
	}
}

func main() {
	options := getOptions()

	run(options)
}
