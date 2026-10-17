package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// Split tokens into individual commands at each pipe.
func splitPipeline(inputs []string) ([][]string, error) {
	var commands [][]string
	var current []string

	for _, token := range inputs {
		if token == "|" {
			if len(current) == 0 {
				return nil, fmt.Errorf("syntax error near unexpected token '|'")
			}

			commands = append(commands, current)
			current = nil
		} else {
			current = append(current, token)
		}
	}

	if len(current) == 0 {
		return nil, fmt.Errorf("syntax error: missing command after pipe")
	}

	commands = append(commands, current)
	return commands, nil
}

func hasPipe(inputs []string) bool {
	for _, token := range inputs {
		if token == "|" {
			return true
		}
	}
	return false
}

func runPipeline(inputs []string) error {
	commands, err := splitPipeline(inputs)
	if err != nil {
		return err
	}

	// Each pipe connects one command's output to the next command's input.
	readers := make([]*io.PipeReader, len(commands)-1)
	writers := make([]*io.PipeWriter, len(commands)-1)

	for i := range readers {
		readers[i], writers[i] = io.Pipe()
	}

	errs := make(chan error, len(commands))

	for i, tokens := range commands {

		go func() {
			var stdin io.Reader = os.Stdin
			var stdout io.Writer = os.Stdout

			if i > 0 {
				stdin = readers[i-1]
				defer readers[i-1].Close()
			}

			if i < len(commands)-1 {
				stdout = writers[i]
				defer writers[i].Close()
			}

			action := tokens[0]
			args := tokens[1:]

			if builtin, ok := BuiltinCmds[action]; ok {
				cmd := builtin
				cmd.stdin = stdin
				cmd.Stdout = stdout
				cmd.Stderr = os.Stderr
				cmd.Run(args)

				errs <- nil
				return
			}

			path, err := exec.LookPath(action)
			if err != nil {
				fmt.Fprintln(os.Stderr, action+": command not found")
				errs <- nil
				return
			}

			cmd := exec.Command(path, args...)
			cmd.Stdin = stdin
			cmd.Stdout = stdout
			cmd.Stderr = os.Stderr

			err = cmd.Run()

			// var exitErr *exec.ExitError
			// if err != nil && !errors.As(err, &exitErr) {
			// 	fmt.Fprintln(os.Stderr, err)
			// }

			errs <- nil
		}()
	}

	for range commands {
		<-errs
	}

	return nil
}
