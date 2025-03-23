package logging

import (
	"io"
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggerLevel int

const (
	Debug LoggerLevel = iota + 1
	Info
	Warn
	Error
	Fatal
)

type Logger struct {
	level       LoggerLevel
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
}

func New(level LoggerLevel, logToFile bool) *Logger {
	var writer io.Writer
	if logToFile {
		writer = io.MultiWriter(os.Stdout, &lumberjack.Logger{
			Filename: "./logs/log.txt",
			MaxSize:  10,
		})
	} else {
		writer = io.MultiWriter(os.Stdout)
	}

	flags := log.Ldate | log.Ltime

	return &Logger{
		level:       level,
		debugLogger: log.New(writer, "[DEBUG] ", flags),
		infoLogger:  log.New(writer, "[INFO] ", flags),
		warnLogger:  log.New(writer, "[WARN] ", flags),
		errorLogger: log.New(writer, "[ERROR] ", flags),
		fatalLogger: log.New(writer, "[FATAL] ", flags),
	}
}

func (l *Logger) SetLevel(level LoggerLevel) {
	l.level = level
}

func (l *Logger) Debug(v ...any) {
	if l.level <= Debug {
		l.debugLogger.Println(v...)
	}
}

func (l *Logger) Debugf(format string, v ...any) {
	if l.level <= Debug {
		l.debugLogger.Printf(format, v...)
	}
}

func (l *Logger) Info(v ...any) {
	if l.level <= Info {
		l.infoLogger.Println(v...)
	}
}

func (l *Logger) Infof(format string, v ...any) {
	if l.level <= Info {
		l.infoLogger.Printf(format, v...)
	}
}

func (l *Logger) Warn(v ...any) {
	if l.level <= Warn {
		l.warnLogger.Println(v...)
	}
}

func (l *Logger) Warnf(format string, v ...any) {
	if l.level <= Warn {
		l.warnLogger.Printf(format, v...)
	}
}

func (l *Logger) Error(v ...any) {
	if l.level <= Error {
		l.errorLogger.Println(v...)
	}
}

func (l *Logger) Errorf(format string, v ...any) {
	if l.level <= Error {
		l.errorLogger.Printf(format, v...)
	}
}

func (l *Logger) Fatal(v ...any) {
	if l.level <= Fatal {
		l.fatalLogger.Fatal(v...)
	}
}

func (l *Logger) Fatalf(format string, v ...any) {
	if l.level <= Fatal {
		l.fatalLogger.Fatalf(format, v...)
	}
}
