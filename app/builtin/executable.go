package builtin

import (
	"os"
	"os/exec"
	"strings"
)

func HandleExecutable(tokens []string, prompt string) {

	for i, s := range tokens {
		tokens[i] = strings.Trim(s, " ")
	}

	target := strings.Trim(tokens[0], " ")

	if _, err := exec.LookPath(target); err == nil {

		cmd := exec.Command(target, tokens[1:]...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}

}
