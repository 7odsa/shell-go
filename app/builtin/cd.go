package builtin

import (
	"fmt"
	"os"
)

var PrevWd string

func HandleCd(tokens []string) {
	if len(tokens) <= 1 {
		return
	}
	newDir := tokens[1]

	if newDir == "~" {
		newDir = PrevWd
	}

	wd, _ := os.Getwd()
	PrevWd = wd

	if err := os.Chdir(newDir); err != nil {
		fmt.Println("cd: " + newDir + ": No such file or directory")
	}
}
