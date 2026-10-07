package builtin

import (
	"fmt"
	"strings"
)

func HandleEcho(tokens []string, prompt string) {

	res := strings.Join(tokens, "")
	res = strings.TrimLeft(res, " ")
	fmt.Println(res)

}
