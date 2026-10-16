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

		// FROM HERE WE WILL CHECK IF THERE IS | OR |& AND IF ANY OF THEM IS THERE WE WILL GET THE STDOUT FROM THE FIRST COMMAND
		// AND PASS IT AS A STDIN TO THE SECCOND COMMAND
		command := Command{}
		command.stdin = os.Stdin
		command.Stdout = os.Stdout
		command.handleCommand(inputs, line)

		printJobs(true)
	}
}

func (command *Command) handleCommand(inputs []string, line string) {

	if i := slices.Index(inputs, "|"); i != -1 {
		r, w, err := os.Pipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		leftInputs := inputs[:i]
		rightInputs := inputs[i+1:]

		if len(leftInputs) == 0 || len(rightInputs) == 0 {
			r.Close()
			w.Close()
			fmt.Fprintln(os.Stderr, "invalid pipeline")
			return
		}

		// For now, support external commands in pipelines.
		leftCmd := exec.Command(leftInputs[0], leftInputs[1:]...)
		leftCmd.Stdin = command.stdin
		leftCmd.Stdout = w
		leftCmd.Stderr = os.Stderr

		rightCmd := exec.Command(rightInputs[0], rightInputs[1:]...)
		rightCmd.Stdin = r
		rightCmd.Stdout = command.Stdout
		rightCmd.Stderr = os.Stderr

		// Start the left process first.
		if err := leftCmd.Start(); err != nil {
			r.Close()
			w.Close()
			fmt.Fprintln(os.Stderr, err)
			return
		}

		// Start the right process before waiting for the left.
		if err := rightCmd.Start(); err != nil {
			w.Close()
			r.Close()
			_ = leftCmd.Wait()
			fmt.Fprintln(os.Stderr, err)
			return
		}

		// Close the parent's copies so EOF can reach the reader.
		_ = w.Close()
		_ = r.Close()

		// Foreground pipeline: wait for both processes.
		leftErr := leftCmd.Wait()
		rightErr := rightCmd.Wait()

		// Non-zero exit statuses are normal command outcomes.
		// Report other errors if needed.
		_ = leftErr
		_ = rightErr

		return
	}

}
