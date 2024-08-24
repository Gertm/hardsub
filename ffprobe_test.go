package main

import (
	"fmt"
	"testing"
)

func TestGetFFprobeInfo(t *testing.T) {
	got, err := GetFFprobeInfo("testvideo2.mkv")
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}
	for _, stream := range got.Streams {
		fmt.Printf("Stream: %s\n", stream.CodecType)
		lang, err := stream.GetLanguage()
		if err != nil {
			continue
		}
		fmt.Printf("Stream language: %s\n", lang)
	}
	got.ShowSubtitles()
	got.ShowChapters()
}
