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

func FastFile(input string, output string) error {
	inputProps := GetVideoPropertiesWithFFProbe(input)
	firstPassArgs := fmt.Sprintf("-i %s -map 0:v -c:v copy -bsf:v h264_mp4toannexb raw.h264", input)
	defer os.RemoveAll("raw.h264")
	err := RunAndParseFfmpeg(firstPassArgs, inputProps)
	if err != nil {
		return err
	}
	ffmpegArgs := fmt.Sprintf("-fflags +genpts -r 36 -i raw.h264 -i %s -map 0:v -c:v copy -map 1:a -af atempo=1.5 -movflags faststart %s", input, output)
	return RunAndParseFfmpeg(ffmpegArgs, inputProps)
}

/* 1.5x
ffmpeg -i $f -map 0:v -c:v copy -bsf:v h264_mp4toannexb raw.h264
ffmpeg -fflags +genpts -r 36 -i raw.h264 -i $f -map 0:v -c:v copy -map 1:a -af atempo=1.5 -movflags faststart FAST_$f
rm raw.h264
*/
