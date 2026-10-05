package builtin

import (
	"fmt"
	"slices"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

func HandleType(tokens []string, builtIns []string) {
	if utils.IsOneWord(tokens) {
		fmt.Println("")
		return
	}

	target := tokens[1]

	if slices.Contains(builtIns, target) {
		fmt.Println(target + " is a shell builtin")
	} else {
		utils.SearchInPaths(tokens)
	}

}
