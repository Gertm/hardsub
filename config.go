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
	FastVersion     bool
	KeepSlowVersion bool
}

func getConfigurationFromArguments() Config {
	workdir, _ := os.Getwd()
	proposedConvertedDir := path.Join(workdir, "converted")
	folder := flag.String("folder", workdir, "The folder to convert all mkvs in. Defaults to the working directory.")
	detox := flag.Bool("detox", true, "Detox the mkv files in the directory first.")
	detoxRemove := flag.String("removewords", "SubsPlease,EMBER", "Remove the words in the comma separated value you specify.")
	extension := flag.String("ext", "mkv", "Look for files of this extension to convert.")
	subslang := flag.String("subslang", "en", "The subs language you want to use. (IETF language tag)")
	subsname := flag.String("subsname", "subtitles", "What the subs name needs to contains.")
	audiolang := flag.String("audiolang", "ja", "The audio language you want to use. (IETF language tag)")
	targetfolder := flag.String("outputfolder", proposedConvertedDir, "The folder to put the converted videos in.")
	originalsfolder := flag.String("move-originals-to", "originals", "The alternative folder you want the originals moved to.")
	keepSubsAfter := flag.Bool("keepsubs", false, "Keep subs in the folder after conversion instead of deleting them.")
	cleanupSubs := flag.Bool("cleansubs", false, "Clean up the subtitles (in the case of srt) to make them render better.")
	extractFonts := flag.Bool("extract-fonts", true, "Extract the fonts from the mkv to use them in the hardcoding.")
	verboseOutput := flag.Bool("v", false, "Show verbose output.")
	crf := flag.Int("crf", 18, "Constant Rate Factor setting for ffmpeg.")
	firstOnly := flag.Bool("first-only", false, "Only convert the first file. (For testing purposes)")
	forOldDevices := flag.Bool("for-old-devices", false, "Use ffmpeg flags to get widest compatibility.")
	fastVersion := flag.Bool("fastversion", false, "Do a second and third pass, making a video at 1.5x the speed.")
	keepSlowVersion := flag.Bool("keep-slow", false, "In case you're making fast versions, keep the slow versions as well.")
	mp4Output := flag.Bool(
		"mp4",
		true,
		"Make MP4 files instead of MKV files. Always doing this because MP4 is better supported.",
	)
	postCmd := flag.String("postcmd", "", "The command to run on completion. Use %%o for the output filename.")
	postSubExtrationCmd := flag.String(
		"postsubextract",
		"",
		"The command to run after we've extracted the subs but before conversion. Use %%s for subs filename.",
	)
	h26xTune := flag.String(
		"h26x-tune",
		"animation",
		"The tuning to use for h26x encoding. (film/animation/fastdecode/zerolatency/none)",
	)
	h26xPreset := flag.String("h26x-preset", "fast", "The preset to use for h26x encoding. (fast/medium/slow/etc..)")
	h265 := flag.Bool("h265", false, "Use H265 encoding.")
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
	VERBOSE = *verboseOutput
	// TODO: this should not be in this function, it doesn't belong here.
	// this should be done before we get which files we need to convert though.
	if *detox {
		fmt.Print("Detoxing folder...")
		detoxWords := strings.Split(*detoxRemove, ",")
		if err := DetoxMkvsInFolder(*folder, detoxWords...); err != nil {
			log.Fatal("Cannot detox folder?!", err)
		}
		fmt.Print("done.\n")
	}

	files, err := os.ReadDir(*folder)
	if err != nil {
		log.Fatal(err)
	}

	config := Config{
		SubsLang:        *subslang,
		AudioLang:       *audiolang,
		SubsName:        *subsname,
		TargetFolder:    *targetfolder,
		OriginalsFolder: *originalsfolder,
		ExtractFonts:    *extractFonts,
		FirstOnly:       *firstOnly,
		Mp4:             *mp4Output,
		H26xTune:        *h26xTune,
		H26xPreset:      *h26xPreset,
		H265:            *h265,
		PostCmd:         *postCmd,
		PostSubExtract:  *postSubExtrationCmd,
		FilesToConvert:  files,
		KeepSubs:        *keepSubsAfter,
		CleanupSubs:     *cleanupSubs,
		FastVersion:     *fastVersion,
		KeepSlowVersion: *keepSlowVersion,
		Crf:             *crf,
		ForOldDevices:   *forOldDevices,
		Extension:       *extension,
	}

	if VERBOSE {
		litter.Dump(config)
	}
	return config
}
