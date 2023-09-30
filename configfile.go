package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/sanity-io/litter"
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
		log.Fatal(err)
	}
	configDir := filepath.Join(ucd, "hardsub")
	if _, err := os.Stat(configDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(configDir, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
	return filepath.Join(configDir, "hardsub.toml")
}

func LoadConfig() {
	configFile := configFilename()
	fp := file.Provider(configFile)
	ReadConfigFile(fp)
	litter.Dump(config)
	fp.Watch(func(event interface{}, err error) {
		if err != nil {
			log.Println("Error in watch callback", err)
		}
		ReadConfigFile(fp)
	})
}

func ReadConfigFile(fp *file.File) {
	koanfConfig = koanf.New(".")
	if err := koanfConfig.Load(fp, toml.Parser()); err != nil {
		log.Fatalln("Error loading config: ", err)
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
