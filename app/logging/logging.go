package logging

import (
	"log"
	"os"
)

type Logger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
}

func New() Logger {
	flags := log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix
	return Logger{
		debugLogger: log.New(os.Stdout, "[DEBUG]", flags),
		infoLogger:  log.New(os.Stdout, "[INFO]", flags),
		warnLogger:  log.New(os.Stdout, "[WARN]", flags),
		errorLogger: log.New(os.Stdout, "[ERROR]", flags),
	}
}

func (l *Logger) Debug(v ...any) {
	l.debugLogger.Println(v...)
}

func (l *Logger) Info(v ...any) {
	l.infoLogger.Println(v...)
}

func (l *Logger) Warn(v ...any) {
	l.warnLogger.Println(v...)
}

func (l *Logger) Error(v ...any) {
	l.errorLogger.Println(v...)
}
