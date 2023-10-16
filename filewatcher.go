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
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

var PollInterval = 1

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

// This is a bit of a naive way of checking if the file is done writing.
// Yet it works quite well in practise for me. Then again, I have quite
// reliable internet, so that helps. So this can certainly be improved.
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
		time.Sleep(time.Duration(PollInterval) * time.Second)
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

func watchForFiles(watchDirectory string, f func() error) {
	count = 0
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	var dir string

	if watchDirectory == "" {
		dir, _ = os.Getwd()
	} else {
		dir = watchDirectory
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
								count -= 1
								if count <= 0 {
									log.Println("No more files being written to.")
									if err := f(); err != nil {
										log.Println(err)
									}
								} else {
									log.Println("Some files are still being written to...", count)
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

	log.Printf("Watching directory %v for incoming .mkv files.\n", dir)
	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
