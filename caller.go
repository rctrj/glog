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

func callers(skip int, levels int) (caller callFrame, stack []callFrame) {
	stack = make([]callFrame, 0)
	rpc := make([]uintptr, levels)
	n := runtime.Callers(skip+2, rpc[:])

	frames := runtime.CallersFrames(rpc[:n])
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

	if len(stack) > 0 {
		caller = stack[0]
	}
	return
}

func (cf callFrame) prettyString() string {
	return fmt.Sprintf("%s @%s: %d", cf.File, filepath.Base(cf.Function), cf.Line)
}
