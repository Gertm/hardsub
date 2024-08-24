package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

var ffprobe_cmdline string = "ffprobe -v quiet -print_format json -show_format -show_streams -show_chapters "

func GetFFprobeInfo(filename string) (*VideoProbeInfo, error) {
	output, err := OutputBytesForCommand(ffprobe_cmdline + filename)
	if err != nil {
		return nil, err
	}
	var result VideoProbeInfo
	err = json.Unmarshal(output, &result)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal ffprobe output: %w", err)
	}
	return &result, nil
}

type VideoProbeInfo struct {
	Streams  []Stream  `json:"streams"`
	Format   Format    `json:"format"`
	Chapters []Chapter `json:"chapters"`
}

type Subtitle struct {
	Index    int
	Title    string
	Language string
}

func (vpi VideoProbeInfo) ShowSubtitles() {
	for _, stream := range vpi.Streams {
		if stream.IsSubtitle() {
			fmt.Printf("CodecName: %s, CodecTagString: %s, CodecTag: %s\n", stream.CodecName, stream.CodecTagString, stream.CodecTag)
			if lang, ok := stream.Tags["language"]; ok {
				fmt.Println("  Language:", lang)
			}
			if title, ok := stream.Tags["title"]; ok {
				fmt.Println("Subs title:", title)
			}
		}
	}
}

func (vpi VideoProbeInfo) ShowChapters() {
	fmt.Printf("Chapter count: %d\n", len(vpi.Chapters))
	for _, chapter := range vpi.Chapters {
		fmt.Println(chapter)
	}
}

type Stream struct {
	Index              int    `json:"index"`
	CodecName          string `json:"codec_name"`
	CodecLongName      string `json:"codec_long_name"`
	Profile            string `json:"profile"`
	CodecType          string `json:"codec_type"`
	CodecTagString     string `json:"codec_tag_string"`
	CodecTag           string `json:"codec_tag"`
	Width              int    `json:"width"`
	Height             int    `json:"height"`
	CodedWidth         int    `json:"coded_width"`
	CodedHeight        int    `json:"coded_height"`
	ClosedCaptions     int    `json:"closed_captions"`
	FilmGrain          int    `json:"film_grain"`
	HasBFrames         int    `json:"has_b_frames"`
	SampleAspectRatio  string `json:"sample_aspect_ratio"`
	DisplayAspectRatio string `json:"display_aspect_ratio"`
	PixFmt             string `json:"pix_fmt"`
	Level              int    `json:"level"`
	ColorRange         string `json:"color_range"`
	ColorSpace         string `json:"color_space"`
	ColorTransfer      string `json:"color_transfer"`
	ColorPrimaries     string `json:"color_primaries"`
	ChromaLocation     string `json:"chroma_location"`
	FieldOrder         string `json:"field_order"`
	Refs               int    `json:"refs"`
	IsAvc              string `json:"is_avc"`
	NalLengthSize      string `json:"nal_length_size"`
	RFrameRate         string `json:"r_frame_rate"`
	AvgFrameRate       string `json:"avg_frame_rate"`
	TimeBase           string `json:"time_base"`
	StartPts           int    `json:"start_pts"`
	StartTime          string `json:"start_time"`
	BitsPerRawSample   string `json:"bits_per_raw_sample"`
	ExtradataSize      int    `json:"extradata_size"`
	Disposition        struct {
		Default         int `json:"default"`
		Dub             int `json:"dub"`
		Original        int `json:"original"`
		Comment         int `json:"comment"`
		Lyrics          int `json:"lyrics"`
		Karaoke         int `json:"karaoke"`
		Forced          int `json:"forced"`
		HearingImpaired int `json:"hearing_impaired"`
		VisualImpaired  int `json:"visual_impaired"`
		CleanEffects    int `json:"clean_effects"`
		AttachedPic     int `json:"attached_pic"`
		TimedThumbnails int `json:"timed_thumbnails"`
		Captions        int `json:"captions"`
		Descriptions    int `json:"descriptions"`
		Metadata        int `json:"metadata"`
		Dependent       int `json:"dependent"`
		StillImage      int `json:"still_image"`
	} `json:"disposition"`
	Tags map[string]string `json:"tags"`
}

func (stream Stream) GetLanguage() (string, error) {
	for k, v := range stream.Tags {
		if strings.TrimSpace(strings.ToLower(k)) == "language" {
			return v, nil
		}
	}
	return "", fmt.Errorf("no language tag found in this stream")
}

func (stream Stream) IsSubtitle() bool {
	return stream.CodecType == "subtitle"
}

func (stream Stream) IsVideo() bool {
	return stream.CodecType == "video"
}

func (stream Stream) IsAudio() bool {
	return stream.CodecType == "video"
}

type Format struct {
	Filename       string            `json:"filename"`
	NbStreams      int               `json:"nb_streams"`
	NbPrograms     int               `json:"nb_programs"`
	FormatName     string            `json:"format_name"`
	FormatLongName string            `json:"format_long_name"`
	StartTime      string            `json:"start_time"`
	Duration       string            `json:"duration"`
	Size           string            `json:"size"`
	BitRate        string            `json:"bit_rate"`
	ProbeScore     int               `json:"probe_score"`
	Tags           map[string]string `json:"tags"`
}

type Chapter struct {
	ID        int               `json:"id"`
	TimeBase  string            `json:"time_base"`
	Start     int               `json:"start"`
	StartTime string            `json:"start_time"`
	End       int64             `json:"end"`
	EndTime   string            `json:"end_time"`
	Tags      map[string]string `json:"tags"`
}
