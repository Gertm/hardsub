package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/gertm/hardsub/subfix"
)

var (
	LastSelectedTracks *SelectedTracks
	VERBOSE            bool
)

func main() {
	config := getConfigurationFromArguments()

	for _, file := range config.FilesToConvert {
		if path.Ext(file.Name()) == "."+config.Extension {
			Log("Need to convert", file.Name())
			fullpath := file.Name()
			convert_file(fullpath, config)
			if config.FirstOnly {
				fmt.Println("Done!")
				return
			}
		}
	}
	fmt.Println("Done!")
}

// TODO: This needs to be split up in smaller chunks, it's way too big now.
func convert_file(videofile string, config Config) {
	Log("Converting", videofile)
	// if config.FastVersion {
	// 	createDirectoryIfNeeded(config.TargetFolder)
	// 	outputFile := path.Join(config.TargetFolder, "FAST_"+path.Base(videofile))
	// 	if err := FastFile(videofile, outputFile); err != nil {
	// 		fmt.Println("Couldn't double the speed of", videofile, err)
	// 	}
	// 	return
	// }
	output, err := SelectTracksWithMkvMerge(videofile, config)
	if err != nil {
		log.Fatal(err)
	}
	LastSelectedTracks = output

	// write the script to convert.
	noext := strings.Replace(videofile, path.Ext(videofile), "", 1)
	var outputFile string
	if config.Mp4 {
		outputFile = path.Join(config.TargetFolder, strings.Replace(path.Base(videofile), ".mkv", ".mp4", 1))
	} else {
		outputFile = path.Join(config.TargetFolder, "HS_"+path.Base(videofile))
	}
	baseVideoFile := path.Base(videofile)
	var subsfile string
	videoCodec := "libx264"
	if config.H265 {
		videoCodec = "libx265"
	}
	h26xTune := ""
	if config.H26xTune == "none" {
		h26xTune = ""
	} else {
		h26xTune = "-tune " + config.H26xTune + " "
	}
	oldDevices := ""
	if config.ForOldDevices {
		oldDevices = " -profile:v baseline -level 3.0 -pix_fmt yuv420p -ac 2 -b:a 128k -movflags faststart "
	}
	vProps := GetVideoPropertiesWithFFProbe(videofile)
	if output.SubtitleType == PICTURE {
		picSubsExtractCommand := fmt.Sprintf(
			"-hide_banner -loglevel error -stats -y -i %s -filter_complex [0:v][0:s:0]overlay[v] -map [v] -map 0:%d -map 0:%d -c:v %s %s %s %s -c:a copy %s",
			videofile,
			output.VideoTrack,
			output.AudioTrack,
			videoCodec,
			fmt.Sprintf("-crf %d", config.Crf),
			fmt.Sprintf("-preset %s", config.H26xPreset),
			h26xTune,
			outputFile,
		)
		Log(picSubsExtractCommand)
		RunAndParseFfmpeg(picSubsExtractCommand, vProps)

	} else {
		// Extracting the subtitle file in case of text based ones, so we can forcibly select the correct one.
		if output.SubtitleType == SSA_ASS {
			subsfile = noext + ".ass"
		}
		if output.SubtitleType == SRT {
			subsfile = noext + ".srt"
		}
		srtSubsExtractCommand := fmt.Sprintf("-y -hide_banner -loglevel error -stats -txt_format text -i %s -map 0:%d %s", videofile, output.SubsTrack, subsfile)
		Log(srtSubsExtractCommand)
		RunAndParseFfmpeg(srtSubsExtractCommand, vProps)
		if !config.KeepSubs {
			defer os.Remove(subsfile)
		}

		if config.PostSubExtract != "" {
			postsubcmd := strings.ReplaceAll(config.PostSubExtract, "%%s", subsfile) + "\n"
			if err := RunBashCommand(postsubcmd); err != nil {
				fmt.Println("Post Sub Extraction Command failed, check your script?\n", err)
			}
		}
		if output.SubtitleType == SRT {
			subfix.FixSubs(subsfile, 22, true, VERBOSE)
		}
		if config.ExtractFonts {
			// writeExtractFontsCommand(config.TargetFolder, videofile, config.ScriptFile)
			extractFonts(config.TargetFolder, videofile)
		} // TODO: Make this entire section template based.
		audioCodec := "copy"
		if config.Mp4 {
			audioCodec = "aac"
		}
		convertCmd := fmt.Sprintf("-y -hide_banner -loglevel error -stats -i %s -map 0:%d -map 0:%d -vf subtitles=%s -c:a %s -c:v %s -crf %d -preset %s %s%s%s",
			videofile, output.VideoTrack, output.AudioTrack, subsfile, audioCodec, videoCodec, config.Crf, config.H26xPreset, h26xTune, oldDevices, outputFile)
		Log("Convert Command:", "ffmpeg", convertCmd)
		fmt.Println("Starting re-encoding...")
		RunAndParseFfmpeg(convertCmd, vProps)
	}
	if config.FastVersion {
		fastOutputFile := strings.ReplaceAll(outputFile, path.Base(outputFile), "FAST_"+path.Base(outputFile))
		fmt.Println(">>>>>>>>> Creating", fastOutputFile, ">>>>>>>>>>>")
		if err := FastFile(outputFile, fastOutputFile); err != nil {
			fmt.Println(err)
			if !config.KeepSlowVersion {
				fmt.Println("Keeping normal speed version because creating the fast version failed.")
			}
		} else {
			if !config.KeepSlowVersion {
				os.RemoveAll(outputFile)
			}
		}

	}
	if config.OriginalsFolder != config.TargetFolder {
		if err := createDirectoryIfNeeded(config.OriginalsFolder); err == nil {
			os.Rename(videofile, path.Join(config.OriginalsFolder, baseVideoFile))
		}
	}
	if config.PostCmd != "" {
		postcommand := strings.ReplaceAll(config.PostCmd, "%%o", outputFile)
		if err := RunBashCommand(postcommand); err != nil {
			fmt.Println("Post command failed, check your script?\n", err)
		}
	}
}
