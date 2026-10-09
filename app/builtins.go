package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
	"strings"
)

type Command struct {
	exec   func(cmd Command, args []string)
	Stdout io.Writer
	Stderr io.Writer
}

func (c Command) Run(args []string) {
	c.exec(c, args)
}

var BuiltinNames = []string{"echo", "type", "pwd", "cd", "exit"}

var BuiltinCmds = map[string]Command{
	"echo": {exec: handleEcho},
	"type": {exec: handleType},
	"pwd":  {exec: handlePwd},
	"cd":   {exec: handleCd},
	"exit": {exec: handleExit},
}

func handleCd(cmd Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(cmd.Stderr, "cd needs at least 1 argument")
		return
	}
	path := args[0]
	if subPath, found := strings.CutPrefix(path, "~"); found {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(cmd.Stderr, err)
			return
		}
		path = homeDir + subPath
	}
	fileInfo, err := os.Stat(path)
	if err != nil || !fileInfo.IsDir() {
		fmt.Fprintf(cmd.Stderr, "cd: %s: No such file or directory\n", path)
		return
	}
	if err := os.Chdir(path); err != nil {
		fmt.Fprintf(cmd.Stderr, "cd: %s: %v\n", path, err)
	}
}

func handleEcho(cmd Command, args []string) {
	fmt.Fprintln(cmd.Stdout, strings.Join(args, " "))
}

func handleType(cmd Command, args []string) {
	if len(args) == 0 {
		return
	}
	arg := args[0]
	if slices.Contains(BuiltinNames, arg) {
		fmt.Fprintln(cmd.Stdout, arg+" is a shell builtin")
	} else if absPath, err := exec.LookPath(arg); err == nil {
		fmt.Fprintln(cmd.Stdout, arg+" is "+absPath)
	} else {
		fmt.Fprintln(cmd.Stdout, arg+": not found")
	}
}

func handlePwd(cmd Command, args []string) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(cmd.Stderr, "pwd error: ", err)
	} else {
		fmt.Fprintln(cmd.Stdout, dir)
	}
}

func handleExit(cmd Command, args []string) {
	os.Exit(0)
}
