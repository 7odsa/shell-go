package builtin

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

func HandleEcho(tokens []string, prompt string) {

	// slices.Contains(tokens, ">")
	// slices.Contains(tokens, "1>")
	// for i, _ := range tokens {
	// 	fmt.Println(tokens[i])
	// }

	if utils.HandleRedirectStdout(tokens, "") {
		return
	}

	res := strings.Join(tokens[1:], " ")

	fmt.Println(res)

}
