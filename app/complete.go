package main

type Comp struct {
	compScripts map[string]string
}

var comp Comp = Comp{compScripts: map[string]string{}}

func (comp *Comp) isCompleter(word string) (string, bool) {
	compScrpt, isComp := comp.compScripts[word]
	return compScrpt, isComp
}
