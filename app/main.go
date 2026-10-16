// package main

// import (
// 	"fmt"
// 	"os"
// 	"strings"

// 	"github.com/chzyer/readline"
// 	"github.com/codecrafters-io/shell-starter-go/app/builtin"
// 	"github.com/codecrafters-io/shell-starter-go/app/utils"
// )

// var builtIns = []string{"exit", "echo", "type", "pwd", "cd", "cat"}

// type BellCompleter struct {
// 	readline.AutoCompleter
// }

// func (bc *BellCompleter) Do(line []rune, pos int) ([][]rune, int) {

// 	newLine, length := bc.AutoCompleter.Do(line, pos)
// 	if len(newLine) == 0 {
// 		fmt.Print("\x07")
// 		os.Stdout.Sync()
// 	}
// 	return newLine, length
// }

// func main() {

// 	allExe := utils.GetAllExecutableFiles()

// 	allPrefixComp := []readline.PrefixCompleterInterface{readline.PcItem("echo"),
// 		readline.PcItem("exit"),
// 	}

// 	for _, exe := range allExe {
// 		allPrefixComp = append(allPrefixComp, readline.PcItem(exe))
// 	}
// 	autoCompleter := readline.NewPrefixCompleter(allPrefixComp...)

// 	bellComp := &BellCompleter{AutoCompleter: autoCompleter}

// 	reader, err := readline.NewEx(&readline.Config{
// 		Prompt:       "$ ",
// 		AutoComplete: bellComp,
// 	})

// 	if err != nil {
// 		panic(err)
// 	}

// 	defer reader.Close()

// 	for {

// 		fmt.Print("$ ")

// 		prompt, err := reader.Readline()

// 		utils.HandleErr(err)

// 		prompt = strings.TrimSpace(prompt)
// 		if utils.IsEmptyPrompt(prompt) {
// 			continue
// 		}

// 		var tokens []string = utils.HandleQuotes(prompt)

// 		hasCommand(tokens, prompt)
// 	}

// }

// func hasCommand(tokens []string, prompt string) bool {
// 	// if len(tokens) > 1 {
// 	// 	tokens[1] = strings.TrimLeft(tokens[1], " ")

// 	// }

// 	switch {

// 	case utils.IsBuiltIn(tokens[0], "exit"):
// 		os.Exit(0)
// 		return true
// 	case utils.IsBuiltIn(tokens[0], "pwd"):
// 		builtin.HandlePwd(tokens)
// 		return true

// 	case utils.IsBuiltIn(tokens[0], "echo"):
// 		builtin.HandleEcho(tokens, strings.TrimLeft(strings.TrimSuffix(prompt, "echo"), " "))
// 		return true
// 	// case utils.IsBuiltIn(tokens[0], "cat"):
// 	// 	builtin.HandleCat(tokens, strings.TrimLeft(strings.TrimSuffix(prompt, "cat"), " "))
// 	// 	return true
// 	case utils.IsBuiltIn(tokens[0], "type"):
// 		builtin.HandleType(tokens, builtIns)
// 		return true

// 	case utils.IsBuiltIn(tokens[0], "cd"):
// 		builtin.HandleCd(tokens)
// 		return true

// 	case utils.IsExecutable(tokens, prompt):
// 		builtin.HandleExecutable(tokens, prompt)
// 		return true
// 	default:
// 		fmt.Println(strings.Join(tokens, "") + ": command not found")
// 		return false
// 	}

// }

package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/chzyer/readline"
)

const prompt = "$ "

