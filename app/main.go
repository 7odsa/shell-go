package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	// "math/rand"
	// "sync"
	// "time"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print
var commands = []string{"exit", "echo", "type"}

func main() {

	reader := bufio.NewReader(os.Stdin)
	for {

		fmt.Print("$ ")

		prompt, err := reader.ReadString('\n')

		handleErr(err)
		prompt = strings.TrimSpace(prompt)

		var tokens []string = strings.Split(prompt, " ")

		hasCommand(tokens)
	}

}

func hasCommand(tokens []string) bool {
	switch {
	case tokens[0] == "":
		return false
	case isExit(tokens[0]):
		handleExit()
		return true
	case isEcho(tokens[0]):
		handleEcho(tokens)
		return true
	case isType(tokens[0]):
		handleType(tokens)
		return true
	default:
		fmt.Println(strings.Join(tokens, "") + ": command not found")
		return false
	}

}

func handleType(tokens []string) {
	if len(tokens) <= 1 {
		fmt.Println("")
		return
	}

	target := tokens[1]

	if slices.Contains(commands, target) {
		fmt.Println(target + " is a shell builtin")
	} else {
		searchInPaths(tokens)
	}

}
func searchInPaths(tokens []string) {
	target := tokens[1]

	if path, err := exec.LookPath(target); err == nil {
		fmt.Println(target + " is " + path)
		return
	}

	fmt.Println(target + ": not found")

}

func handleEcho(tokens []string) {
	if len(tokens) > 1 {
		fmt.Println(strings.Join(tokens[1:], " "))
	} else {
		fmt.Println("nothing")

	}
}

func handleExit() {
	os.Exit(0)

}

func isType(prompt string) bool {
	return hasPrefixCommand(prompt, "type")
}
func isEcho(prompt string) bool {
	return hasPrefixCommand(prompt, "echo")
}
func hasPrefixCommand(command string, prefix string) bool {
	return strings.HasPrefix(command, prefix)
}
func isExit(command string) bool {
	return command == "exit"
}

func removeNewlineFromCommand(command string) string {
	if command[:len(command)-1] == "\n" {
		return command[:len(command)-1]
	} else {
		return command
	}

}

func handleErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
		os.Exit(1)

	}
}
