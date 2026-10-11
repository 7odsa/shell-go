package main

type Comp struct {
	compScripts map[string]string
	isCompleter func(comp Comp, word string) bool
}

var comp Comp = Comp{compScripts: map[string]string{}, isCompleter: func(comp Comp, word string) bool {
	_, isComp := comp.compScripts[word]
	return isComp
}}
