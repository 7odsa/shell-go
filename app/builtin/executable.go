package builtin

import (
	"os"
	"os/exec"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

func HandleExecutable(tokens []string, prompt string) {

	if utils.HandleRedirectStdout(tokens[1:], "") {
		return
	}

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
