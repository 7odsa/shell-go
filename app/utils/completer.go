package utils

type Completer struct {
	BuiltIns []string
}

func (c *Completer) Do(line []rune, pos int) ([][]rune, int) {
	println(string(line))
	println(pos)

	return nil, 0
}
