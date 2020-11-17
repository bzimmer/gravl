package console

import "strings"

type Accumulator struct {
	lines []string
}

func NewAccumulator() *Accumulator {
	return &Accumulator{
		lines: make([]string, 0),
	}
}

func (a *Accumulator) String() string {
	if len(a.lines) == 0 {
		return ""
	}
	return strings.Join(a.lines, "")
}

func (a *Accumulator) Push(line string) {
	a.lines = append(a.lines, line)
}

func (a *Accumulator) Reset() {
	a.lines = make([]string, 0)
}
