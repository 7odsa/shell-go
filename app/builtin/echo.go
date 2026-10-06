package builtin

import (
	"fmt"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

func HandleEcho(tokens []string, prompt string) {
	if utils.IsOneWord(tokens) {
		fmt.Println("nothing")
	}
	fmt.Println(handleSingleQuotes(prompt))

}
func handleSingleQuotes(prompt string) string {
	var word string = ""
	// var res []string

	var isSingleQuote bool = false
	var isFirstSpace bool = true

	for _, char := range prompt {

		if char == '\'' {
			isSingleQuote = !isSingleQuote
			if !isSingleQuote {
				isFirstSpace = true
				// res = append(res, word)
				// word = ""
			}
		} else {
			if (isSingleQuote) || char != ' ' {
				isFirstSpace = true
				word += string(char)

			} else if char == ' ' && isFirstSpace {
				isFirstSpace = false
				word += string(char)
			}
		}

	}
	return word
}
