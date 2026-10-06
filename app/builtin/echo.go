package builtin

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

func HandleEcho(tokens []string, prompt string) {
	if utils.IsOneWord(tokens) {
		fmt.Println("")
		return
	}

	res := strings.Join(utils.HandleSingleQuotes(prompt), "")
	fmt.Println(res)

}
