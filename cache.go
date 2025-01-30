package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"reflect"
)

const FILENAME = "filelist.json"

type Cache struct {
	Options *Options       `json:"options"`
	Files   map[string]int `json:"files"`
}

func (c *Cache) Check(options *Options) bool {
	return c.Options.Forced == options.Forced &&
		c.Options.SkipSrt == options.SkipSrt &&
		c.Options.StripFormatting == options.StripFormatting &&
		reflect.DeepEqual(c.Options.Langs, options.Langs) &&
		reflect.DeepEqual(c.Options.ForcedTitles, options.ForcedTitles)
}

func (c *Cache) Validate(options *Options) {
	if !c.Check(options) {
		c.Options = options
		clear(c.Files)
	}
}

func (c *Cache) Save() {
	fullName := filepath.Join(c.Options.DataDir, FILENAME)
	bytes, _ := json.MarshalIndent(c, "", "  ")
	os.WriteFile(fullName, bytes, os.ModePerm)
}

func ReadCache(options *Options) *Cache {
	fullName := filepath.Join(options.DataDir, FILENAME)
	cache := &Cache{Files: map[string]int{}}

	bytes, err := os.ReadFile(fullName)
	if err == nil {
		json.Unmarshal(bytes, cache)
		cache.Validate(options)
		log.Println("Cache size:", len(cache.Files))
	} else {
		log.Println("Cache (re)initialized")
	}

	return cache
}
