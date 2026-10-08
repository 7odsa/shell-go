package builtin

import (
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

func HandleCat(tokens []string, prompt string) {

	if utils.HandleRedirectStdout(tokens, "") {
		println(prompt)
		return
	}

	for _, token := range tokens {
		if strings.Trim(token, " ") == "" {
			continue
		}
		token = strings.Trim(string(token), " ")
		content, err := os.ReadFile(token)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(string(content))
	}

}
