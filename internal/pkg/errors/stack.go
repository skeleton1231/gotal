package errors

import (
	"fmt"
	"runtime"
)

// callers returns a stack trace captured at the point it's called.
func callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first few callers

	st := &stack{pcs[0:n]}
	return st
}

// stack represents a stack of program counters.
type stack struct {
	pcs []uintptr
}

// Format formats the stack trace.
func (s *stack) Format(st fmt.State, verb rune) {
	for _, pc := range s.pcs {
		f := runtime.FuncForPC(pc)
		if f == nil {
			continue
		}
		file, line := f.FileLine(pc)
		fmt.Fprintf(st, "\n\t%s:%d", file, line)
	}
}
