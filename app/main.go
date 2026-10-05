package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	// "math/rand"
	// "sync"
	// "time"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print
var commands [2]string = [2]string{"exit", "echo"}

func main() {

	reader := bufio.NewReader(os.Stdin)
	for {

		fmt.Print("$ ")

		command, err := reader.ReadString('\n')

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)

		}

		command = removeNewlineFromCommand(command)

		if hasCommand(command) == false {
			fmt.Println(command + ": command not found")
		}
	}

}

func hasCommand(command string) bool {
	switch {
	case isExit(command):
		handleExit(command)
		return true
	case isEcho(command):
		handleEcho(command)
		return true
	case isType(command):
		handleType(command)
		return true
	}

	return false

}

func handleType(command string) {

	if len(command) <= 4 {
		fmt.Println(": nothing")
		return
	}
	switch {
	case (isExit(command) || isEcho(command)):
		fmt.Println(command[5:] + " is a shell builtin")
		return
	}
	fmt.Println(command[5:] + ": not found")

}

func handleEcho(command string) {
	if len(command) > 4 {
		fmt.Println(command[5:])
	}

}

func handleExit(command string) {
	os.Exit(0)

}

func isType(command string) bool {
	return hasPrefixCommand(command, "type")
}
func isEcho(command string) bool {
	return hasPrefixCommand(command, "echo")
}
func hasPrefixCommand(command string, prefix string) bool {
	return strings.HasPrefix(command, prefix)
}
func isExit(command string) bool {
	return command == "exit"
}

func removeNewlineFromCommand(command string) string {
	return command[:len(command)-1]

}
