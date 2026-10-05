package builtin

import (
	"fmt"
	"os"
)

func HandlePwd(tokens []string) {
	path, _ := os.Getwd()

	fmt.Println(path)
}
