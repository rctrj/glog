package fmtcore

import log "github.com/rctrj/glog"

type Config struct {
	log.Config
}

func NewConfig(levelStr, stripLevelStr string, prettyPrint, displayPrefix bool) Config {
	return Config{
		log.NewConfig(levelStr, stripLevelStr, prettyPrint, displayPrefix),
	}
}

func New(conf Config) *log.Logger {
	return log.New(conf.Config, newCore(conf.PrettyPrint))
}
