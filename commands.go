/*
	Copyright (C) 2023 Gert Meulyzer

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
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
