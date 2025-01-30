package main

import "github.com/microcosm-cc/bluemonday"

func Strip(data string) string {
	return bluemonday.StripTagsPolicy().Sanitize(data)
}
