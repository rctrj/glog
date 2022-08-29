package fmtcore

import log "github.com/rctrj/glog"

type Config struct {
	log.Config
}

func NewConfig(levelStr, stripLevelStr string, prettyPrint, displayPrefix bool) Config {
	return NewConfigFromLogConfig(log.NewConfig(levelStr, stripLevelStr, prettyPrint, displayPrefix))
}

func NewConfigFromLogConfig(in log.Config) Config {
	return Config{log.Config{
		Level:         in.Level,
		StripLevel:    in.StripLevel,
		PrettyPrint:   in.PrettyPrint,
		DisplayPrefix: in.DisplayPrefix,
	}}
}

func New(conf Config) *log.Logger {
	return log.New(conf.Config, newCore(conf.PrettyPrint))
}
