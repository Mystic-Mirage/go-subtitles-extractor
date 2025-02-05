package main

import "github.com/martinlindhe/subtitles"

func Strip(data string) string {
	res, err := subtitles.NewFromSRT(data)
	if err != nil {
		return data
	}

	res.FilterCaptions("html")
	return res.AsSRT()
}
