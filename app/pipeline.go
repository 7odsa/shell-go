package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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

	n := len(commands)
	if n == 1 {
		return nil
	}

	type stage struct {
		cmd     *exec.Cmd
		builtin bool
	}

	pipes := make([][2]*os.File, n-1)

	for i := range pipes {
		r, w, err := os.Pipe()
		if err != nil {
			for j := 0; j < i; j++ {
				pipes[j][0].Close()
				pipes[j][1].Close()
			}
			return err
		}
		pipes[i] = [2]*os.File{r, w}
	}

	closePipes := func() {
		for _, p := range pipes {
			p[0].Close()
			p[1].Close()
		}
	}

	stages := make([]stage, n)
	errs := make(chan error, n)

	// Configure and start external commands first.
	for i, tokens := range commands {
		if _, ok := BuiltinCmds[tokens[0]]; ok {
			stages[i].builtin = true
			continue
		}

		path, err := exec.LookPath(tokens[0])
		if err != nil {
			closePipes()
			return fmt.Errorf("%s: command not found", tokens[0])
		}

		cmd := exec.Command(path, tokens[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if i > 0 {
			cmd.Stdin = pipes[i-1][0]
		}
		if i < n-1 {
			cmd.Stdout = pipes[i][1]
		}

		if err := cmd.Start(); err != nil {
			closePipes()
			return err
		}

		stages[i].cmd = cmd
	}

	// Run built-ins concurrently, using their own pipe endpoints.
	for i, tokens := range commands {
		if !stages[i].builtin {
			continue
		}

		i, tokens := i, tokens

		go func() {
			cmd := BuiltinCmds[tokens[0]]
			cmd.stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if i > 0 {
				cmd.stdin = pipes[i-1][0]
			}
			if i < n-1 {
				cmd.Stdout = pipes[i][1]
			}

			cmd.Run(tokens[1:])

			// Close the endpoint owned by this built-in after it returns.
			if i > 0 {
				pipes[i-1][0].Close()
			}
			if i < n-1 {
				pipes[i][1].Close()
			}

			errs <- nil
		}()
	}

	// The parent must close its copies so readers can receive EOF.
	// Child processes and built-in goroutines retain their assigned endpoints.
	closePipes()

	// Wait for all external processes and built-ins.
	for i := range stages {
		if stages[i].builtin {
			continue
		}
		errs <- stages[i].cmd.Wait()
	}

	builtinCount := 0
	for _, s := range stages {
		if s.builtin {
			builtinCount++
		}
	}
	for i := 0; i < builtinCount; i++ {
		<-errs
	}

	return nil
}

// Keep this helper available if you want to check whether a pipeline
// contains more than one command.
func pipelineCommandCount(inputs []string) int {
	return len(strings.Fields(strings.Join(inputs, " ")))
}
