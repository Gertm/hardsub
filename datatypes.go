/*
Copyright 2023 Gert Meulyzer

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	// "os"
	"time"
)

type (
	SubsType     int
	MappedTracks []int
)

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
	Container struct {
		Type       string `json:"type"`
		Properties struct {
			DateLocal             time.Time `json:"date_local"`
			DateUtc               time.Time `json:"date_utc"`
			MuxingApplication     string    `json:"muxing_application"`
			SegmentUID            string    `json:"segment_uid"`
			WritingApplication    string    `json:"writing_application"`
			ContainerType         int       `json:"container_type"`
			Duration              int64     `json:"duration"`
			IsProvidingTimestamps bool      `json:"is_providing_timestamps"`
		} `json:"properties"`
		Recognized bool `json:"recognized"`
		Supported  bool `json:"supported"`
	} `json:"container"`
	FileName    string        `json:"file_name"`
	Attachments []interface{} `json:"attachments"`
	Chapters    []struct {
		NumEntries int `json:"num_entries"`
	} `json:"chapters"`
	Errors     []interface{} `json:"errors"`
	GlobalTags []interface{} `json:"global_tags"`
	TrackTags  []interface{} `json:"track_tags"`
	Tracks     []struct {
		Codec       string `json:"codec"`
		Type        string `json:"type"`
		Properties3 struct {
			CodecID                   string `json:"codec_id"`
			ContentEncodingAlgorithms string `json:"content_encoding_algorithms"`
			Language                  string `json:"language"`
			LanguageIetf              string `json:"language_ietf"`
			TrackName                 string `json:"track_name"`
			CodecPrivateLength        int    `json:"codec_private_length"`
			Number                    int    `json:"number"`
			UID                       int64  `json:"uid"`
			DefaultTrack              bool   `json:"default_track"`
			EnabledTrack              bool   `json:"enabled_track"`
			ForcedTrack               bool   `json:"forced_track"`
		} `json:"properties,omitempty"`
		Properties2 struct {
			CodecID                   string `json:"codec_id"`
			ContentEncodingAlgorithms string `json:"content_encoding_algorithms"`
			Language                  string `json:"language"`
			LanguageIetf              string `json:"language_ietf"`
			CodecPrivateLength        int    `json:"codec_private_length"`
			MinimumTimestamp          int    `json:"minimum_timestamp"`
			Number                    int    `json:"number"`
			UID                       int64  `json:"uid"`
			DefaultTrack              bool   `json:"default_track"`
			EnabledTrack              bool   `json:"enabled_track"`
			ForcedTrack               bool   `json:"forced_track"`
		} `json:"properties,omitempty"`
		Properties struct {
			CodecID            string `json:"codec_id"`
			CodecPrivateData   string `json:"codec_private_data"`
			DisplayDimensions  string `json:"display_dimensions"`
			Language           string `json:"language"`
			LanguageIetf       string `json:"language_ietf"`
			Packetizer         string `json:"packetizer"`
			PixelDimensions    string `json:"pixel_dimensions"`
			CodecPrivateLength int    `json:"codec_private_length"`
			DefaultDuration    int    `json:"default_duration"`
			DisplayUnit        int    `json:"display_unit"`
			MinimumTimestamp   int    `json:"minimum_timestamp"`
			Number             int    `json:"number"`
			UID                int64  `json:"uid"`
			DefaultTrack       bool   `json:"default_track"`
			EnabledTrack       bool   `json:"enabled_track"`
			ForcedTrack        bool   `json:"forced_track"`
		} `json:"properties,omitempty"`
		Properties0 struct {
			CodecID                string `json:"codec_id"`
			Language               string `json:"language"`
			LanguageIetf           string `json:"language_ietf"`
			AudioChannels          int    `json:"audio_channels"`
			AudioSamplingFrequency int    `json:"audio_sampling_frequency"`
			CodecPrivateLength     int    `json:"codec_private_length"`
			DefaultDuration        int    `json:"default_duration"`
			MinimumTimestamp       int    `json:"minimum_timestamp"`
			Number                 int    `json:"number"`
			UID                    int    `json:"uid"`
			DefaultTrack           bool   `json:"default_track"`
			EnabledTrack           bool   `json:"enabled_track"`
			ForcedTrack            bool   `json:"forced_track"`
		} `json:"properties,omitempty"`
		Properties1 struct {
			CodecID                string `json:"codec_id"`
			Language               string `json:"language"`
			LanguageIetf           string `json:"language_ietf"`
			AudioChannels          int    `json:"audio_channels"`
			AudioSamplingFrequency int    `json:"audio_sampling_frequency"`
			CodecPrivateLength     int    `json:"codec_private_length"`
			DefaultDuration        int    `json:"default_duration"`
			MinimumTimestamp       int    `json:"minimum_timestamp"`
			Number                 int    `json:"number"`
			UID                    int    `json:"uid"`
			DefaultTrack           bool   `json:"default_track"`
			EnabledTrack           bool   `json:"enabled_track"`
			ForcedTrack            bool   `json:"forced_track"`
		} `json:"properties,omitempty"`
		ID int `json:"id"`
	} `json:"tracks"`
	Warnings                    []interface{} `json:"warnings"`
	IdentificationFormatVersion int           `json:"identification_format_version"`
}
