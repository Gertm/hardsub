package main

import (
	"fmt"
	"testing"
)

func TestGetFFprobeInfo(t *testing.T) {
	got, err := GetFFprobeInfo("/home/gert/Videos/S02E01-A_New_Roar_3289D5B8_.mkv")
	if err != nil {
		t.Fail()
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
}
