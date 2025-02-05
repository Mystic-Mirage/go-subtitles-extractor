package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/sirupsen/logrus"
)

const EXT = ".srt"

type Options struct {
	Libraries       []string `json:"-"`
	Sleep           int      `json:"-"`
	Forced          bool     `json:"forced"`
	SkipSrt         bool     `json:"skip_srt"`
	StripFormatting bool     `json:"strip_formatting"`
	Langs           []string `json:"langs"`
	DataDir         string   `json:"-"`
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

type File struct {
	Path     string
	UnixNano int64
}

type FileSet mapset.Set[File]

func run(options *Options) {
	cache := ReadCache(options)
	extExclude := []string{EXT, ".nfo", ".txt"}

	for {
		files := mapset.NewSet[File]()

		for _, lib := range options.Libraries {
			filepath.WalkDir(lib, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return nil
				}

				info, err := d.Info()
				if err != nil {
					return nil
				}

				modTime := info.ModTime()
				if info.Mode().IsRegular() {
					if time.Since(modTime) > 5*time.Minute {
						files.Add(File{Path: path, UnixNano: modTime.UnixNano()})
					} else {
						files.Add(File{Path: path, UnixNano: cache.Files[path].UnixNano()})
					}
				}
				return nil
			})
		}

		cachedFiles := cache.FileSet()

		for _, file := range files.Difference(cachedFiles).ToSlice() {
			if endsWith(strings.ToLower(file.Path), extExclude) || strings.Contains(file.Path, "-TdarrCacheFile-") {
				log.Println("Skipping:", file.Path)
			} else {
				log.Println("Processing:", file.Path)
				SaveSubtitles(file.Path, options)
			}
		}

		if !files.Equal(cachedFiles) {
			cache.Save(files)
		}

		time.Sleep(time.Duration(options.Sleep) * time.Minute)
	}
}

func main() {
	logrus.SetLevel(logrus.PanicLevel)

	options := getOptions()

	run(options)
}
