package fmt

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rctrj/glog"
	"reflect"
)

type (
	Mode  int
	Color uint8
)

const (
	ModeDefault Mode = iota
	ModeJson
)

const (
	ColorBlack Color = iota + 30
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
)

type Option func(*Driver)

type Driver struct {
	Mode        Mode
	PrettyPrint struct {
		Enabled bool
		Colors  map[log.Level]Color
	}
}

func New(options ...Option) *Driver {
	d := &Driver{
		Mode: ModeDefault,
		PrettyPrint: struct {
			Enabled bool
			Colors  map[log.Level]Color
		}{
			Enabled: false,
			Colors:  defaultColors(),
		},
	}

	for _, option := range options {
		option(d)
	}

	return d
}

func WithColors(colors map[log.Level]Color) Option {
	return func(d *Driver) {
		d.PrettyPrint.Colors = colors
	}
}

func WithMode(mode Mode) Option {
	return func(d *Driver) {
		d.Mode = mode
	}
}

func WithPrettyPrint(enable bool) Option {
	return func(d *Driver) {
		d.PrettyPrint.Enabled = enable
	}
}

func (d Driver) Log(_ context.Context, level log.Level, fields map[string]any) {
	switch d.Mode {
	case ModeJson:
		d.logJson(level, fields)
	default:
		d.logDefault(level, fields)
	}
}

func (d Driver) logJson(level log.Level, fields map[string]any) {
	bytes, err := json.Marshal(fields)
	if err != nil {
		panic(err)
	}

	message := string(bytes)
	if d.PrettyPrint.Enabled {
		color := d.PrettyPrint.Colors[level]
		message = color.Apply(message)
	}

	fmt.Println(message)
}

func (d Driver) logDefault(level log.Level, fields map[string]any) {
	if d.PrettyPrint.Enabled {
		fields = d.prettifyMap(level, fields)
	}

	fmt.Println(fields)
}

func (d Driver) prettifyMap(level log.Level, in map[string]any) map[string]any {
	if !d.PrettyPrint.Enabled {
		return in
	}

	for key, value := range in {
		color := d.PrettyPrint.Colors[level]
		prettifiedKey := color.Apply(key)
		prettifiedValue := value

		if reflect.ValueOf(value).Kind() != reflect.Struct {
			prettifiedValue = color.Apply(fmt.Sprintf("%v", value))
		}

		delete(in, key)
		in[prettifiedKey] = prettifiedValue
	}

	return in
}

func defaultColors() map[log.Level]Color {
	return map[log.Level]Color{
		log.LevelTrace: ColorCyan,
		log.LevelDebug: ColorMagenta,
		log.LevelInfo:  ColorBlue,
		log.LevelWarn:  ColorYellow,
		log.LevelError: ColorRed,
	}
}

func (c Color) Apply(in string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), in)
}
