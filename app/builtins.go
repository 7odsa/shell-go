package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

type Command struct {
	exec   func(cmd Command, args []string)
	Stdout io.Writer
	Stderr io.Writer
	stdin  io.Reader
}

type Job struct {
	number    int
	cmd       *exec.Cmd
	done      chan struct{}
	inputLine string
}

var jobs []*Job

func (c Command) Run(args []string) {
	c.exec(c, args)
}

var BuiltinNames = []string{"echo", "type", "pwd", "cd", "exit", "complete", "jobs", "history", "declare"}

var BuiltinCmds = map[string]Command{
	"echo":     {exec: handleEcho},
	"type":     {exec: handleType},
	"pwd":      {exec: handlePwd},
	"cd":       {exec: handleCd},
	"exit":     {exec: handleExit},
	"complete": {exec: handleComplete},
	"jobs":     {exec: handleJobs},
	"history":  {exec: handleHistory},
	"declare":  {exec: handleDeclare},
}

var declareVars map[string]string = make(map[string]string)

func handleDeclare(cmd Command, args []string) {
	if len(args) <= 0 {
		return
	}

	if args[0] == "-p" && len(args) > 1 {
		if val, isVar := declareVars[args[1]]; !isVar {
			fmt.Fprintln(cmd.Stdout, "declare: "+args[1]+": not found")
		} else {
			fmt.Fprintln(cmd.Stdout, "declare -- "+args[1]+"=\""+val+"\"")
		}
	}
	if i := strings.Index(args[0], "="); i != -1 {
		if !(isLetter(args[0][0]) || args[0][0] == '_') {
			fmt.Fprintln(cmd.Stdout, "declare: `"+args[0]+"': not a valid identifier")
			return
		}
		key := args[0][:i]
		val := args[0][i+1:]
		declareVars[key] = val
	}
}
func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func handleHistory(cmd Command, args []string) {

	limit := len(utils.SavedHistory)
	if len(args) != 0 {

		if len(args) > 1 {
			switch args[0] {
			case "-r":
				utils.SavedHistory = append(utils.SavedHistory, utils.ReadDataFromFile(args[1])...)
				return
			case "-w":
				data := strings.Join(utils.SavedHistory, "\n")
				data = data + "\n"
				utils.AppendDataToaFile([]byte(data), args[1])
				return
			case "-a":
				data := strings.Join(utils.SavedHistory[utils.HistoryAppended:], "\n")
				data = data + "\n"
				utils.AppendDataToaFile([]byte(data), args[1])
				utils.HistoryAppended = len(utils.SavedHistory)
				return
			}
		}

		num, err := strconv.Atoi(args[0])

		if err == nil {
			limit = num
		}

	}

	for n := len(utils.SavedHistory) - limit; n < len(utils.SavedHistory); n++ {
		fmt.Fprint(cmd.Stdout, "    "+strconv.Itoa(n+1)+"  "+utils.SavedHistory[n]+"\n")

	}

}

func handleJobs(cmd Command, args []string) {

	printJobs(false)

}

func printJobs(doneOnly bool) {
	for i := 0; i < len(jobs); i++ {
		job := jobs[i]

		var chanState string
		select {
		case <-job.done:
			chanState = "Done"
		default:
			chanState = "Running"
		}
		if chanState == "Running" && doneOnly {
			continue
		}

		leadingChar := " "
		if i == len(jobs)-1 {
			leadingChar = "+"
		}
		if i == len(jobs)-2 {
			leadingChar = "-"
		}

		jobInfo := "[" + strconv.Itoa(job.number) + "]" + leadingChar + "  " + chanState
		sizeOfJobinfo := len(jobInfo)
		for i := 0; i < 27-sizeOfJobinfo; i++ {
			jobInfo = jobInfo + " "
		}
		// job.number -= numOfDeleted

		if chanState == "Done" {
			jobInfo = jobInfo + strings.TrimSuffix(job.inputLine, "&")
			jobs = slices.Delete(jobs, i, i+1)
			i--
		} else {
			jobInfo = jobInfo + job.inputLine
		}
		fmt.Println(jobInfo)

	}
}

func handleComplete(cmd Command, args []string) {
	if len(args) == 0 {
		return
	}

	flag := args[0]

	switch flag {
	case "-p":
		if len(args) < 2 {
			return
		}
		command := args[1]
		if comp, isComp := comp.compScripts[command]; isComp {
			println("complete -C '" + comp + "' " + command)
		} else {
			println("complete: " + command + ": no completion specification")
		}

	case "-C":
		if len(args) < 2 {
			return
		}
		completer, command := strings.Join(args[1:len(args)-1], " "), args[len(args)-1]
		comp.compScripts[command] = completer
		// To Be Done
	case "-r":
		if len(args) < 2 {
			return
		}
		command := strings.TrimSpace(args[len(args)-1])
		delete(comp.compScripts, command)
		// To Be Done
	}

}

func handleCd(cmd Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(cmd.Stderr, "cd needs at least 1 argument")
		return
	}
	path := args[0]
	if subPath, found := strings.CutPrefix(path, "~"); found {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(cmd.Stderr, err)
			return
		}
		path = homeDir + subPath
	}
	fileInfo, err := os.Stat(path)
	if err != nil || !fileInfo.IsDir() {
		fmt.Fprintf(cmd.Stderr, "cd: %s: No such file or directory\n", path)
		return
	}
	if err := os.Chdir(path); err != nil {
		fmt.Fprintf(cmd.Stderr, "cd: %s: %v\n", path, err)
	}
}

func handleEcho(cmd Command, args []string) {
	fmt.Fprintln(cmd.Stdout, strings.Join(args, " "))
}

func handleType(cmd Command, args []string) {
	if len(args) == 0 {
		return
	}
	arg := args[0]
	if slices.Contains(BuiltinNames, arg) {
		fmt.Fprintln(cmd.Stdout, arg+" is a shell builtin")
	} else if absPath, err := exec.LookPath(arg); err == nil {
		fmt.Fprintln(cmd.Stdout, arg+" is "+absPath)
	} else {
		fmt.Fprintln(cmd.Stdout, arg+": not found")
	}
}

func handlePwd(cmd Command, args []string) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(cmd.Stderr, "pwd error: ", err)
	} else {
		fmt.Fprintln(cmd.Stdout, dir)
	}
}

func handleExit(cmd Command, args []string) {
	deferClosing()
	os.Exit(0)
}
