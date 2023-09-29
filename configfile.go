package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/sanity-io/litter"
)

var (
	globalConfig = koanf.New(".")
	config       = Config{}
)

func ReadConfigFile() {
	ucd, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	configDir := filepath.Join(ucd, "hardsub")
	if _, err := os.Stat(configDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(configDir, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
	configFile := filepath.Join(configDir, "hardsub.toml")
	if err := globalConfig.Load(file.Provider(configFile), toml.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}
	globalConfig.Unmarshal("", &config)
	litter.Dump(config)
}
