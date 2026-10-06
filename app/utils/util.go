package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func HandleQuotes(prompt string) []string {
	var word string = ""
	var res []string

	var isSingleQuote bool = false
	var isDoubleQuote bool = false
	var isFirstSpace bool = true

	for i := 0; i < len(prompt); i++ {
		char := prompt[i]
		switch char {
		case '\\':
			isFirstSpace = true

			for j := i; j < len(prompt)-2 && prompt[j] == '\\'; j += 2 {
				char = prompt[j+1]
				word += string(char)
				if prompt[j+2] != '\\' {
					i = j + 1
					break
				}

			}

		case '"':
			isDoubleQuote = !isDoubleQuote
			if word != "" {
				res = append(res, word)

			}
			word = ""
		case '\'':
			if !isDoubleQuote {
				isSingleQuote = !isSingleQuote
				if word != "" {
					res = append(res, word)

				}
				word = ""
			} else {
				word += string(char)
			}
		default:
			if isSingleQuote || isDoubleQuote || char != ' ' {
				isFirstSpace = true
				word += string(char)

			} else if char == ' ' && isFirstSpace {
				isFirstSpace = false
				word += string(char)
			}
		}

	}
	if word != "" {
		res = append(res, word)

	}

	return res
}

func SearchInPaths(tokens []string) {
	target := tokens[1]

	if path, err := exec.LookPath(target); err == nil {
		fmt.Println(target + " is " + path)
		return
	}

	fmt.Println(target + ": not found")

}

func IsExecutable(tokens []string) bool {
	target := tokens[0]

	if _, err := exec.LookPath(target); err == nil {
		return true
	}
	return false
}
func IsEmptyPrompt(tokens []string) bool {
	return tokens[0] == ""
}

// func isType(prompt string) bool {
// 	return hasPrefixCommand(prompt, "type")
// }

func HasPrefixCommand(command string, prefix string) bool {
	return strings.HasPrefix(command, prefix)
}

func IsBuiltIn(command string, builtIn string) bool {
	return (command == builtIn)
}

func RemoveNewlineFromCommand(command string) string {
	if command[:len(command)-1] == "\n" {
		return command[:len(command)-1]
	} else {
		return command
	}

}
func IsOneWord(tokens []string) bool {
	return len(tokens) <= 1
}
func HandleErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
		os.Exit(1)

	}
}
