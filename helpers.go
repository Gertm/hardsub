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
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/sanity-io/litter"
)

func logV(format string, v ...any) {
	if config.Verbose {
		fmt.Printf(format, v...)
	}
}

func Log(msg ...any) {
	if config.Verbose {
		fmt.Println(msg...)
	}
}

func LogError(format string, v ...any) {
	fmt.Fprintf(os.Stderr, format, v...)
}

func LogErrorln(msgs ...any) {
	fmt.Fprintln(os.Stderr, msgs...)
}

func DetoxFilename(filename string, remove ...string) string {
	baseName := path.Base(filename)
	for _, r := range remove {
		baseName = strings.ReplaceAll(baseName, r, "")
	}
	var sb strings.Builder
	justWroteUnderscore := true // don't write underscores at the start
	for _, ch := range baseName {
		switch ch {
		case 32, 95:
			if !justWroteUnderscore {
				sb.WriteString("_")
				justWroteUnderscore = true
			}
			continue
		}
		if (ch > 47 && ch < 58) ||
			(ch > 64 && ch < 91) || ch == 45 ||
			(ch > 96 && ch < 123) || ch == 46 {
			sb.WriteRune(ch)
			justWroteUnderscore = false
		}

	}
	return path.Join(path.Dir(filename), sb.String())
}

