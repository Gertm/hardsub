package main

import (
	"fmt"
	"testing"

	"github.com/sanity-io/litter"
)

func TestGetFFprobeInfo(t *testing.T) {
	got, err := GetFFprobeInfo("/home/gert/Videos/S02E01-A_New_Roar_3289D5B8_.mkv")
	if err != nil {
		t.Fail()
	}
	litter.Dump(got)
	for _, stream := range got.Streams {
		fmt.Printf("Stream: %s\n", stream.CodecType)
	}
}
