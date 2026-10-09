package main

import (
	"fmt"
	"os"
)

func parseRedirect(args []string) (rest []string, stdoutFile *os.File, stderrFile *os.File, err error) {
	var stdout, stderr string
	var isAppendStdout, isAppendStderr bool

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case ">", "1>":
			if len(args) > i+1 {
				stdout = args[i+1]
				i++
			}
		case ">>", "1>>":
			if len(args) > i+1 {
				stdout = args[i+1]
				isAppendStdout = true
				i++
			}
		case "2>":
			if len(args) > i+1 {
				stderr = args[i+1]
				i++
			}
		case "2>>":
			if len(args) > i+1 {
				stderr = args[i+1]
				isAppendStderr = true
				i++
			}
		default:
			rest = append(rest, args[i])
		}
	}

	if stdout == "" && stderr == "" {
		return args, nil, nil, nil
	}

	if stdout != "" {
		flag := os.O_WRONLY | os.O_CREATE
		if isAppendStdout {
			flag |= os.O_APPEND
		} else {
			flag |= os.O_TRUNC
		}
		stdoutFile, err = os.OpenFile(stdout, flag, 0666)
		if err != nil {
			fmt.Println("Error creating output file: ", err)
			return nil, nil, nil, err
		}
	}

	if stderr != "" {
		flag := os.O_WRONLY | os.O_CREATE
		if isAppendStderr {
			flag |= os.O_APPEND
		} else {
			flag |= os.O_TRUNC
		}
		stderrFile, err = os.OpenFile(stderr, flag, 0666)
		if err != nil {
			fmt.Println("Error creating err file: ", err)
			return nil, nil, nil, err
		}
	}

	return rest, stdoutFile, stderrFile, nil
}