func main() {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       prompt,
		AutoComplete: NewCustomCompleter(prompt),
		HistoryFile:  "/tmp/readline.tmp",
		HistoryLimit: 20,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()
	rl.CaptureExitSignal()

	for {
		isBackground := false

		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}
		inputs := tokenize(strings.TrimSpace(line))

		if len(inputs) == 0 {
			continue
		}

		if slices.Contains(inputs, "|") {
			runPipeline(inputs)
			continue
		}

		action := inputs[0]
		if inputs[len(inputs)-1] == "&" && action != "echo" {
			isBackground = true
			inputs = inputs[:len(inputs)-1]
		}
		args := inputs[1:]

		args, outFile, errFile, err := parseRedirect(args)
		if err != nil {
			fmt.Println("Error parsing redirect: ", err)
			continue
		}
		stdout := os.Stdout
		stderr := os.Stderr
		if outFile != nil {
			stdout = outFile
		}
		if errFile != nil {
			stderr = errFile
		}

		if cmd, isBuiltin := BuiltinCmds[action]; isBuiltin {
			cmd.Stdout = stdout
			cmd.Stderr = stderr
			cmd.Run(args)

		} else if _, err := exec.LookPath(action); err == nil {
			cmd := exec.Command(action, args...)
			cmd.Stdout = stdout
			cmd.Stderr = stderr

			if !isBackground {

				if err := cmd.Run(); err != nil {
					var exitErr *exec.ExitError
					// Only print non-exit message
					if isExit := errors.As(err, &exitErr); !isExit {
						fmt.Fprintln(stderr, err)
					}
				}

			} else {
				if err := cmd.Start(); err != nil {
					fmt.Fprintln(stderr, err)
				} else {
					job := &Job{number: len(jobs) + 1, cmd: cmd, done: make(chan struct{}), inputLine: line}
					jobs = append(jobs, job)

					go func() {
						cmd.Wait()
						close(job.done)
					}()
					fmt.Printf("[%d] %d\n", job.number, cmd.Process.Pid)
				}
			}

		} else {
			fmt.Fprintln(os.Stderr, action+": command not found")
		}

		printJobs(true)

		// Clean up opened file
		if outFile != nil {
			outFile.Close()
		}
		if errFile != nil {
			errFile.Close()
		}
	}
}

func runPipeline(inputs []string) bool {
	// Find the pipe operator.
	pipeIndex := -1

	for i, token := range inputs {
		if token == "|" {
			if pipeIndex != -1 {
				fmt.Fprintln(os.Stderr, "multiple pipes are not supported yet")
				return true
			}
			pipeIndex = i
		}
	}

	// No pipeline: let the normal command handler run.
	if pipeIndex == -1 {
		return false
	}

	// Validate both sides.
	if pipeIndex == 0 || pipeIndex == len(inputs)-1 {
		fmt.Fprintln(os.Stderr, "invalid pipeline")
		return true
	}

	left := inputs[:pipeIndex]
	right := inputs[pipeIndex+1:]

	// Parse redirections independently for each command.
	leftArgs, leftOut, leftErr, err := parseRedirect(left[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing redirect:", err)
		return true
	}

	rightArgs, rightOut, rightErr, err := parseRedirect(right[1:])
	if err != nil {
		closeFile(leftOut)
		closeFile(leftErr)
		fmt.Fprintln(os.Stderr, "Error parsing redirect:", err)
		return true
	}

	// Create the connection between the processes.
	r, w, err := os.Pipe()
	if err != nil {
		closeFile(leftOut)
		closeFile(leftErr)
		closeFile(rightOut)
		closeFile(rightErr)
		fmt.Fprintln(os.Stderr, err)
		return true
	}

	defer r.Close()
	defer w.Close()

	// Set up the left command.
	leftCmd := exec.Command(left[0], leftArgs...)
	leftCmd.Stdout = w
	leftCmd.Stderr = os.Stderr

	if leftOut != nil {
		leftCmd.Stdout = leftOut
	}
	if leftErr != nil {
		leftCmd.Stderr = leftErr
	}

	// Set up the right command.
	rightCmd := exec.Command(right[0], rightArgs...)
	rightCmd.Stdin = r
	rightCmd.Stdout = os.Stdout
	rightCmd.Stderr = os.Stderr

	if rightOut != nil {
		rightCmd.Stdout = rightOut
	}
	if rightErr != nil {
		rightCmd.Stderr = rightErr
	}

	defer closeFile(leftOut)
	defer closeFile(leftErr)
	defer closeFile(rightOut)
	defer closeFile(rightErr)

	// Start the left process first.
	if err := leftCmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return true
	}

	// Start the right process without waiting for the left.
	if err := rightCmd.Start(); err != nil {
		_ = w.Close()
		_ = r.Close()
		_ = leftCmd.Wait()
		fmt.Fprintln(os.Stderr, err)
		return true
	}

	// The parent must release its copies of the pipe handles.
	// Otherwise, the reader might never receive EOF.
	_ = w.Close()
	_ = r.Close()

	// Wait for both processes.
	leftErrWait := leftCmd.Wait()
	rightErrWait := rightCmd.Wait()

	// Child processes normally print their own errors.
	// Don't print exec.ExitError again.
	for _, waitErr := range []error{leftErrWait, rightErrWait} {
		var exitErr *exec.ExitError
		if waitErr != nil && !errors.As(waitErr, &exitErr) {
			fmt.Fprintln(os.Stderr, waitErr)
		}
	}

	return true
}

func closeFile(file *os.File) {
	if file != nil {
		_ = file.Close()
	}
}
