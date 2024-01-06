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
	"reflect"
	"testing"
)

func Test_logV(t *testing.T) {
	type args struct {
		format string
		v      []any
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logV(tt.args.format, tt.args.v...)
		})
	}
}

func TestLog(t *testing.T) {
	type args struct {
		msg []any
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Log(tt.args.msg...)
		})
	}
}

func TestLogError(t *testing.T) {
	type args struct {
		format string
		v      []any
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LogError(tt.args.format, tt.args.v...)
		})
	}
}

func TestLogErrorln(t *testing.T) {
	type args struct {
		msgs []any
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LogErrorln(tt.args.msgs...)
		})
	}
}

func TestDetoxFilename(t *testing.T) {
	type args struct {
		filename string
		remove   []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			"simple test",
			args{
				"/this/is/a[test]/file[123].mkv",
				[]string{},
			},
			"/this/is/a[test]/file123.mkv",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DetoxFilename(tt.args.filename, tt.args.remove...); got != tt.want {
				t.Errorf("DetoxFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetoxMkvsInDirectory(t *testing.T) {
	type args struct {
		dirname string
		remove  []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DetoxMkvsInDirectory(tt.args.dirname, tt.args.remove...); (err != nil) != tt.wantErr {
				t.Errorf("DetoxMkvsInDirectory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRunBashCommand(t *testing.T) {
	type args struct {
		cmd string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RunBashCommand(tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("RunBashCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_createDirectoryIfNeeded(t *testing.T) {
	type args struct {
		dirName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createDirectoryIfNeeded(tt.args.dirName); (err != nil) != tt.wantErr {
				t.Errorf("createDirectoryIfNeeded() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_copyFontsToLocalFontsDir(t *testing.T) {
	type args struct {
		sourcedir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := copyFontsToLocalFontsDir(tt.args.sourcedir); (err != nil) != tt.wantErr {
				t.Errorf("copyFontsToLocalFontsDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_copyFile(t *testing.T) {
	type args struct {
		src string
		dst string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			copyFile(tt.args.src, tt.args.dst)
		})
	}
}

func Test_refreshFonts(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			refreshFonts()
		})
	}
}

func Test_extractFonts(t *testing.T) {
	type args struct {
		workingdir string
		videofile  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := extractFonts(tt.args.workingdir, tt.args.videofile); (err != nil) != tt.wantErr {
				t.Errorf("extractFonts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetConfigFilenameForVideo(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetConfigFilenameForVideo(tt.args.path); got != tt.want {
				t.Errorf("GetConfigFilenameForVideo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelectTracksWithMkvMerge(t *testing.T) {
	type args struct {
		path   string
		config Config
	}
	tests := []struct {
		name    string
		args    args
		want    *SelectedTracks
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SelectTracksWithMkvMerge(tt.args.path, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectTracksWithMkvMerge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SelectTracksWithMkvMerge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindInPath(t *testing.T) {
	type args struct {
		exe string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindInPath(tt.args.exe)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindInPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FindInPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileExists(tt.args.filename); got != tt.want {
				t.Errorf("FileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
