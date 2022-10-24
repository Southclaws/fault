package fault

import (
	"fmt"
	"runtime"
	"strings"
)

type Stack []StackFrame

func (s Stack) String() string {
	var str []string
	for _, f := range s {
		str = append(str, f.String())
	}
	return strings.Join(str, "\n")
}

type StackFrame struct {
	File string
	Line int
}

func (f *StackFrame) String() string {
	return fmt.Sprintf("%s:%d", f.File, f.Line)
}

func callers(skip int) *stack {
	const depth = 64
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	var st stack = pcs[0 : n-2] // todo: change this to filtering out runtime instead of hardcoding n-2
	return &st
}

type stack []uintptr

func (s *stack) get() Stack {
	var stackFrames []StackFrame

	frames := runtime.CallersFrames(*s)
	for {
		frame, more := frames.Next()
		stackFrames = append(stackFrames, StackFrame{
			File: frame.File,
			Line: frame.Line,
		})
		if !more {
			break
		}
	}

	return stackFrames
}
