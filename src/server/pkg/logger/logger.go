package logger

import (
	"io"
	"log"
	"sync"

	"ivankprod.ru/src/server/pkg/utils"
)

// Logger interface
type ILogger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	Fatalln(v ...interface{})
}

// Logger struct
type Logger struct {
	logger *log.Logger
}

var (
	l    *log.Logger
	once sync.Once
)

// Singleton constructor
func New(w io.Writer) ILogger {
	once.Do(func() {
		l = log.New(w, "", 0)
	})

	return Logger{l}
}

// Getting logger instance
func Get() ILogger {
	return Logger{l}
}

// printf
func (l Logger) Printf(format string, v ...interface{}) {
	l.logger.SetPrefix(utils.TimeMSK_ToLocaleString() + " ")
	l.logger.Printf(format, v...)
}

// println
func (l Logger) Println(v ...interface{}) {
	l.logger.SetPrefix(utils.TimeMSK_ToLocaleString() + " ")
	l.logger.Println(v...)
}

// fatalln
func (l Logger) Fatalln(v ...interface{}) {
	l.logger.SetPrefix(utils.TimeMSK_ToLocaleString() + " ")
	l.logger.Fatalln(v...)
}
