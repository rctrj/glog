package log

import "fmt"

type color uint8
type Level string
type StripLevel string

const (
	LevelTrace Level = "trace"
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
	LevelFatal Level = "fatal"
)

func (l Level) String() string {
	return string(l)
}

func (l Level) IsTrace() bool { return l == LevelTrace }
func (l Level) IsDebug() bool { return l == LevelDebug }
func (l Level) IsInfo() bool  { return l == LevelInfo }
func (l Level) IsWarn() bool  { return l == LevelWarn }
func (l Level) IsError() bool { return l == LevelError }
func (l Level) IsFatal() bool { return l == LevelFatal }

const (
	StripLevelNone          StripLevel = "none"           //no stripping is done
	StripLevelMinimal       StripLevel = "minimal"        //keys like caller stack, prefix, etc are pruned
	StripLevelMessageOnly   StripLevel = "message_only"   //only message is printed. Everything else is pruned
	StripLevelMinimumFields StripLevel = "minimum_fields" //only message and fields specifically marked to be minimal data are printed. Rest are pruned
)

func (s StripLevel) isNone() bool          { return s == StripLevelNone }
func (s StripLevel) isMinimal() bool       { return s == StripLevelMinimal }
func (s StripLevel) isMessageOnly() bool   { return s == StripLevelMessageOnly }
func (s StripLevel) isMinimumFields() bool { return s == StripLevelMinimumFields }

const (
	Black color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

var levelScore = map[Level]int{
	LevelFatal: 60,
	LevelError: 50,
	LevelWarn:  40,
	LevelInfo:  30,
	LevelDebug: 20,
	LevelTrace: 10,
}

var levelColor = map[Level]color{
	LevelTrace: Cyan,
	LevelDebug: Magenta,
	LevelInfo:  Blue,
	LevelWarn:  Yellow,
	LevelError: Red,
	LevelFatal: Red,
}

func (l Level) greaterThan(c Level) bool      { return l.score() > c.score() }
func (l Level) greaterThanEqual(c Level) bool { return l.score() >= c.score() }

func (l Level) score() int {
	if score, ok := levelScore[l]; ok {
		return score
	} else {
		return 10
	}
}

func (l Level) color() color {
	if c, ok := levelColor[l]; ok {
		return c
	} else {
		return Cyan
	}
}

func (c color) Add(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), s)
}
