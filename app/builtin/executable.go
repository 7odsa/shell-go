package builtin

import (
	"os"
	"os/exec"
)

func HandleExecutable(tokens []string) {
	target := tokens[0]
	if _, err := exec.LookPath(target); err == nil {

		cmd := exec.Command(target, tokens[1:]...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}

}
