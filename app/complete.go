package main

type Comp struct {
	compScripts map[string]string
	isCompleter func(comp Comp, word string) (string, bool)
}

var comp Comp = Comp{compScripts: map[string]string{}, isCompleter: func(comp Comp, word string) (string, bool) {
	compScrpt, isComp := comp.compScripts[word]
	return compScrpt, isComp
}}
