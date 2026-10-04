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

func main() {

	reader := bufio.NewReader(os.Stdin)
	for {

		fmt.Print("$ ")
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)

		}
		command = command[:len(command)-1]
		checkExit(command)

		if strings.HasPrefix(command, "echo") {
			if len(command) > 4 {
				fmt.Println(command[5:])
			}
		} else {
			fmt.Println(command + ": command not found")
		}
	}

}

func checkExit(command string) {
	if command == "exit" {
		os.Exit(0)
	}
}
