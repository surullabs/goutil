package logging

import (
	"log"
	"os"
)

type Logger struct {
	*log.Logger
	V *log.Logger
}

var PrintLogger = log.New(os.Stdout, "", 0)
var NilLogger = log.New(NilWriter(0), "", 0)
var StdoutLogger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

type NilWriter int

func (NilWriter) Write(p []byte) (n int, err error) { return len(p), nil }

func New(l *log.Logger, v *log.Logger) *Logger { return &Logger{l, v} }
