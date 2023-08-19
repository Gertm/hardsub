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

func timestampForFrameInVideo(videoFile, frameFile string) (string, error) {
	ffmpegCommand := fmt.Sprintf("-i %s -loop 1 -i %s -an -filter_complex \"blend=difference:shortest=1,blackframe=99:32\" -f null -", videoFile, frameFile)
	return "", nil
}
