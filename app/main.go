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
	"strconv"
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
		isbBackground := false

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

		action := inputs[0]
		if inputs[len(inputs)-1] == "&" && action != "echo" {
			isbBackground = true
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
			if !isbBackground {

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

					select {
					case <-job.done:
						chanState = "Done"
					default:
						chanState = "Running"
					}

					jobInfo := "[" + strconv.Itoa(job.number) + "]+ <PID>" // jobInfo := "[" + strconv.Itoa(job.number) + "]+  " + chanState
					// var chanState string
					// sizeOfJobinfo := len(jobInfo)
					// for i := 0; i < 27-sizeOfJobinfo; i++ {
					// 	jobInfo = jobInfo + " "
					// }
					// jobInfo = jobInfo + job.inputLine
					fmt.Println(jobInfo)
					go func() {
						cmd.Wait()
						close(job.done)
					}()

				}
			}

		} else {
			fmt.Fprintln(os.Stderr, action+": command not found")
		}

		// Clean up opened file
		if outFile != nil {
			outFile.Close()
		}
		if errFile != nil {
			errFile.Close()
		}
	}
}
