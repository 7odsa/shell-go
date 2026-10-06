package builtin

import (
	"fmt"
	"slices"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

func HandleType(tokens []string, builtIns []string) {

	target := tokens[1]

	if slices.Contains(builtIns, target) && target != "cat" {
		fmt.Println(target + " is a shell builtin")
	} else {
		utils.SearchInPaths(tokens)
	}

}
