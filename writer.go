package log

import (
	"context"
	"fmt"
	"time"
)

func (l *Logger) Trace(ctx context.Context, msg string, entries ...*entry) {
	l.write(ctx, LevelTrace, msg, entries)
}

func (l *Logger) Debug(ctx context.Context, msg string, entries ...*entry) {
	l.write(ctx, LevelDebug, msg, entries)
}

func (l *Logger) Info(ctx context.Context, msg string, entries ...*entry) {
	l.write(ctx, LevelInfo, msg, entries)
}

func (l *Logger) Warn(ctx context.Context, msg string, entries ...*entry) {
	l.write(ctx, LevelWarn, msg, entries)
}

func (l *Logger) Error(ctx context.Context, msg string, err error, entries ...*entry) {
	l.write(ctx, LevelError, msg, addErr(err, entries))
}

func (l *Logger) Fatal(ctx context.Context, msg string, err error, entries ...*entry) {
	l.write(ctx, LevelFatal, msg, addErr(err, entries))
}

func (l *Logger) write(ctx context.Context, level Level, msg string, entries []*entry) {
	if l.conf.Level.greaterThan(level) {
		return
	}

	entriesFromCtx := extractEntriesFromCtx(ctx)
	extraParams := l.extraParams()
	entries = append(entries, entriesFromCtx...)
	entries = append(entries, extraParams...)

	msg = l.modifyMsg(msg, level)
	prunedFields := l.pruneEntries(entries)

	l.log(ctx, level, msg, prunedFields)
}

func (l *Logger) modifyMsg(msg string, level Level) string {
	if l.conf.DisplayPrefix {
		msg = fmt.Sprintf("%s%s", l.printPrefix, msg)
	}

	if !l.conf.PrettyPrint {
		return msg
	}

	return level.color().Add(msg)
}

func (l *Logger) pruneEntries(entries []*entry) (remaining []*entry) {
	sl := l.conf.StripLevel

	if sl.isMessageOnly() {
		return
	}

	for _, e := range entries {
		if e == nil {
			continue
		}

		if sl.isMinimumFields() && !e.isMinimalInfo {
			continue
		}

		remaining = append(remaining, e)
	}

	return
}

func (l *Logger) extraParams() []*entry {
	if !l.conf.StripLevel.isNone() {
		return []*entry{}
	}

	caller, callStack := l.getCallers()
	return []*entry{
		Entry(KeyLogTag, l.tag),
		Entry(KeyLogPrefix, l.prefix),
		Entry(KeyLogLevel, l.conf.Level),
		Entry(KeyCaller, caller),
		Entry(KeyCallStack, callStack),
		Entry(KeyTimestamp, time.Now()),
	}
}

func (l *Logger) log(ctx context.Context, level Level, message string, entries []*entry) {
	target := l.core.Log
	logEntries := make([]IEntry, len(entries))
	for i, e := range entries {
		logEntries[i] = e
	}
	target(ctx, level, message, logEntries)
}

func (l *Logger) getCallers() (caller interface{}, callStack interface{}) { //caller, callFrames
	mainCaller, callFrames := callers(4, 4)

	// string format to reduce space used
	frameString := make([]string, 0)
	for _, x := range callFrames {
		frameString = append(frameString, x.prettyString())
	}

	return mainCaller.prettyString(), frameString
}

func addErr(err error, fields []*entry) []*entry {
	return append(fields, ErrorEntry(err))
}
