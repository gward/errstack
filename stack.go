package errstack

import (
	"fmt"
	"runtime"
)

type StackTrace []uintptr

func Callers(skip int) StackTrace {
	stack := make([]uintptr, 32)
	depth := runtime.Callers(skip + 1, stack)
	return StackTrace(stack[0:depth])
}

func (stack StackTrace) FormatStack(msg string, stopFunction *string) []string {
	lines := make([]string, 0, 1 + 2*len(stack))
	if msg != "" {
		lines = append(lines, msg + ":")
	}
	frames := runtime.CallersFrames(stack)
	for {
		frame, more := frames.Next()

		lines = append(lines, fmt.Sprintf("%s()", frame.Function))
		lines = append(lines, fmt.Sprintf("\t%s:%d", frame.File, frame.Line))
		if stopFunction != nil && frame.Function == *stopFunction {
			break
		}
		if !more {
			break
		}
	}
	return lines
}
