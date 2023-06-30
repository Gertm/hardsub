package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/schollz/progressbar/v3"
)

type VideoProperties struct {
	NrOfVideoFrames int
	Filename        string
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

func DoubleSpeedFile(input string, output string) error {
	inputProps := GetVideoPropertiesWithFFProbe(input)
	inputProps.NrOfVideoFrames = inputProps.NrOfVideoFrames / 2
	ffmpegArgs := fmt.Sprintf("-i %s -filter_complex [0:v]setpts=0.5*PTS[v];[0:a]atempo=2.0[a] -map [v] -map [a] %s", input, output)
	return RunAndParseFfmpeg(ffmpegArgs, inputProps)
}

/* 1.5x
ffmpeg -i $f -map 0:v -c:v copy -bsf:v h264_mp4toannexb raw.h264
ffmpeg -fflags +genpts -r 36 -i raw.h264 -i $f -map 0:v -c:v copy -map 1:a -af atempo=1.5 -movflags faststart FAST_$f
rm raw.h264
*/
