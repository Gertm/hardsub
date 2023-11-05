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

func TestIntroFramesForFilename(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    IntroBoundaries
		wantErr bool
	}{
		// TODO: Add test cases.
		{"undeadunluck",
			args{
				"Undead_Unluck_-_05.mkv",
			},
			IntroBoundaries{
				Begin: "/home/gert/.config/hardsub/Undead_Unluck_begin.png",
				End:   "/home/gert/.config/hardsub/Undead_Unluck_end.png",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IntroFramesForFilename(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("IntroFramesForFilename() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntroFramesForFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}
