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
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"

	"github.com/sanity-io/litter"
)

type Config struct {
	AudioLang       string
	ScriptPath      string
	SubsLang        string
	SubsName        string
	TargetFolder    string
	OriginalsFolder string
	H26xTune        string
	H26xPreset      string
	PostCmd         string
	PostSubExtract  string
	Extension       string
	RemoveWords     string
	SourceFolder    string
	FilesToConvert  []fs.DirEntry
	Crf             int
	ScpPort         int
	ExtractFonts    bool
	FirstOnly       bool
	Mkv             bool
	H265            bool
	RunDirectly     bool
	KeepSubs        bool
	CleanupSubs     bool
	Verbose         bool
	ForOldDevices   bool
	FastVersion     bool
	KeepSlowVersion bool
	Detox           bool
	CutStart        string
	CutEnd          string
	OnlyCut         bool
	DumpFramesAt    string
	File            string
	ForceAudioTrack int
	ForceSubsTrack  int
	WathForFiles    bool
}

func getConfigurationFromArguments() Config {
	config := Config{}
	workdir, _ := os.Getwd()
	proposedConvertedDir := path.Join(workdir, "converted")
	flag.StringVar(&config.SourceFolder, "folder", workdir, "The folder to convert all mkvs in. Defaults to the working directory.")
	flag.BoolVar(&config.Detox, "detox", true, "Detox the mkv files in the directory first.")
	flag.StringVar(&config.RemoveWords, "removewords", "SubsPlease,EMBER", "Remove the words in the comma separated value you specify.")
	flag.StringVar(&config.Extension, "ext", "mkv", "Look for files of this extension to convert.")
	flag.StringVar(&config.SubsLang, "subslang", "en", "The subs language you want to use. (IETF language tag)")
	flag.StringVar(&config.SubsName, "subsname", "subtitles", "What the subs name needs to contains.")
	flag.StringVar(&config.AudioLang, "audiolang", "ja", "The audio language you want to use. (IETF language tag)")
	flag.StringVar(&config.TargetFolder, "outputfolder", proposedConvertedDir, "The folder to put the converted videos in.")
	flag.StringVar(&config.OriginalsFolder, "move-originals-to", "originals", "The alternative folder you want the originals moved to.")
	flag.BoolVar(&config.KeepSubs, "keepsubs", false, "Keep subs in the folder after conversion instead of deleting them.")
	flag.BoolVar(&config.CleanupSubs, "cleansubs", false, "Clean up the subtitles (in the case of srt) to make them render better.")
	flag.BoolVar(&config.ExtractFonts, "extract-fonts", true, "Extract the fonts from the mkv to use them in the hardcoding.")
	flag.BoolVar(&config.Verbose, "v", false, "Show verbose output.")
	flag.IntVar(&config.Crf, "crf", 18, "Constant Rate Factor setting for ffmpeg.")
	flag.BoolVar(&config.FirstOnly, "first-only", false, "Only convert the first file. (For testing purposes)")
	flag.BoolVar(&config.ForOldDevices, "for-old-devices", false, "Use ffmpeg flags to get widest compatibility.")
	flag.BoolVar(&config.FastVersion, "fastversion", false, "Do a second and third pass, making a video at 1.5x the speed.")
	flag.BoolVar(&config.KeepSlowVersion, "keep-slow", false, "In case you're making fast versions, keep the slow versions as well.")
	flag.BoolVar(&config.Mkv, "mkv", false, "Make MKV files instead of MP4 files.")
	flag.StringVar(&config.PostCmd, "postcmd", "", "The command to run on completion. Use %%o for the output filename.")
	flag.StringVar(&config.PostSubExtract, "postsubextract", "", "The command to run after sub extraction, before conversion. Use %%s for subs filename.")
	flag.StringVar(&config.H26xTune, "h26x-tune", "animation", "The tuning to use for h26x encoding. (film/animation/fastdecode/zerolatency/none)")
	flag.StringVar(&config.H26xPreset, "h26x-preset", "fast", "The preset to use for h26x encoding. (fast/medium/slow/etc..)")
	flag.BoolVar(&config.H265, "h265", false, "Use H265 encoding.")
	flag.StringVar(&config.File, "file", "", "The specific file to operate on for cutting and frame dumping.")
	flag.BoolVar(&config.OnlyCut, "onlycut", false, "Only cut, don't convert.")
	flag.BoolVar(&config.WathForFiles, "watch", false, "Watch for files in the folder and convert them as they appear.")
	flag.StringVar(&config.CutStart, "cutstart", "", "A jpg of the frame to look for in the video to determine the start of the fragment to cut out.")
	flag.StringVar(&config.CutEnd, "cutend", "", "A jpg of the frame to look for in the video to determine the end of the fragment to cut out.")
	flag.StringVar(&config.DumpFramesAt, "dumpframesat", "", "A comma-separated list of timestamps you want to make jpg dumps for.")
	flag.IntVar(&config.ForceAudioTrack, "force-audio-track", -1, "Force the audio track to use. (for example: 4)")
	flag.IntVar(&config.ForceSubsTrack, "force-subs-track", -1, "Force the subs track to use. (for example: 3)")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `
 This tool is for converting mkv files to an mp4 that has the subs hardcoded into the video.
 For when your television doesn't like the subs in the mkv and doesn't want to display them, or
 displays them too quickly...
 
 Available options:
 `)
		flag.PrintDefaults()
	}

	flag.Parse()
	VERBOSE = config.Verbose
	// TODO: this should not be in this function, it doesn't belong here.
	// this should be done before we get which files we need to convert though.
	if config.Detox {
		fmt.Print("Detoxing folder...")
		detoxWords := strings.Split(config.RemoveWords, ",")
		if err := DetoxMkvsInFolder(config.SourceFolder, detoxWords...); err != nil {
			log.Fatal("Cannot detox folder?!", err)
		}
		fmt.Print("done.\n")
	}
	files, err := os.ReadDir(config.SourceFolder)
	if err != nil {
		log.Fatal(err)
	}
	config.FilesToConvert = files
	if VERBOSE {
		litter.Dump(config)
	}
	return config
}

func PrepareFolderForConversion(config *Config) {
	if config.Detox {
		fmt.Print("Detoxing folder...")
		detoxWords := strings.Split(config.RemoveWords, ",")
		if err := DetoxMkvsInFolder(config.SourceFolder, detoxWords...); err != nil {
			log.Fatal("Cannot detox folder?!", err)
		}
		fmt.Print("done.\n")
	}
	files, err := os.ReadDir(config.SourceFolder)
	if err != nil {
		log.Fatal(err)
	}
	config.FilesToConvert = files
	if VERBOSE {
		litter.Dump(config)
	}
}
