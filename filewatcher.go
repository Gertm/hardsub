package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

func startProcess(path string, arguments string) error {
	if arguments == "" {
		cmd := exec.Command(path)
		cmd.Stdout = os.Stdout
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		return cmd.Wait()
	} else {
		cmd := exec.Command(path, strings.Split(arguments, " ")...)
		cmd.Stdout = os.Stdout
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		return cmd.Wait()
	}
}

func waitForUploadToFinish(file string) error {
	var size int64
	size = 0
	sameSizeCount := 0
	log.Printf("Waiting for write operations to stop on %v\n", file)
	defer func() {
		count--
		monitored_files.Delete(file)
	}()
	for {
		time.Sleep(1 * time.Second)
		fi, err := os.Stat(file)
		if err != nil {
			return err
		}
		currentSize := fi.Size()
		if currentSize > size {
			size = currentSize
			sameSizeCount = 0
			continue
		}
		if currentSize == fi.Size() {
			sameSizeCount += 1
		}
		if sameSizeCount == 3 {
			return nil
		}
	}
}

var (
	count           int
	monitored_files sync.Map
)

func watchForFiles(watchfolder string) {
	count = 0
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	var dir string

	if watchfolder == "" {
		dir, _ = os.Getwd()
	} else {
		dir = watchfolder
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					if _, ok := monitored_files.Load(event.Name); ok {
						// already watching this file
						log.Printf("Already watching %v\n", event.Name)
						return
					} else {
						monitored_files.Store(event.Name, true)
						count++
					}
					go func() {
						if strings.HasSuffix(strings.ToLower(event.Name), ".mkv") {
							log.Printf("A new file is being written: %v\n", event.Name)
							err := waitForUploadToFinish(event.Name)
							if err != nil {
								// Could not wait for the file correctly. Something must have gone awry.
								log.Println(err)
							} else {
								if count <= 0 {
									log.Println("No more files being written to.")
									// start the converting here.
								} else {
									log.Printf(".")
									time.Sleep(3 * time.Second)
								}
							}
						}
					}()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	log.Printf("Watching folder %v for incoming .mkv files.\n", dir)
	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
