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
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

type VideoProperties struct {
	Filename        string
	NrOfVideoFrames int
	Duration        string
}

func GetVideoPropertiesWithFFProbe(filename string) VideoProperties {
	ffprobe, err := FindInPath("ffprobe")
	if err != nil {
		fmt.Println("No ffprobe found:", err)
		return VideoProperties{}
	}
	ffprobeCommand := ffprobe + " -v error -select_streams v:0 -count_packets -show_entries stream=nb_read_packets -print_format csv " + filename
	// output, err := OutputForBashCommand(ffprobeCommand)
	// ffprobe -i video -show_entries format=duration -v quiet -sexagesimal -of csv
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
	durationCommand := fmt.Sprintf("%s -i %s -show_entries format=duration -v quiet -sexagesimal -of csv", ffprobe, filename)
	output, err = OutputForCommand(durationCommand)
	if err != nil {
		fmt.Println("Could not get duration with ffprobe")
	}
	duration := strings.TrimSpace(strings.Split(string(output), ",")[1])
	return VideoProperties{
		Filename:        filename,
		NrOfVideoFrames: packets,
		Duration:        duration,
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

func SearchForFrame(videoFile, frameImage string) (time.Duration, error) {
	// ffmpeg -loglevel info -i video.mkv -loop 1 -i frameImage.jpg -an -filter_complex "blend=difference:shortest=1,blackframe=98:32" -f null -
	args := fmt.Sprintf("-loglevel info -i %s -loop 1 -i %s -an -filter_complex blend=difference:shortest=1,blackframe=98:32 -f null -progress - -", videoFile, frameImage)

	cmd := exec.Command("ffmpeg", strings.Split(args, " ")...)

	stderr, _ := cmd.StderrPipe() // for some reason ffmpeg outputs to stderr only.
	cmd.Start()
	// [Parsed_blackframe_1 @ 0x9f27880] frame:2470 pblack:100 pts:103020 t:103.020000 type:I last_keyframe:2448
	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		if strings.HasPrefix(m, "t:") {
			if s, err := strconv.ParseFloat(m[2:], 64); err == nil {
				cmd.Process.Kill()
				return time.ParseDuration(fmt.Sprintf("%fs", s))
			}
		}
	}
	return 0, fmt.Errorf("cannot find frame")
}

func formatDuration(d time.Duration) string {
	t := time.Unix(0, 0).UTC()
	return t.Add(d).Format("15:04:05.000")
}

// returns the full name of the videoFile with the fragment cut out.
func cutFromVideo(ts_start, ts_end time.Duration, videoFile string) (string, error) {
	noIntroFile := strings.ReplaceAll(videoFile, path.Base(videoFile), "NOINTRO_"+path.Base(videoFile))
	firstPart := strings.ReplaceAll(videoFile, path.Base(videoFile), "first_"+path.Base(videoFile))
	lastPart := strings.ReplaceAll(videoFile, path.Base(videoFile), "last_"+path.Base(videoFile))
	videoProps := GetVideoPropertiesWithFFProbe(videoFile)
	start := formatDuration(ts_start)
	end := formatDuration(ts_end)
	// first make the pre-fragment video
	firstArgs := fmt.Sprintf("-y -i %s -ss 00:00:00 -to %s -c:v copy -c:a aac %s", videoFile, start, firstPart)
	lastArgs := fmt.Sprintf("-y -i %s -ss %s -to %s -c:v copy -c:a aac %s", videoFile, end, videoProps.Duration, lastPart)
	concatInput := fmt.Sprintf("file '%s'\nfile '%s'", firstPart, lastPart)
	os.WriteFile("concat.txt", []byte(concatInput), 0x644)
	// defer os.RemoveAll("concat.txt")
	concatArgs := fmt.Sprintf("-y -f concat -i concat.txt -c copy %s", noIntroFile)
	fmt.Println("ffmpeg", firstArgs, "\nffmpeg", lastArgs, "\nffmpeg", concatArgs)
	fmt.Println("Cutting first part...")
	if err := RunAndParseFfmpeg(firstArgs, videoProps); err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("Cutting second part...")
	if err := RunAndParseFfmpeg(lastArgs, videoProps); err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("Concatenating the two pieces...")
	if err := RunAndParseFfmpeg(concatArgs, videoProps); err != nil {
		fmt.Println(err)
		return "", err
	}
	return noIntroFile, nil
}

func CutFragmentFromVideo(startFrameFile, endFrameFile, videoFile string) (string, error) {
	fmt.Println("Looking for start of fragment...")
	start, err := SearchForFrame(videoFile, startFrameFile)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("Found start frame at", start)
	fmt.Println("Looking for end of fragment...")
	stop, err := SearchForFrame(videoFile, endFrameFile)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Printf("Cutting out fragment between %v and %v\n", start, stop)
	return cutFromVideo(start, stop, videoFile)
}

func DumpFrameFromVideoAt(videoFile, time string) (string, error) {
	timestr := strings.ReplaceAll(time, ":", "-")

	jpegname := strings.ReplaceAll(videoFile, path.Base(videoFile), "FRAME_"+timestr+"_"+strings.ReplaceAll(path.Base(videoFile), path.Ext(videoFile), ".png"))

	fmt.Println(timestr, jpegname)
	cmd := "-ss " + time + " -i " + videoFile + " -frames:v 1 " + jpegname
	fmt.Println("ffmpeg", cmd)
	err := RunAndParseFfmpeg(cmd, GetVideoPropertiesWithFFProbe(videoFile))
	return jpegname, err
}

func AddChaptersToVideo() {
	/*
		https://ikyle.me/blog/2020/add-mp4-chapters-ffmpeg

		;FFMETADATA1
		major_brand=isom
		minor_version=512
		compatible_brands=isomiso2avc1mp41
		encoder=Lavf60.3.100

		[CHAPTER]
		TIMEBASE=1/1000 // this needs to stay this way because of ffmpeg. Only milliseconds allowed.
		START=1
		END=448000
		title=Pre Intro

		[CHAPTER]
		TIMEBASE=1/1000
		START=448001
		END= 3883999
		title=Intro

		[CHAPTER]
		TIMEBASE=1/1000
		START=3884000
		END=4418000
		title=Post Intro


	*/
}
