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
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/sanity-io/litter"
)

type IntroBoundaries struct {
	Begin string
	End   string
}

type Arguments struct {
	SourceDirectory string `koanf:"sourcedir"`
	CutStart        string `koanf:"cutstart"`
	CutEnd          string `koanf:"cutend"`
	OnlyCut         bool   `koanf:"onlycut"`
	DumpFramesAt    string `koanf:"dumpframesat"`
	File            string `koanf:"file"`
	ForceAudioTrack int    `koanf:"forceaudiotrack"`
	ForceSubsTrack  int    `koanf:"forcesubstrack"`
}

type Config struct {
	AudioLang          string `koanf:"audiolang" toml:"audiolang" comment:"The audio language you want to use in the ouput video. (IETF language tag)"`
	SubsLang           string `koanf:"subslang" toml:"subslang" comment:"The subs language you want to use. (IETF language tag)"`
	SubsName           string `koanf:"subsname" toml:"subsname" comment:"What the subtitle trackname needs to contains."`
	TargetDirectory    string `koanf:"targetdir" toml:"targetdir" comment:"Where to put the converted videos."`
	OriginalsDirectory string `koanf:"originalsdir" toml:"originalsdir" comment:"Where to move the original files to."`
	H26xTune           string `koanf:"h26xtune" toml:"h26xtune" comment:"The tuning to use for h26x encoding. (film/animation/fastdecode/zerolatency/none)"`
	H26xPreset         string `koanf:"h26xpreset" toml:"h26xpreset" comment:"The preset to use for h26x encoding. (fast/medium/slow/etc..)"`
	PostCmd            string `koanf:"postcmd" toml:"postcmd" comment:"The command to run on completion. Use %%o for the output filename."`
	PostSubExtract     string `koanf:"postsubextract" toml:"postsubextract" comment:"The command to run after sub extraction, before conversion. Use %%s for subs filename."`
	Extension          string `koanf:"extension" toml:"extension" comment:"Look for files of this extension to convert. (You really want to set this to mkv)"`
	RemoveWords        string `koanf:"removewords" toml:"removewords" comment:"When detoxing, remove the words in the comma separated value you specify."`
	filesToConvert     []fs.DirEntry
	Crf                int                        `koanf:"crf" toml:"crf" comment:"Constant Rate Factor setting for ffmpeg."`
	ExtractFonts       bool                       `koanf:"extractfonts" toml:"extractfonts" comment:"Extract the fonts from the mkv to use them in the hardcoding."`
	FirstOnly          bool                       `koanf:"firstonly" toml:"firstonly" comment:"Only convert the first file. (For testing purposes)"`
	Mkv                bool                       `koanf:"mkv" toml:"mkv" comment:"Make MKV files instead of MP4 files."`
	H265               bool                       `koanf:"h265" toml:"h265" comment:"Use H265 encoding. Check if your CPU can do H265 encoding first, or this will be very slow."`
	KeepSubs           bool                       `koanf:"keepsubs" toml:"keepsubs" comment:"Keep subs in the directory after conversion instead of deleting them."`
	CleanupSubs        bool                       `koanf:"cleanupsubs" toml:"cleanupsubs" comment:"Clean up the subtitles (in the case of srt) to make them render better. Sometimes they render too big, use this in that case."`
	Verbose            bool                       `koanf:"verbose" toml:"verbose" comment:"Give more output about what's going on."`
	ForOldDevices      bool                       `koanf:"forolddevices" toml:"forolddevices" comment:"Use ffmpeg flags to get widest compatibility. (yuv stuff)"`
	FastVersion        bool                       `koanf:"fastversion" toml:"fastversion" comment:"Do a second and third pass, making a video at 1.5x the speed."`
	KeepSlowVersion    bool                       `koanf:"keepslowversion" toml:"keepslowversion" comment:"When making a fast version, don't delete the slow one."`
	Detox              bool                       `koanf:"detox" toml:"detox" comment:"Remove all 'weird' characters from the filename. (you want this)"`
	WatchForFiles      bool                       `koanf:"watchforfiles" toml:"watchforfiles" comment:"Watch for files in the directory and convert them as they appear."`
	IntroFrames        map[string]IntroBoundaries `koanf:"introframes" toml:"introframes" comment:"The locations of the intro beginning and ending frames for specific series."`
	arguments          Arguments                  `koanf:"arguments"`
}

func (c Config) FfmpegParametersForCutting(inputFile, outputFile string) string {
	sb := strings.Builder{}
	sb.WriteString("-y -hide_banner -loglevel error -stats -i ")
	sb.WriteString(inputFile)
	sb.WriteString(" -c:a aac")
	videoCodec := "libx264"
	if config.H265 {
		videoCodec = "libx265"
	}
	sb.WriteString(" -c:v " + videoCodec)
	sb.WriteString(fmt.Sprintf(" -crf %d", config.Crf))
	h26xTune := ""
	if config.H26xTune == "none" {
		h26xTune = ""
	} else {
		h26xTune = " -tune " + config.H26xTune + " "
	}
	sb.WriteString(h26xTune)
	sb.WriteString(" -preset " + config.H26xPreset)

	if config.ForOldDevices {
		sb.WriteString(" -profile:v baseline -level 3.0 -pix_fmt yuv420p -ac 2 -b:a 128k -movflags faststart ")
	}
	sb.WriteString(outputFile)
	// convertCmd := fmt.Sprintf("-y -hide_banner -loglevel error -stats -i %s -map 0:%d -map 0:%d -vf subtitles=%s -c:a %s -c:v %s -crf %d -preset %s %s%s%s",
	// 	videofile, output.VideoTrack, output.AudioTrack, subsfile, audioCodec, videoCodec, config.Crf, config.H26xPreset, h26xTune, oldDevices, outputFile)
	return sb.String()
}

func DefaultConfig() Config {
	return Config{
		AudioLang:          "ja",
		SubsLang:           "en",
		SubsName:           "subtitles",
		TargetDirectory:    "converted",
		OriginalsDirectory: "originals",
		H26xTune:           "animation",
		H26xPreset:         "fast",
		PostCmd:            "",
		PostSubExtract:     "",
		Extension:          "mkv",
		RemoveWords:        "SubsPlease,EMBER",
		Crf:                18,
		ExtractFonts:       true,
		FirstOnly:          false,
		Mkv:                false,
		H265:               false,
		KeepSubs:           false,
		CleanupSubs:        false,
		Verbose:            false,
		ForOldDevices:      false,
		FastVersion:        false,
		KeepSlowVersion:    false,
		Detox:              true,
		WatchForFiles:      false,
	}
}

func PrepareDirectoryForConversion(config *Config) {
	if config.Detox {
		fmt.Print("Detoxing directory...")
		detoxWords := strings.Split(config.RemoveWords, ",")
		if err := DetoxMkvsInDirectory(config.arguments.SourceDirectory, detoxWords...); err != nil {
			LogError("detoxing directory failed: %s %v", config.arguments.SourceDirectory, err)
			os.Exit(1)
		}
		fmt.Print("done.\n")
	}
	files, err := os.ReadDir(config.arguments.SourceDirectory)
	if err != nil {
		LogError("cannot read directory: %s %v", config.arguments.SourceDirectory, err)
	}
	config.filesToConvert = files
	if config.Verbose {
		litter.Dump(config)
	}
}
