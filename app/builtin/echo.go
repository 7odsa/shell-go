package builtin

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

func HandleEcho(tokens []string) {
	if utils.IsOneWord(tokens) {
		fmt.Println("nothing")
	}

	fmt.Println(strings.Join(tokens[1:], " "))
}
