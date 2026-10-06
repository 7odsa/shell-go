package builtin

import (
	"fmt"
	"os"
)

func HandleCd(tokens []string) {

	newDir := tokens[1]

	if newDir == "~" {
		newDir, _ = os.UserHomeDir()
	}

	if err := os.Chdir(newDir); err != nil {
		fmt.Println("cd: " + newDir + ": No such file or directory")
	}
}
