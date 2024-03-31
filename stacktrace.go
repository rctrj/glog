package log

import (
	"fmt"
	"path/filepath"
	"runtime"
)

type callFrame struct {
	File     string
	Function string
	Line     int
}

func callers(skip int, levels int) (stack []callFrame) {
	stack = make([]callFrame, 0)
	callers := make([]uintptr, levels)
	n := runtime.Callers(skip+2, callers[:])

	frames := runtime.CallersFrames(callers[:n])
	for {
		frame, more := frames.Next()
		stack = append(stack, callFrame{
			File:     frame.File,
			Function: frame.Function,
			Line:     frame.Line,
		})

		if !more {
			break
		}
	}

	return
}

func (c callFrame) prettyString() string {
	return fmt.Sprintf("%s @%s: %d", c.File, filepath.Base(c.Function), c.Line)
}
