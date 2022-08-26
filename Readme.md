# glog

A simple logger for go where you can provide your custom core that'll be used to log

Features:
- Log Levels { Trace, Debug, Info, Warn, Error, Fatal }
- Stripping Levels { None, Message_Only, Minimal, Minimum_Fields }
- Custom Core
- Pretty Printing

See core implementation in `fmtcore` package where a `fmt` core is implemented
