package log

import (
	"fmt"
	"strings"
)

type Config struct {
	Level         Level
	StripLevel    StripLevel
	PrettyPrint   bool
	DisplayPrefix bool
}

type Logger struct {
	conf Config
	core core

	tag         string
	prefix      string
	printPrefix string
}

func New(conf Config, core core) *Logger {
	return newLogger(conf, core)
}

func (l *Logger) NewWithTag(tag string) *Logger {
	logger := newLogger(l.conf, l.core)
	logger.SetTag(tag)
	return logger
}

func (l *Logger) NewWithPrefix(prefix string) *Logger {
	logger := newLogger(l.conf, l.core)
	logger.SetTag(l.tag)
	logger.SetPrefix(prefix)
	return logger
}

func (l *Logger) NewWithTagAndPrefix(tag, prefix string) *Logger {
	logger := newLogger(l.conf, l.core)
	logger.SetTag(tag)
	logger.SetPrefix(prefix)
	return logger
}

func (l *Logger) NewWithoutTag() *Logger {
	return newLogger(l.conf, l.core)
}

// ResetPrefix you can use this as part of defer
// Example:
//
//	func foo(logger *log.Logger) {
//		defer logger.ResetPrefix()
//		logger.setPrefix("foo")
//		... do something here, and all logs will include the prefix
//	}
//
// Note: Be extremely careful using prefix in a class with async methods
func (l *Logger) ResetPrefix() {
	l.SetPrefix("")
}

func (l *Logger) SetTag(tag string) {
	l.updatePrintPrefix(tag, l.prefix)
}

func (l *Logger) SetPrefix(prefix string) {
	l.updatePrintPrefix(l.tag, prefix)
}

func (l *Logger) updatePrintPrefix(tag, prefix string) {
	tag = trim(tag)
	prefix = trim(prefix)

	printPrefix := ""
	if tag != "" {
		printPrefix = fmt.Sprintf("[%s] ", tag)
	}

	if prefix != "" {
		printPrefix = fmt.Sprintf("%s%s: ", printPrefix, prefix)
	}

	l.tag = tag
	l.prefix = prefix
	l.printPrefix = printPrefix
}

func trim(str string) string {
	return strings.Trim(str, " ")
}

func newLogger(conf Config, core core) *Logger {
	return &Logger{
		conf: conf,
		core: core,
	}
}
