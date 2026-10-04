package main

import (
	"bufio"
	"fmt"
	"os"
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
		if command == "exit" {
			os.Exit(0)
		}
		fmt.Println(command + ": command not found")
	}

}
