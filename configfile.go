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
	"log"
	"os"
	"path/filepath"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	flag "github.com/spf13/pflag"
)

var (
	koanfConfig = koanf.New(".")
	config      = Config{}
)

func InitConfig() {
	// TODO: Add a second possible config file .hardsub.toml in the current folder that could also
	// be read after the general config.
	// Then read the command line arguments as well and override the config again.
	if !FileExists(configFilename()) {
		SaveDefaultConfig()
	}
	LoadConfig()
}

func configFilename() string {
	ucd, err := os.UserConfigDir()
	if err != nil {
		LogError("could not get user config dir. What kind of system are you running?", err)
		os.Exit(1)
	}
	configDir := filepath.Join(ucd, "hardsub")
	if _, err := os.Stat(configDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(configDir, os.ModePerm)
		if err != nil {
			LogError("can't make the config folder:", configDir, err)
			os.Exit(1)
		}
	}
	return filepath.Join(configDir, "hardsub.toml")
}

func LoadConfig() {
	configFile := configFilename()
	fp := file.Provider(configFile)
	ReadConfigFile(fp)
	f := flag.NewFlagSet("config", flag.ExitOnError)
	f.Usage = func() {
		fmt.Print("hardsub - A tool to burn subtitles into videos.\n\n")
		fmt.Println("Most of the configuration is in " + configFile)
		fmt.Print("\nExtra command line options:\n\n")
		fmt.Println(f.FlagUsages())
	}
	f.String("file", "", "The specific file to operate on for cutting and frame dumping.")
	f.Bool("onlycut", false, "Only cut, don't convert.")
	f.String("cutstart", "", "A jpg of the frame to look for in the video to determine the start of the fragment to cut out.")
	f.String("cutend", "", "A jpg of the frame to look for in the video to determine the end of the fragment to cut out.")
	f.String("dumpframesat", "", "A comma-separated list of timestamps you want to make jpg dumps for.")
	f.Int("force-audio-track", -1, "Force the audio track to use. (for example: 4)")
	f.Int("force-subs-track", -1, "Force the subs track to use. (for example: 3)")
	wd, _ := os.Getwd()
	f.String("sourcefolder", wd, "The folder in which to look for videos.")
	f.Parse(os.Args[1:])

	ka := koanf.New(".")
	if err := ka.Load(posflag.Provider(f, ".", ka), nil); err != nil {
		LogError("cannot load the command line arguments:", err)
		os.Exit(1)
	}
	config.arguments = Arguments{}
	config.arguments.File = ka.String("file")
	config.arguments.OnlyCut = ka.Bool("onlycut")
	config.arguments.CutStart = ka.String("cutstart")
	config.arguments.CutEnd = ka.String("cutend")
	config.arguments.DumpFramesAt = ka.String("dumpframesat")
	config.arguments.ForceAudioTrack = ka.Int("force-audio-track")
	config.arguments.ForceSubsTrack = ka.Int("force-subs-track")
	config.arguments.SourceFolder = ka.String("sourcefolder")
	if config.arguments.SourceFolder == "" {
		config.arguments.SourceFolder = wd
	}
}

func ReadConfigFile(fp *file.File) {
	koanfConfig = koanf.New(".")
	if err := koanfConfig.Load(fp, toml.Parser()); err != nil {
		LogError("cannot load config file:", err)
		os.Exit(1)
	}
	koanfConfig.Unmarshal("", &config)
}

func SaveDefaultConfig() {
	defConfig := DefaultConfig()
	d := koanf.New(".")
	if err := d.Load(structs.Provider(defConfig, "koanf"), nil); err != nil {
		panic(err)
	}
	b, err := d.Marshal(toml.Parser())
	if err != nil {
		panic(err)
	}
	os.WriteFile(configFilename(), b, 0o666)
}
