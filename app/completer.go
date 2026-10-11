package main

import (
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/chzyer/readline"
)

type BuiltinCompleter struct {
	prompt  string
	lastTab string // prefix of the previous ambiguous TAB; second TAB on it lists matches
}

func NewCustomCompleter(prompt string) *BuiltinCompleter {
	return &BuiltinCompleter{prompt: prompt}
}

// Do prints the match list itself and returns nothing for readline to draw,
// because readline's own menu renders below a redrawn prompt, not above it.
func (c *BuiltinCompleter) Do(line []rune, pos int) ([][]rune, int) {
	s := strings.Split(string(line[:pos]), " ")
	prefix := s[len(s)-1]

	if compScrpt, isComp := comp.isCompleter(comp, s[0]); isComp {

		cmd := exec.Command(compScrpt)

		output, err := cmd.Output()
		println(string(output))
		if err != nil {
			return suffixOf(string(output), "")
		}

	}

	var matches []string
	var isCommand bool = true
	if len(s) <= 1 {
		matches = listCommands(prefix, os.Getenv("PATH"))
	} else {
		isCommand = false
		matches, prefix = listFiles(prefix)
	}

	switch {
	case len(matches) == 0:
		return c.ring()
	case len(matches) == 1:
		if isCommand || strings.Contains(matches[0], ".") {
			return suffixOf(matches[0]+" ", prefix)
		} else {
			return suffixOf(matches[0], prefix)
		}
	}

	if common := commonPrefix(matches); len(common) > len(prefix) {

		if isCommand || strings.Contains(common, ".") {
			return suffixOf(common, prefix)
		} else {
			return suffixOf(common, prefix)
		}

	}

	if c.lastTab != prefix {
		c.lastTab = prefix
		return c.ring()
	}

	fmt.Fprintf(readline.Stdout, "\r\n%s\r\n%s%s", strings.Join(matches, "  "), c.prompt, string(line[:pos]))
	return nil, 0
}

func (c *BuiltinCompleter) ring() ([][]rune, int) {
	readline.Stdout.Write([]byte("\x07"))
	return nil, 0
}

func suffixOf(full, prefix string) ([][]rune, int) {
	return [][]rune{[]rune(strings.TrimPrefix(full, prefix))}, len(prefix)
}

func listCommands(prefix, pathEnv string) []string {
	names := map[string]struct{}{}
	for name := range BuiltinCmds {
		names[name] = struct{}{}
	}
	for _, name := range executableNames(pathEnv) {
		names[name] = struct{}{}
	}

	matches := slices.DeleteFunc(slices.Collect(maps.Keys(names)), func(name string) bool {
		return !strings.HasPrefix(name, prefix)
	})
	slices.Sort(matches)
	return matches
}

func listFiles(prefix string) ([]string, string) {

	dir, file := "", prefix
	if strings.Contains(prefix, "/") {
		dir, file = filepath.Split(prefix)
		// println("&&")
	}

	names := map[string]struct{}{}

	for _, name := range fileNames(dir) {
		if !strings.Contains(name, ".") {
			names[name+"/"] = struct{}{}
		} else {
			names[name] = struct{}{}

		}
	}

	matches := slices.DeleteFunc(slices.Collect(maps.Keys(names)), func(name string) bool {
		return !strings.HasPrefix(name, file)
	})

	slices.Sort(matches)
	return matches, file
}

func fileNames(dir string) []string {
	if dir == "" {
		dir = "."
	}
	dirEnt, _ := os.ReadDir(dir)

	var names []string
	for _, entry := range dirEnt {
		names = append(names, (entry.Name()))

	}

	return names
}

func executableNames(pathEnv string) []string {
	var names []string
	for _, dir := range filepath.SplitList(pathEnv) {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			info, err := entry.Info()
			if err != nil || info.IsDir() || info.Mode().Perm()&0111 == 0 {
				continue
			}
			names = append(names, entry.Name())
		}
	}
	return names
}

func commonPrefix(sortedMatches []string) string {
	first, last := sortedMatches[0], sortedMatches[len(sortedMatches)-1]
	i := 0
	for i = range len(first) {
		if first[i] != last[i] {
			break
		}
		i++
	}
	return first[:i]
}
