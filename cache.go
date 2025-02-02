package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"slices"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
)

const FILENAME = "filelist.json"

type Files map[string]time.Time

type Cache struct {
	Options *Options `json:"options"`
	Files   Files    `json:"files"`
}

func (c *Cache) FileSet() FileSet {
	set := mapset.NewSet[File]()
	for path, modTime := range c.Files {
		set.Add(File{Path: path, UnixNano: modTime.UnixNano()})
	}
	return set
}

func (c *Cache) Check(options *Options) bool {
	return c.Options.Forced == options.Forced &&
		c.Options.SkipSrt == options.SkipSrt &&
		c.Options.StripFormatting == options.StripFormatting &&
		slices.Equal(c.Options.Langs, options.Langs) &&
		slices.Equal(c.Options.ForcedTitles, options.ForcedTitles)
}

var ErrOptionsMismatch error = errors.New("ErrOptionsMismatch")

func (c *Cache) Validate(options *Options) error {
	if !c.Check(options) {
		return ErrOptionsMismatch
	}
	return nil
}

func (c *Cache) Save(files FileSet) {
	clear(c.Files)
	for _, file := range files.ToSlice() {
		c.Files[file.Path] = time.Unix(0, file.UnixNano)
	}
	fullName := filepath.Join(c.Options.DataDir, FILENAME)
	bytes, _ := json.MarshalIndent(c, "", "  ")
	os.WriteFile(fullName, bytes, 0o644)
}

func ReadCache(options *Options) *Cache {
	fullName := filepath.Join(options.DataDir, FILENAME)
	cache := &Cache{Files: Files{}}

	bytes, err := os.ReadFile(fullName)
	if err == nil {
		json.Unmarshal(bytes, cache)
		err = cache.Validate(options)
	}

	if err == nil {
		log.Println("Cache size:", len(cache.Files))
		return cache
	} else {
		cache.Options = options
		clear(cache.Files)
		log.Println("Cache (re)initialized")
	}

	return cache
}
