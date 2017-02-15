package lib

import (
	"io"
	"log"
)

type loglevel uint

const (
	Info loglevel = iota
	Warn
	Debug
	Dev
)

// Log is a global logger. It may be called from anywhere in the CLI. It has a mutex that allows
// it to be called from multiple goroutines simulaneously
var Log *logger

type logger struct {
	*log.Logger
	//debug bool
	level loglevel
}

// InitLog initializes the global logger
func InitLog(out io.Writer) {
	Log = new(logger)
	Log.Logger = log.New(out, "", log.LstdFlags)
	Log.SetOutput(out)
}

func (l *logger) SetLevel(level uint) {
	l.level = loglevel(level)
}

func (l *logger) Devln(v ...interface{}) {
	if l.level >= Dev {
		l.Println(v...)
	}
}

func (l *logger) Devf(format string, v ...interface{}) {
	if l.level >= Dev {
		l.Printf(format, v...)
		l.Println()
	}
}

func (l *logger) Debugln(v ...interface{}) {
	if l.level >= Debug {
		l.Println(v...)
	}
}

func (l *logger) Debugf(format string, v ...interface{}) {
	if l.level >= Debug {
		l.Printf(format, v...)
		l.Println()
	}
}

func (l *logger) Warnln(v ...interface{}) {
	if l.level >= Warn {
		l.Println(v...)
	}
}

func (l *logger) Warnf(format string, v ...interface{}) {
	if l.level >= Warn {
		l.Printf(format, v...)
	}
}
