package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

const FILENAME = "filelist.json"

type Files map[string]time.Time

type Cache struct {
	Options *Options `json:"options"`
	Files   Files    `json:"files"`
}

func (c *Cache) Check(options *Options) bool {
	return c.Options.Forced == options.Forced &&
		c.Options.SkipSrt == options.SkipSrt &&
		c.Options.StripFormatting == options.StripFormatting &&
		reflect.DeepEqual(c.Options.Langs, options.Langs) &&
		reflect.DeepEqual(c.Options.ForcedTitles, options.ForcedTitles)
}

var ErrOptionsMismatch error = errors.New("ErrOptionsMismatch")

func (c *Cache) Validate(options *Options) error {
	if !c.Check(options) {
		return ErrOptionsMismatch
	}
	return nil
}

func (c *Cache) Save(files Files) {
	c.Files = files
	fullName := filepath.Join(c.Options.DataDir, FILENAME)
	bytes, _ := json.MarshalIndent(c, "", "  ")
	os.WriteFile(fullName, bytes, os.ModePerm)
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
