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
	"os"
	"path"
	"strings"

	"github.com/gertm/hardsub/subfix"
)

var LastSelectedTracks *SelectedTracks

func main() {
	InitConfig()
	LoadConfig()

	// Do we really need to start if these aren't available?
	for _, exe := range []string{"ffmpeg", "ffprobe", "mkvmerge"} {
		if _, err := FindInPath(exe); err != nil {
			fmt.Printf("Need to have %s on $PATH to work.\n", exe)
			return
		}
		if config.Verbose {
			fmt.Println("✅ Found", exe)
		}
	}

	if config.arguments.File != "" {
		if config.arguments.OnlyCut && config.arguments.CutStart != "" && config.arguments.CutEnd != "" {
			_, err := CutFragmentFromVideo(config)
			if err != nil {
				LogError("could not cut fragment from output: %s\n", err)
			}
		}
		if config.arguments.DumpFramesAt != "" {
			timestamps := strings.FieldsFunc(config.arguments.DumpFramesAt, func(c rune) bool { return c == ',' })
			// ffmpeg -ss 00:01:00 -i input.mp4 -frames:v 1 output.png
			fmt.Println(timestamps)
			for _, timestamp := range timestamps {
				if name, err := DumpFrameFromVideoAt(config.arguments.File, timestamp); err != nil {
					fmt.Fprintln(os.Stderr, name, err)
				}
			}
		}
		return
	}
	if config.WatchForFiles {
		watchForFiles(config.arguments.SourceFolder, func() error {
			PrepareFolderForConversion(&config)
			return ConvertAllTheThings(config)
		})
	} else {

		if config.Detox {
			Log("Detoxing folder...")
			detoxWords := strings.Split(config.RemoveWords, ",")
			if err := DetoxMkvsInFolder(config.arguments.SourceFolder, detoxWords...); err != nil {
				LogErrorln("Cannot detox folder?!", err)
			}
			Log("done.")
		}
		files, err := os.ReadDir(config.arguments.SourceFolder)
		if err != nil {
			LogErrorln("Could not read files from", config.arguments.SourceFolder)
			os.Exit(1)
		}
		config.filesToConvert = files
		if err := ConvertAllTheThings(config); err != nil {
			LogErrorln("Something went wrong while converting:", err)
			os.Exit(1)
		}
	}
	fmt.Println("Done!")
}

func ConvertAllTheThings(config Config) error {
	for _, file := range config.filesToConvert {
		if path.Ext(file.Name()) == "."+config.Extension {
			Log("Need to convert", file.Name())
			fullpath := file.Name()
			_, err := convert_file(fullpath, config)
			if err != nil {
				return err
			}
			if config.FirstOnly {
				Log("Done!")
				return nil
			}
		}
	}
	return nil
}

// TODO: This needs to be split up in smaller chunks, it's way too big now.
// Returns the converted filename and an error.
func convert_file(videofile string, config Config) (string, error) {
	Log("Converting", videofile)
	output, err := SelectTracksWithMkvMerge(videofile, config)
	if err != nil {
		LogErrorln("Could not select tracks with mkvmerge:", err)
		os.Exit(1)
	}
	LastSelectedTracks = output
	// write the script to convert.
	noext := strings.Replace(videofile, path.Ext(videofile), "", 1)
	var outputFile string
	if config.Mkv {
		outputFile = path.Join(config.TargetFolder, "HS_"+path.Base(videofile))
	} else {
		outputFile = path.Join(config.TargetFolder, strings.Replace(path.Base(videofile), ".mkv", ".mp4", 1))
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
		if err := RunAndParseFfmpeg(picSubsExtractCommand, vProps); err != nil {
			return "", fmt.Errorf("error while extracting picture subs: %w", err)
		}

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
		if err := RunAndParseFfmpeg(srtSubsExtractCommand, vProps); err != nil {
			return "", fmt.Errorf("error while extracting subs: %w", err)
		}
		if !config.KeepSubs {
			defer os.Remove(subsfile)
		}

		if config.PostSubExtract != "" {
			postsubcmd := strings.ReplaceAll(config.PostSubExtract, "%%s", subsfile) + "\n"
			if err := RunBashCommand(postsubcmd); err != nil {
				LogErrorln("Post Sub Extraction Command failed, check your script?\n", err)
			}
		}
		if output.SubtitleType == SRT {
			subfix.FixSubs(subsfile, 22, true, config.Verbose)
		}
		if config.ExtractFonts {
			if err := extractFonts(config.TargetFolder, videofile); err != nil {
				return "", fmt.Errorf("error extracting subs: %w", err)
			}
		} // TODO: Make this entire section template based.
		audioCodec := "copy"
		if !config.Mkv {
			audioCodec = "aac"
		}
		convertCmd := fmt.Sprintf("-y -hide_banner -loglevel error -stats -i %s -map 0:%d -map 0:%d -vf subtitles=%s -c:a %s -c:v %s -crf %d -preset %s %s%s%s",
			videofile, output.VideoTrack, output.AudioTrack, subsfile, audioCodec, videoCodec, config.Crf, config.H26xPreset, h26xTune, oldDevices, outputFile)
		Log("Convert Command:", "ffmpeg", convertCmd)
		fmt.Println("Starting re-encoding...")
		if err := RunAndParseFfmpeg(convertCmd, vProps); err != nil {
			return "", fmt.Errorf("error running the conversion for %s: %w", videofile, err)
		}
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
				outputFile = fastOutputFile
			}
		}
	}

	if !config.arguments.OnlyCut && config.arguments.CutStart != "" && config.arguments.CutEnd != "" {
		cutFile, err := CutFragmentFromVideo(config)
		if err != nil {
			fmt.Printf("could not cut fragment from output: %s\n", err)
		} else {
			outputFile = cutFile
		}
	}

	if config.PostCmd != "" {
		postcommand := strings.ReplaceAll(config.PostCmd, "%%o", outputFile)
		if err := RunBashCommand(postcommand); err != nil {
			fmt.Println("Post command failed, check your script?\n", err)
		}
	}

	if config.OriginalsFolder != config.TargetFolder {
		if err := createDirectoryIfNeeded(config.OriginalsFolder); err == nil {
			movedFile := path.Join(config.OriginalsFolder, baseVideoFile)
			os.Rename(videofile, movedFile)
		}
	}
	return outputFile, nil
}
