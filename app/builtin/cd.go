package builtin

import (
	"fmt"
	"os"
)

func HandleCd(tokens []string) {
	if len(tokens) <= 1 {
		return
	}
	newDirectory := tokens[1]

	if err := os.Chdir(newDirectory); err != nil {

		fmt.Println("cd: " + newDirectory + ": No such file or directory")
	}
}
