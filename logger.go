package log

import (
	"context"
	"github.com/rctrj/glog/slices"
)

type ctxKey int

const ctxKeyStackTrace ctxKey = iota

func withStackTrace(ctx context.Context, stackTrace []callFrame) context.Context {
	return context.WithValue(ctx, ctxKeyStackTrace, stackTrace)
}

func stackTraceFrom(ctx context.Context) []callFrame {
	return ctx.Value(ctxKeyStackTrace).([]callFrame)
}

type Level int

const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
)

type (
	Config struct {
		Level  Level
		Fields struct {
			Keys struct {
				Level      string
				Message    string
				StackTrace string
			}
			ToPrint struct {
				StackTrace     bool
				FromContext    bool
				MessageContext bool
			}
		}
		StackTraceLevels int
	}
	Contextualizer interface {
		LogEntries(ctx context.Context) map[string]any
	}
	Driver interface {
		Log(ctx context.Context, level Level, fields map[string]any)
	}
	Logger struct {
		config          Config
		contextualizers []Contextualizer
		driver          Driver
	}
	Option func(*Logger)
)

func New(driver Driver, opts ...Option) *Logger {
	config := defaultConfig()

	out := &Logger{
		config: config,
		driver: driver,
	}

	for _, opt := range opts {
		opt(out)
	}

	return out
}

func WithContextualizer(contextualizers ...Contextualizer) Option {
	return func(l *Logger) {
		l.contextualizers = append(l.contextualizers, contextualizers...)
	}
}

func WithFieldStackTrace(enable bool) Option {
	return func(l *Logger) {
		l.config.Fields.ToPrint.StackTrace = enable
	}
}

func WithFieldFromContext(enable bool) Option {
	return func(l *Logger) {
		l.config.Fields.ToPrint.FromContext = enable
	}
}

func WithFieldKeyLevel(key string) Option {
	return func(l *Logger) {
		l.config.Fields.Keys.Level = key
	}
}

func WithFieldKeyMessage(key string) Option {
	return func(l *Logger) {
		l.config.Fields.Keys.Message = key
	}
}

func WithFieldKeyStackTrace(key string) Option {
	return func(l *Logger) {
		l.config.Fields.Keys.StackTrace = key
	}
}

func WithFieldMessageContext(enable bool) Option {
	return func(l *Logger) {
		l.config.Fields.ToPrint.MessageContext = enable
	}
}

func WithLevel(level Level) Option {
	return func(l *Logger) {
		l.config.Level = level
	}
}

func WithStackTraceLevels(levels int) Option {
	return func(l *Logger) {
		l.config.StackTraceLevels = levels
	}
}

type Entry struct {
	key          string
	value        any
	isContextual bool
}

func KV(key string, value any) Entry {
	return Entry{
		key:   key,
		value: value,
	}
}

func KVContextual(key string, value any) Entry {
	return Entry{
		key:          key,
		value:        value,
		isContextual: true,
	}
}

func Error(err error) Entry {
	return KV("error", err)
}

func (l Logger) Trace(ctx context.Context, message string, entries ...Entry) {
	l.log(ctx, LevelTrace, message, entries...)
}

func (l Logger) Debug(ctx context.Context, message string, entries ...Entry) {
	l.log(ctx, LevelDebug, message, entries...)
}

func (l Logger) Info(ctx context.Context, message string, entries ...Entry) {
	l.log(ctx, LevelInfo, message, entries...)
}

func (l Logger) Warn(ctx context.Context, message string, entries ...Entry) {
	l.log(ctx, LevelWarn, message, entries...)
}

func (l Logger) Error(ctx context.Context, message string, entries ...Entry) {
	l.log(ctx, LevelError, message, entries...)
}

func (l Logger) log(ctx context.Context, level Level, message string, entries ...Entry) {
	if level < l.config.Level {
		// do not log as level is lower than the configured level
		return
	}

	stackTrace := callers(2, 4)
	ctx = withStackTrace(ctx, stackTrace)

	entriesToPrint := l.entriesToPrint(ctx, level, message, entries)
	l.driver.Log(ctx, level, entriesToPrint)
}

func (l Logger) entriesToPrint(ctx context.Context, level Level, message string, in []Entry) map[string]any {
	out := make(map[string]any, len(in))

	toPrint := l.config.Fields.ToPrint
	keys := l.config.Fields.Keys

	addEntry := func(check bool, key string, value any) {
		if check {
			out[key] = value
		}
	}

	for _, entry := range in {
		if entry.isContextual {
			addEntry(toPrint.FromContext, entry.key, entry.value)
		} else {
			addEntry(true, entry.key, entry.value)
		}
	}

	if toPrint.FromContext {
		for _, contextualizer := range l.contextualizers {
			for key, value := range contextualizer.LogEntries(ctx) {
				// true as that's already validated before entering the for loop
				addEntry(true, key, value)
			}
		}
	}

	if toPrint.StackTrace {
		stackTrace := stackTraceFrom(ctx)
		stackTraceReadable := slices.Map(stackTrace, callFrame.prettyString)
		addEntry(true, keys.StackTrace, stackTraceReadable)
	}

	addEntry(true, keys.Message, message)
	addEntry(true, keys.Level, level.String())
	return out
}

func defaultConfig() Config {
	return Config{
		Level: LevelInfo,
		Fields: struct {
			Keys struct {
				Level      string
				Message    string
				StackTrace string
			}
			ToPrint struct {
				StackTrace     bool
				FromContext    bool
				MessageContext bool
			}
		}{
			Keys: struct {
				Level      string
				Message    string
				StackTrace string
			}{
				Level:      "level",
				Message:    "message",
				StackTrace: "stack_trace",
			},
			ToPrint: struct {
				StackTrace     bool
				FromContext    bool
				MessageContext bool
			}{
				StackTrace:     true,
				FromContext:    true,
				MessageContext: true,
			},
		},
		StackTraceLevels: 4,
	}
}

func (l Level) String() string {
	switch l {
	case LevelTrace:
		return "trace"
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	default:
		return "unknown"
	}
}
