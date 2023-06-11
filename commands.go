package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-multierror"
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

func ExecuteMultipleCommands(cmds []string) error {
	var result *multierror.Error
	for _, cmd := range cmds {
		fmt.Println("Running ", cmd)
		result = multierror.Append(result, ExecuteCommand(cmd))
	}
	return result.ErrorOrNil()
}
