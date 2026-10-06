package builtin

import (
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

func HandleCat(tokens []string, prompt string) {
	if utils.IsOneWord(tokens) {
		fmt.Println("nothing")
	}
	handledTokens := utils.HandleSingleQuotes(prompt)

	for _, token := range handledTokens {
		if strings.Trim(token, " ") == "" {
			continue
		}
		token = strings.Trim(string(token), " ")
		content, err := os.ReadFile(token)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Print(string(content))
	}
	fmt.Println()

}
