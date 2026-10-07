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

	res := strings.Join(tokens, "")

	if utils.HandleRedirectStdout(tokens, res) {
		return
	}

	fmt.Println(res)

}
