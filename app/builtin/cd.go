package builtin

import (
	"fmt"
	"os"
)

func HandleCd(tokens []string) {
	if len(tokens) <= 1 {
		return
	}
	newDir := tokens[1]

	if newDir == "~" {
		newDir, _ = os.UserHomeDir()
	}

	if err := os.Chdir(newDir); err != nil {
		fmt.Println("cd: " + newDir + ": No such file or directory")
	}
}
