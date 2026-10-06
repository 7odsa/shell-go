package builtin

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

func HandleCd(tokens []string) {
	if utils.IsOneWord(tokens) {
		fmt.Println("")
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
