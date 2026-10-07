package builtin

import (
	"os"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

func HandleExecutable(tokens []string, prompt string) {
	handledTokens := (utils.HandleQuotes(prompt))
	target := handledTokens[0]

	println(target)
	if _, err := exec.LookPath(target); err == nil {

		cmd := exec.Command(target, handledTokens[1:]...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}

}
