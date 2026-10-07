package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/builtin"
	"github.com/codecrafters-io/shell-starter-go/app/utils"
	// "math/rand"
	// "sync"
	// "time"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print
var builtIns = []string{"exit", "echo", "type", "pwd", "cd", "cat"}

func main() {

	reader := bufio.NewReader(os.Stdin)
	for {

		fmt.Print("$ ")

		prompt, err := reader.ReadString('\n')

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
	switch {

	case utils.IsBuiltIn(tokens[0], "exit"):
		os.Exit(0)
		return true
	case utils.IsBuiltIn(tokens[0], "pwd"):
		builtin.HandlePwd(tokens)
		return true

	case utils.IsBuiltIn(tokens[0], "echo"):
		builtin.HandleEcho(tokens[1:], strings.TrimLeft(strings.TrimPrefix(prompt, "echo"), " "))
		return true
	case utils.IsBuiltIn(tokens[0], "cat"):
		builtin.HandleCat(tokens[1:], strings.TrimLeft(strings.TrimPrefix(prompt, "cat"), " "))
		return true
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
