package builtin

import (
	"fmt"
	"strings"
)

func HandleEcho(tokens []string, prompt string) {

	res := strings.Join(tokens, "")
	fmt.Println(res)

}
