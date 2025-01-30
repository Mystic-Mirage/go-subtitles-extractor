package main

import (
	"encoding/json"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

type Tags struct {
	Language string `json:"language"`
	Title    string `json:"title"`
}

type Disposition struct {
	Forced          int `json:"forced"`
	HearingImpaired int `json:"hearing_impaired"`
}

type Stream struct {
	Index       int          `json:"index"`
	Tags        *Tags        `json:"tags"`
	Disposition *Disposition `json:"disposition"`
	CodecName   string       `json:"codec_name"`
	Width       int          `json:"width"`
}

func (s *Stream) Subtitles(fileName string, forcedTitles []string) *Subtitles {
	language := s.Tags.Language
	if language == "" {
		language = "und"
	}

	title := strings.ToLower(s.Tags.Title)
	forced := s.Disposition.Forced > 0
	if !forced {
		for _, forcedTitle := range forcedTitles {
			forced = strings.Contains(title, forcedTitle)
			if forced {
				break
			}
		}
	}
	sdh := s.Disposition.HearingImpaired > 0 || strings.Contains(title, "sdh")

	return &Subtitles{
		FileName: fileName,
		Index:    s.Index,
		Language: language,
		Codec:    s.CodecName,
		Bitmap:   s.Width > 0,
		Forced:   forced,
		Sdh:      sdh,
	}
}

type ProbeData struct {
	Streams []*Stream `json:"streams"`
}

type Subtitles struct {
	FileName string
	Index    int
	Language string
	Codec    string
	Bitmap   bool
	Forced   bool
	Sdh      bool
}

func (s *Subtitles) Path() string {
	suffix := "." + s.Language

	if s.Sdh {
		suffix += ".sdh"
	}

	if s.Forced {
		suffix += ".forced"
	}

	suffix += EXT

	return strings.TrimSuffix(s.FileName, filepath.Ext(s.FileName)) + suffix
}

func GetSubtitles(fileName string, forcedTitles []string) []*Subtitles {
	subtitles := []*Subtitles{}

	stdErr := strings.Builder{}
	command := exec.Command(
		"ffprobe",
		"-v",
		"error",
		"-print_format",
		"json",
		"-show_streams",
		"-select_streams",
		"s",
		fileName,
	)
	command.Stderr = &stdErr

	output, err := command.Output()
	if err != nil {
		log.Print("FFPROBE error:", stdErr.String())
		return subtitles
	}

	probeData := &ProbeData{}
	json.Unmarshal(output, probeData)

	for _, stream := range probeData.Streams {
		subtitles = append(subtitles, stream.Subtitles(fileName, forcedTitles))
	}

	return subtitles
}
