/*
		Copyright (C) 2023 Gert Meulyzer

	    This program is free software: you can redistribute it and/or modify
	    it under the terms of the GNU General Public License as published by
	    the Free Software Foundation, either version 3 of the License, or
	    (at your option) any later version.

	    This program is distributed in the hope that it will be useful,
	    but WITHOUT ANY WARRANTY; without even the implied warranty of
	    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	    GNU General Public License for more details.

	    You should have received a copy of the GNU General Public License
	    along with this program.  If not, see <https://www.gnu.org/licenses/>.
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
	ScpHost         string
	ScpUser         string
	ScpPrivKeyPath  string
	ScpTargetDir    string
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
	flag.StringVar(&config.ScpHost, "scpHost", "", "The hostname of the server you want to SCP to after successfull conversion.")
	flag.StringVar(&config.ScpUser, "scpUser", "", "The username when doing SCP.")
	flag.IntVar(&config.ScpPort, "scpPort", 22, "The port to use when doing SCP.")
	flag.StringVar(&config.ScpPrivKeyPath, "scpPrivKeyPath", "", "The location of your private key file for SCP.")
	flag.StringVar(&config.ScpTargetDir, "scpTargetDir", "", "The remote folder you want to SCP into.")

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
