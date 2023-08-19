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
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func addEnvironment(c *exec.Cmd, envStr string) {
	c.Env = append(os.Environ(), envStr)
}

func setDefaultEnvironment(c *exec.Cmd) {
	addEnvironment(c, "DISPLAY=:0")
}

func ExecuteCommand(cmd string) error {
	parts, err := SplitCommand(cmd)
	if err != nil {
		return err
	}
	c := exec.Command(parts[0], parts[1:]...)
	setDefaultEnvironment(c)
	return c.Run()
}

// split maar hou rekening met double quotes
func SplitCommand(cmd string) ([]string, error) {
	r := csv.NewReader(strings.NewReader(cmd))
	r.Comma = ' '
	fields, err := r.Read()
	if err != nil {
		return []string{}, err
	}
	return fields, nil
}

func OutputForCommand(cmd string) (string, error) {
	parts, err := SplitCommand(cmd)
	if err != nil {
		return "", err
	}
	raw, err := exec.Command(parts[0], parts[1:]...).Output()
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func OutputForCommandLst(cmd []string) (string, error) {
	c := exec.Command(cmd[0], cmd[1:]...)
	setDefaultEnvironment(c)
	raw, err := c.Output()
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func OutputForCommandLst2(cmd ...string) (string, error) {
	return OutputForCommandLst(cmd)
}

func OutputForBashCommand(cmd string) (string, error) {
	if VERBOSE {
		fmt.Println("bash -c", cmd)
	}
	c := exec.Command("bash", "-c", cmd)
	setDefaultEnvironment(c)
	raw, err := c.Output()
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func WriteBashScriptForCommand(filename, cmd string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString("#!/bin/bash\n")
	f.WriteString("set -e\n")
	f.WriteString(cmd + "\n")
	return nil
}
