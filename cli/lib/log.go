package lib

import (
	"io"
	"log"
)

// Log is a global logger. It may be called from anywhere in the CLI. It has a mutex that allows
// it to be called from multiple goroutines simulaneously
var Log logger

type logger struct {
	*log.Logger
	debug bool
}

// InitLog initializes the global logger
func InitLog(out io.Writer) {
	Log = logger{
		Logger: log.New(out, "", log.LstdFlags),
	}
	Log.SetOutput(out)
}

func (l logger) SetDebug(b bool) {
	l.debug = b
}

func (l logger) Debugln(v ...interface{}) {
	if l.debug {
		l.Println(v...)
	}
}

func (l logger) Debugf(format string, v ...interface{}) {
	if l.debug {
		l.Printf(format, v...)
	}
}
