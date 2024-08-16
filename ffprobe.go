package main

import "time"

var ffprobe_cmdline string = "ffprobe -v quiet -print_format json -show_format -show_streams -show_format"

func GetFFprobeInfo(filename string) FfprobeOutput {

}

type FfprobeOutput struct {
	Streams  []Stream  `json:"streams"`
	Format   Format    `json:"format"`
	Chapters []Chapter `json:"chapters"`
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
	Tags Tags `json:"tags"`
}

type Tags struct {
	BPSEng                      string `json:"BPS-eng"`
	DURATIONEng                 string `json:"DURATION-eng"`
	NUMBEROFFRAMESEng           string `json:"NUMBER_OF_FRAMES-eng"`
	NUMBEROFBYTESEng            string `json:"NUMBER_OF_BYTES-eng"`
	STATISTICSWRITINGAPPEng     string `json:"_STATISTICS_WRITING_APP-eng"`
	STATISTICSWRITINGDATEUTCEng string `json:"_STATISTICS_WRITING_DATE_UTC-eng"`
	STATISTICSTAGSEng           string `json:"_STATISTICS_TAGS-eng"`
}

type Format struct {
	Filename       string `json:"filename"`
	NbStreams      int    `json:"nb_streams"`
	NbPrograms     int    `json:"nb_programs"`
	FormatName     string `json:"format_name"`
	FormatLongName string `json:"format_long_name"`
	StartTime      string `json:"start_time"`
	Duration       string `json:"duration"`
	Size           string `json:"size"`
	BitRate        string `json:"bit_rate"`
	ProbeScore     int    `json:"probe_score"`
	Tags           struct {
		Encoder      string    `json:"encoder"`
		CreationTime time.Time `json:"creation_time"`
	} `json:"tags"`
}

type Chapter struct {
	ID        int    `json:"id"`
	TimeBase  string `json:"time_base"`
	Start     int    `json:"start"`
	StartTime string `json:"start_time"`
	End       int64  `json:"end"`
	EndTime   string `json:"end_time"`
	Tags      struct {
		Title string `json:"title"`
	} `json:"tags"`
}
