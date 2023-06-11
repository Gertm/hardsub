package main

import (
	"io/fs"
	"os"
	"time"
)

type Config struct {
	AudioLang       string
	ScriptPath      string
	SubsLang        string
	SubsName        string
	ScriptFile      *os.File
	TargetFolder    string
	OriginalsFolder string
	ExtractFonts    bool
	FirstOnly       bool
	Mp4             bool
	H26xTune        string
	H26xPreset      string
	H265            bool
	PostCmd         string
	PostSubExtract  string
	RunDirectly     bool
	FilesToConvert  []fs.DirEntry
	KeepSubs        bool
	CleanupSubs     bool
	Verbose         bool
	Crf             int
	ForOldDevices   bool
	Extension       string
	DoubleSpeed     bool
}

type SubsType int
type MappedTracks []int

const (
	SRT SubsType = iota
	SSA_ASS
	PICTURE
)

type SelectedTracks struct {
	VideoTrack   int
	AudioTrack   int
	SubsTrack    int
	SubtitleType SubsType
}

func (lst MappedTracks) contains(i int) bool {
	for _, a := range lst {
		if a == i {
			return true
		}
	}
	return false
}

type MkvMergeOutput struct {
	Attachments []interface{} `json:"attachments"`
	Chapters    []struct {
		NumEntries int `json:"num_entries"`
	} `json:"chapters"`
	Container struct {
		Properties struct {
			ContainerType         int       `json:"container_type"`
			DateLocal             time.Time `json:"date_local"`
			DateUtc               time.Time `json:"date_utc"`
			Duration              int64     `json:"duration"`
			IsProvidingTimestamps bool      `json:"is_providing_timestamps"`
			MuxingApplication     string    `json:"muxing_application"`
			SegmentUID            string    `json:"segment_uid"`
			WritingApplication    string    `json:"writing_application"`
		} `json:"properties"`
		Recognized bool   `json:"recognized"`
		Supported  bool   `json:"supported"`
		Type       string `json:"type"`
	} `json:"container"`
	Errors                      []interface{} `json:"errors"`
	FileName                    string        `json:"file_name"`
	GlobalTags                  []interface{} `json:"global_tags"`
	IdentificationFormatVersion int           `json:"identification_format_version"`
	TrackTags                   []interface{} `json:"track_tags"`
	Tracks                      []struct {
		Codec      string `json:"codec"`
		ID         int    `json:"id"`
		Properties struct {
			CodecID            string `json:"codec_id"`
			CodecPrivateData   string `json:"codec_private_data"`
			CodecPrivateLength int    `json:"codec_private_length"`
			DefaultDuration    int    `json:"default_duration"`
			DefaultTrack       bool   `json:"default_track"`
			DisplayDimensions  string `json:"display_dimensions"`
			DisplayUnit        int    `json:"display_unit"`
			EnabledTrack       bool   `json:"enabled_track"`
			ForcedTrack        bool   `json:"forced_track"`
			Language           string `json:"language"`
			LanguageIetf       string `json:"language_ietf"`
			MinimumTimestamp   int    `json:"minimum_timestamp"`
			Number             int    `json:"number"`
			Packetizer         string `json:"packetizer"`
			PixelDimensions    string `json:"pixel_dimensions"`
			UID                int64  `json:"uid"`
		} `json:"properties,omitempty"`
		Type        string `json:"type"`
		Properties0 struct {
			AudioChannels          int    `json:"audio_channels"`
			AudioSamplingFrequency int    `json:"audio_sampling_frequency"`
			CodecID                string `json:"codec_id"`
			CodecPrivateLength     int    `json:"codec_private_length"`
			DefaultDuration        int    `json:"default_duration"`
			DefaultTrack           bool   `json:"default_track"`
			EnabledTrack           bool   `json:"enabled_track"`
			ForcedTrack            bool   `json:"forced_track"`
			Language               string `json:"language"`
			LanguageIetf           string `json:"language_ietf"`
			MinimumTimestamp       int    `json:"minimum_timestamp"`
			Number                 int    `json:"number"`
			UID                    int    `json:"uid"`
		} `json:"properties,omitempty"`
		Properties1 struct {
			AudioChannels          int    `json:"audio_channels"`
			AudioSamplingFrequency int    `json:"audio_sampling_frequency"`
			CodecID                string `json:"codec_id"`
			CodecPrivateLength     int    `json:"codec_private_length"`
			DefaultDuration        int    `json:"default_duration"`
			DefaultTrack           bool   `json:"default_track"`
			EnabledTrack           bool   `json:"enabled_track"`
			ForcedTrack            bool   `json:"forced_track"`
			Language               string `json:"language"`
			LanguageIetf           string `json:"language_ietf"`
			MinimumTimestamp       int    `json:"minimum_timestamp"`
			Number                 int    `json:"number"`
			UID                    int    `json:"uid"`
		} `json:"properties,omitempty"`
		Properties2 struct {
			CodecID                   string `json:"codec_id"`
			CodecPrivateLength        int    `json:"codec_private_length"`
			ContentEncodingAlgorithms string `json:"content_encoding_algorithms"`
			DefaultTrack              bool   `json:"default_track"`
			EnabledTrack              bool   `json:"enabled_track"`
			ForcedTrack               bool   `json:"forced_track"`
			Language                  string `json:"language"`
			LanguageIetf              string `json:"language_ietf"`
			MinimumTimestamp          int    `json:"minimum_timestamp"`
			Number                    int    `json:"number"`
			UID                       int64  `json:"uid"`
		} `json:"properties,omitempty"`
		Properties3 struct {
			CodecID                   string `json:"codec_id"`
			CodecPrivateLength        int    `json:"codec_private_length"`
			ContentEncodingAlgorithms string `json:"content_encoding_algorithms"`
			DefaultTrack              bool   `json:"default_track"`
			EnabledTrack              bool   `json:"enabled_track"`
			ForcedTrack               bool   `json:"forced_track"`
			Language                  string `json:"language"`
			LanguageIetf              string `json:"language_ietf"`
			Number                    int    `json:"number"`
			TrackName                 string `json:"track_name"`
			UID                       int64  `json:"uid"`
		} `json:"properties,omitempty"`
	} `json:"tracks"`
	Warnings []interface{} `json:"warnings"`
}
