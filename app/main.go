package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/app/builtin"
	"github.com/codecrafters-io/shell-starter-go/app/utils"
	// "math/rand"
	// "sync"
	// "time"
)

var builtIns = []string{"exit", "echo", "type", "pwd", "cd", "cat"}

func main() {

	autoCompleter := readline.NewPrefixCompleter(
		readline.PcItem("echo"),
		readline.PcItem("exit"),
	)

	reader, err := readline.NewEx(&readline.Config{
		Prompt:       "$ ",
		AutoComplete: autoCompleter,
	})

	if err != nil {
		panic(err)
	}

	defer reader.Close()

	for {

		fmt.Print("$ ")

		prompt, err := reader.Readline()

		utils.HandleErr(err)

		prompt = strings.TrimSpace(prompt)
		if utils.IsEmptyPrompt(prompt) {
			continue
		}

		var tokens []string = utils.HandleQuotes(prompt)

		hasCommand(tokens, prompt)
	}

}

func hasCommand(tokens []string, prompt string) bool {
	// if len(tokens) > 1 {
	// 	tokens[1] = strings.TrimLeft(tokens[1], " ")

	// }

	switch {

	case utils.IsBuiltIn(tokens[0], "exit"):
		os.Exit(0)
		return true
	case utils.IsBuiltIn(tokens[0], "pwd"):
		builtin.HandlePwd(tokens)
		return true

	case utils.IsBuiltIn(tokens[0], "echo"):
		builtin.HandleEcho(tokens, strings.TrimLeft(strings.TrimSuffix(prompt, "echo"), " "))
		return true
	// case utils.IsBuiltIn(tokens[0], "cat"):
	// 	builtin.HandleCat(tokens, strings.TrimLeft(strings.TrimSuffix(prompt, "cat"), " "))
	// 	return true
	case utils.IsBuiltIn(tokens[0], "type"):
		builtin.HandleType(tokens, builtIns)
		return true

	case utils.IsBuiltIn(tokens[0], "cd"):
		builtin.HandleCd(tokens)
		return true

	case utils.IsExecutable(tokens, prompt):
		builtin.HandleExecutable(tokens, prompt)
		return true
	default:
		fmt.Println(strings.Join(tokens, "") + ": command not found")
		return false
	}

}
