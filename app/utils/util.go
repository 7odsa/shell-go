package utils

import (
	"fmt"
	"log"
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

	var isBackslach bool = false

	for _, char := range prompt {

		if isBackslach {
			isBackslach = !isBackslach
			word += string(char)
			continue
		}

		switch char {
		case '\\':
			if isSingleQuote {
				isFirstSpace = true
				word += string(char)

			} else {
				isBackslach = !isBackslach
			}
			// isFirstSpace = true
			// for j := i; j < len(prompt)-2 && prompt[j] == '\\'; j += 2 {
			// 	char = prompt[j+1]
			// 	word += string(char)
			// 	if prompt[j+2] != '\\' {
			// 		if char == ' ' {
			// 			isFirstSpace = false
			// 		}
			// 		i = j + 1
			// 		break
			// 	}
			// }
		case '"':

			// if isSingleQuote || isDoubleQuote || char != ' ' {
			// 	isFirstSpace = true
			// 	word += string(char)
			// 	break
			// }
			if !isSingleQuote {
				isDoubleQuote = !isDoubleQuote

			} else {
				word += string(char)

			}
			// if word != "" {
			// 	res = append(res, word)
			// }
			// word = ""
		case '\'':

			if !isDoubleQuote {
				isSingleQuote = !isSingleQuote
				// if word != "" {
				// 	// println(string(word[len(word)-1]))

				// 	res = append(res, word)

				// }
				// word = ""
			} else {
				word += string(char)
			}
		default:

			isBackslach = false
			if isSingleQuote || isDoubleQuote || char != ' ' {
				isFirstSpace = true
				word += string(char)

			} else if char == ' ' && isFirstSpace {
				isFirstSpace = false
				if word != "" {
					res = append(res, word)
				}
				word = ""
			}
		}

	}
	if word != "" {
		res = append(res, word)

	}

	// println(strings.Join(res,"&&"))
	// println(strings.Join(res, "&&"))

	return res
}

func HandleRedirectStdout(tokens []string, prompt string) bool {

	indexOp := -1
	for i, s := range tokens {
		if strings.Trim(s, " ") == ">" || strings.Trim(s, " ") == "1>" {
			indexOp = i
			break
		}
	}
	if indexOp == -1 {
		return false
	}

	fileName := tokens[indexOp+1]
	args := tokens[:indexOp]

	cmd := exec.Command(args[0], args[1:]...)

	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer file.Close()

	cmd.Stdout = file
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
}

func SearchInPaths(tokens []string) {
	target := tokens[1]

	if path, err := exec.LookPath(target); err == nil {
		fmt.Println(target + " is " + path)
		return
	}

	fmt.Println(target + ": not found")

}

func IsExecutable(tokens []string, prompt string) bool {

	target := strings.Trim(tokens[0], " ")

	if _, err := exec.LookPath(target); err == nil {
		return true
	}
	return false
}
func IsEmptyPrompt(prompt string) bool {
	return prompt == ""
}

// func isType(prompt string) bool {
// 	return hasPrefixCommand(prompt, "type")
// }

func HasPrefixCommand(command string, prefix string) bool {
	return strings.HasPrefix(command, prefix)
}

func IsBuiltIn(command string, builtIn string) bool {
	target := strings.Trim(command, " ")
	return (target == builtIn)
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
