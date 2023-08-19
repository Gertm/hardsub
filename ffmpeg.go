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
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/schollz/progressbar/v3"
)

type VideoProperties struct {
	Filename        string
	NrOfVideoFrames int
}

func GetVideoPropertiesWithFFProbe(filename string) VideoProperties {
	ffprobe, err := FindInPath("ffprobe")
	if err != nil {
		fmt.Println("No ffprobe found:", err)
		return VideoProperties{}
	}
	ffprobeCommand := ffprobe + " -v error -select_streams v:0 -count_packets -show_entries stream=nb_read_packets -print_format csv " + filename
	// output, err := OutputForBashCommand(ffprobeCommand)
	output, err := OutputForCommand(ffprobeCommand)
	if err != nil {
		fmt.Println("Could not get VideoProperties with FFprobe:", err)
		return VideoProperties{}
	}
	nrOfPackets := strings.TrimSpace(strings.Split(string(output), ",")[1])
	packets, err := strconv.Atoi(nrOfPackets)
	if err != nil {
		fmt.Println("Number of packets not a number?!", "|"+nrOfPackets+"|")
		return VideoProperties{}
	}
	return VideoProperties{
		Filename:        filename,
		NrOfVideoFrames: packets,
	}
}

func GetVideoParamsFromFFMpeg(filename string) VideoProperties {
	fmt.Println("Getting video properties of", filename)
	cmd := exec.Command("ffmpeg", "-i", filename)
	stderr, _ := cmd.StderrPipe()
	cmd.Start()
	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanLines)
	var SawVideoStream bool

	for scanner.Scan() {
		m := scanner.Text()
		if strings.Contains(m, "Stream") {
			if strings.Contains(m, "Video") {
				SawVideoStream = true
			} else {
				SawVideoStream = false
			}
		}

		if SawVideoStream && strings.Contains(m, "NUMBER_OF_FRAMES") {
			fmt.Println("Saw NUMBER_OF_FRAMES")
			frameStr := strings.TrimSpace(strings.Split(m, ":")[1])
			frames, err := strconv.Atoi(frameStr)
			fmt.Println(frameStr, frames)
			if err != nil {
				return VideoProperties{}
			}
			return VideoProperties{Filename: filename, NrOfVideoFrames: frames}
		}
	}
	return VideoProperties{Filename: filename}
}

func RunAndParseFfmpeg(args string, prop VideoProperties) error {
	bar := progressbar.NewOptions(
		prop.NrOfVideoFrames,
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionSetDescription(prop.Filename),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionShowDescriptionAtLineEnd(),
		progressbar.OptionFullWidth(),
	)
	Log("ffmpeg", args)
	cmd := exec.Command("ffmpeg", strings.Split(args, " ")...)

	stderr, _ := cmd.StderrPipe() // for some reason ffmpeg outputs to stderr only.
	cmd.Start()

	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanWords)
	nextIsFrame := false
	for scanner.Scan() {
		m := scanner.Text()
		if nextIsFrame {
			nextIsFrame = false
			curFrame, err := strconv.Atoi(m)
			if err == nil {
				bar.Set(curFrame)
				// if we're not showing a progress bar yet, show progression of frames encoded.
				if curFrame < bar.GetMax()/100 {
					fmt.Printf("\r%d/%d : %s", curFrame, prop.NrOfVideoFrames, prop.Filename)
				}
				continue
			} else {
				fmt.Println(err)
			}

		}
		if strings.HasPrefix(m, "frame=") {
			if len(m) > 6 {
				// need to extract frames now, since no space separates the frames and
				// and the label 'frame='
				curFrame, err := strconv.Atoi(m[6:])
				if err == nil {
					bar.Set(curFrame)
					continue
				}
			} else {
				nextIsFrame = true
			}
		} else {
			nextIsFrame = false
		}
	}
	cmd.Wait()
	if !cmd.ProcessState.Success() {
		return fmt.Errorf("exitcode %d", cmd.ProcessState.ExitCode())
	}
	fmt.Printf("\n")
	return nil
}

func FastFile(inputFilePath string, outputFilePath string) error {
	inputProps := GetVideoPropertiesWithFFProbe(inputFilePath)
	firstPassArgs := fmt.Sprintf("-i %s -map 0:v -c:v copy -bsf:v h264_mp4toannexb raw.h264", inputFilePath)
	defer os.RemoveAll("raw.h264")
	err := RunAndParseFfmpeg(firstPassArgs, inputProps)
	if err != nil {
		return err
	}
	ffmpegArgs := fmt.Sprintf("-fflags +genpts -r 36 -i raw.h264 -i %s -map 0:v -c:v copy -map 1:a -af atempo=1.5 -movflags faststart %s", inputFilePath, outputFilePath)
	return RunAndParseFfmpeg(ffmpegArgs, inputProps)
}

/* 1.5x
ffmpeg -i $f -map 0:v -c:v copy -bsf:v h264_mp4toannexb raw.h264
ffmpeg -fflags +genpts -r 36 -i raw.h264 -i $f -map 0:v -c:v copy -map 1:a -af atempo=1.5 -movflags faststart FAST_$f
rm raw.h264
*/

func SearchForFrame(videoFile, frameImage string) (int, error) {
	// ffmpeg -loglevel info -i video.mkv -loop 1 -i frameImage.jpg -an -filter_complex "blend=difference:shortest=1,blackframe=98:32" -f null -
	ffmpegArgs := fmt.Sprintf("ffmpeg -loglevel info -i %s -loop 1 -i %s -an -filter_complex \"blend=difference:shortest=1,blackframe=98:32\" -f null -", videoFile, frameImage)
	fmt.Println(ffmpegArgs)
	return 0, nil
}
