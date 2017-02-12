package lib

import (
	"io"
	"log"
)

type loglevel uint

const (
	info loglevel = iota
	warn
	debug
)

// Log is a global logger. It may be called from anywhere in the CLI. It has a mutex that allows
// it to be called from multiple goroutines simulaneously
var Log *logger

type logger struct {
	*log.Logger
	debug bool
	level loglevel
}

// InitLog initializes the global logger
func InitLog(out io.Writer) {
	Log = new(logger)
	Log.Logger = log.New(out, "", log.LstdFlags)
	Log.SetOutput(out)
}

func (l *logger) SetLevel(level loglevel) {
	l.level = level
}

func (l *logger) SetDebug(b bool) {
	l.debug = b
}

func (l *logger) Debugln(v ...interface{}) {
	if l.debug {
		l.Println(v...)
	}
}

func (l *logger) Debugf(format string, v ...interface{}) {
	if l.debug {
		l.Printf(format, v...)
		l.Println()
	}
}

func (l *logger) Warnln(v ...interface{}) {
	if l.level >= warn {
		l.Println(v...)
	}
}

func (l *logger) Warnf(format string, v ...interface{}) {
	if l.level >= warn {
		l.Printf(format, v...)
	}
}