func DetoxMkvsInFolder(foldername string, remove ...string) error {
	toxic, err := os.ReadDir(foldername)
	if err != nil {
		LogErrorln("cannot read the filenames in folder", foldername, err)
		os.Exit(1)
	}
	for _, f := range toxic {
		fullname := path.Join(foldername, f.Name())
		if strings.ToLower(path.Ext(fullname)) == ".mkv" {
			dt := DetoxFilename(fullname, remove...)
			if f.Name() != dt {
				err := os.Rename(fullname, dt)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func RunBashCommand(cmd string) error {
	c := exec.Command("bash", "-c", cmd)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		return err
	}
	return nil
}

func createDirectoryIfNeeded(dirName string) error {
	src, err := os.Stat(dirName)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dirName, 0o755)
		if errDir != nil {
			return err
		}
		return nil
	}

	if src.Mode().IsRegular() {
		return fmt.Errorf("%s is a file", dirName)
	}

	if src.Mode().IsDir() {
		return nil
	}

	return errors.New("shouldn't really get here")
}

func copyFontsToLocalFontsDir(sourcedir string) error {
	files, err := os.ReadDir(sourcedir)
	if err != nil {
		LogErrorln("Cannot read the folder we just created?!", err)
	}
	userHome, err := os.UserHomeDir()
	if err != nil {
		LogErrorln("Cannot get our home directory!", err)
	}
	dotFonts := path.Join(userHome, ".fonts")
	for _, file := range files {
		if strings.Index(strings.ToLower(file.Name()), ".ttf") > 0 ||
			strings.Index(strings.ToLower(file.Name()), ".otf") > 0 {
			logV("Copying %s to %s", file.Name(), dotFonts)
			copyFile(path.Join(sourcedir, file.Name()), path.Join(dotFonts, file.Name()))
		}
	}
	return nil
}

func copyFile(src, dst string) {
	fin, err := os.Open(src)
	if err != nil {
		LogErrorln("cannot open source file for copying", src, err)
	}
	defer fin.Close()

	fout, err := os.Create(dst)
	if err != nil {
		LogErrorln("cannot create target file for copy:", dst, err)
		log.Fatal(err)
	}
	defer fout.Close()

	_, err = io.Copy(fout, fin)

	if err != nil {
		LogError("cannot copy %s to %s", src, dst)
		log.Fatal(err)
	}
}

func refreshFonts() {
	exec.Command("fc-cache", "-f", "-v").Wait()
}

func extractFonts(workingdir, videofile string) error {
	attachmentsFolder := path.Join(workingdir, "attachments")
	if err := createDirectoryIfNeeded(attachmentsFolder); err != nil {
		return err
	}
	err := os.MkdirAll(attachmentsFolder, os.ModePerm)
	if err != nil {
		log.Printf("Cannot create %s, skipping font extraction.\n%s\n", attachmentsFolder, err)
		// the video conversion will work without the custom fonts, so we don't need to fail on this.
		return nil
	}
	// ffmpeg -dump_attachment:t "" -i input.mkv
	currentDir, _ := os.Getwd() // let's assume we can know where we are.
	os.Chdir(attachmentsFolder)
	defer os.Chdir(currentDir)
	exec.Command("ffmpeg", "-dump_attachment:t", "", "-i", videofile).Output()
	// copy all fonts to the ~/.fonts folder
	if err := copyFontsToLocalFontsDir(attachmentsFolder); err != nil {
		return err
	}
	refreshFonts()
	return nil
}

func GetConfigFilenameForVideo(path string) string {
	return strings.Replace(path, filepath.Ext(path), ".hcConfig", 1)
}

// , config Config, output *SelectedTracks
func SelectTracksWithMkvMerge(path string, config Config) (*SelectedTracks, error) {
	Log("Getting tracks with mkvmerge...", path)
	output := SelectedTracks{
		AudioTrack: -1,
		VideoTrack: -1,
		SubsTrack:  -1,
	}
	raw, err := exec.Command("mkvmerge", "-J", path).Output()
	if err != nil {
		fmt.Println("Running mkvmerge failed")
		return &SelectedTracks{}, err
	}
	var audioTracks []int
	var subsTracks []int
	jsonparser.ArrayEach(raw, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		trackType, _ := jsonparser.GetString(value, "type")
		codec, _ := jsonparser.GetString(value, "codec")
		codec_id, _ := jsonparser.GetString(value, "properties", "codec_id")
		id, _ := jsonparser.GetInt(value, "id")
		Log(id, trackType, codec, codec_id)
		switch trackType {
		case "video":
			if output.VideoTrack == -1 {
				output.VideoTrack = int(id)
			}
		case "audio":
			lang, _ := jsonparser.GetString(value, "properties", "language_ietf")
			audioTracks = append(audioTracks, int(id))
			if lang == config.AudioLang && output.AudioTrack == -1 {
				output.AudioTrack = int(id)
			}
		case "subtitles":
			lang, _ := jsonparser.GetString(value, "properties", "language_ietf")
			Log("Config subslang", config.SubsLang)
			Log("Language", lang)
			subsTracks = append(subsTracks, int(id))

			if strings.HasPrefix(lang, config.SubsLang) && output.SubsTrack == -1 {
				trackName, err := jsonparser.GetString(value, "properties", "track_name")
				Log("Trackname", trackName)
				if err == nil {
					if !strings.Contains(strings.ToLower(trackName), "songs") {
						output.SubsTrack = int(id)
					} else {
						fmt.Println(err)
					}
				} else {
					Log("no name for the track found, so let's assume it's the right one for now")
					output.SubsTrack = int(id)
				}
				switch codec_id {
				case "S_TEXT/ASS", "S_TEXT/SSA", "SAA/ASS":
					if config.Verbose {
						fmt.Printf("%s has SSA subs\n", path)
					}
					output.SubtitleType = SSA_ASS
				case "S_TEXT/UTF8":
					if config.Verbose {
						fmt.Printf("%s has SRT subtitles\n", path)
					}
					output.SubtitleType = SRT
				case "S_HDMV/PGS", "S_IMAGE/BMP", "S_DVDSUB", "S_VOBSUB":
					if config.Verbose {
						fmt.Printf("%s has picture based subtitles\n", path)
					}
					output.SubtitleType = PICTURE
				}
			}
		}
	}, "tracks")
	if output.AudioTrack == -1 && len(audioTracks) == 1 {
		output.AudioTrack = audioTracks[0]
	}
	if output.SubsTrack == -1 && len(subsTracks) == 1 {
		output.SubsTrack = subsTracks[0]
	}
	if config.arguments.ForceAudioTrack != -1 {
		output.AudioTrack = config.arguments.ForceAudioTrack
	}
	if config.arguments.ForceSubsTrack != -1 {
		output.SubsTrack = config.arguments.ForceSubsTrack
	}
	if config.Verbose {
		litter.Dump(output)
	}
	return &output, nil
}

func FindInPath(exe string) (string, error) {
	path, err := exec.LookPath(exe)
	if err != nil {
		fmt.Printf("Could not find %s in $PATH\n", exe)
		return "", err
	}
	return path, nil
}

func FileExists(filename string) bool {
	_, error := os.Stat(filename)
	return !errors.Is(error, os.ErrNotExist)
}
