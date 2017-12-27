package main

import (
	"fmt"
	"log"
	"os"
)

// Logging Levels
var LOG_LEVEL = map[string]int{"DEBUG": 10,
	"INFO":  20,
	"ERROR": 30,
	"FATAL": 40}

// implements leveled logging
type LevelLog struct {
	LOG_LEVEL int
	Logger    *log.Logger
}

type LevelLogger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
	LogLevelf(int, string, ...interface{})

	Debug(string)
	Info(string)
	Error(string)
	Fatal(string)
	LogLevel(int, string)
}

func (l *LevelLog) Infof(formatString string, args ...interface{}) {
	if l.LOG_LEVEL <= LOG_LEVEL["INFO"] {
		l.Logger.Printf(formatString, args...)
	}
}

func (l *LevelLog) Debugf(formatString string, args ...interface{}) {
	if l.LOG_LEVEL <= LOG_LEVEL["DEBUG"] {
		l.Logger.Printf(formatString, args...)
	}
}

func (l *LevelLog) Errorf(formatString string, args ...interface{}) {
	if l.LOG_LEVEL <= LOG_LEVEL["ERROR"] {
		l.Logger.Printf(formatString, args...)
	}
}

func (l *LevelLog) Fatalf(formatString string, args ...interface{}) {
	if l.LOG_LEVEL <= LOG_LEVEL["FATAL"] {
		l.Logger.Fatalf(formatString, args...)
	}
}

func (l *LevelLog) LogLevelf(level int, formatString string, args ...interface{}) {
	if l.LOG_LEVEL <= level {
		l.Logger.Fatalf(formatString, args...)
	}
}

func (l *LevelLog) Info(formatString string) {
	if l.LOG_LEVEL <= LOG_LEVEL["INFO"] {
		l.Logger.Print(formatString)
	}
}

func (l *LevelLog) Debug(formatString string) {
	if l.LOG_LEVEL <= LOG_LEVEL["DEBUG"] {
		l.Logger.Print(formatString)
	}
}

func (l *LevelLog) Error(formatString string) {
	if l.LOG_LEVEL <= LOG_LEVEL["ERROR"] {
		l.Logger.Print(formatString)
	}
}

func (l *LevelLog) Fatal(formatString string) {
	if l.LOG_LEVEL <= LOG_LEVEL["FATAL"] {
		l.Logger.Fatal(formatString)
	}
}

func (l *LevelLog) LogLevel(level int, formatString string) {
	if l.LOG_LEVEL <= level {
		l.Logger.Fatal(formatString)
	}
}

// initializes the logger
func initLogger(logLevelString string, logOutputPath string, prefix string) (LevelLogger, error) {
	logLevel, ok := LOG_LEVEL[logLevelString]
	if !ok {
		log.Printf("Error for log level: %+v\nValid Error Levels: \n", ok)
		for word, _ := range LOG_LEVEL {
			log.Printf("\t%s\n", word)
		}
		return nil, fmt.Errorf("Unrecognized Log Level: %s", logLevelString)
	}
	tLogger := new(LevelLog)
	tLogger.LOG_LEVEL = logLevel
	logOutput := os.Stderr
	var err error
	if logOutputPath != "stderr" {
		logOutput, err = os.OpenFile(logOutputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return nil, fmt.Errorf("Failed to open log '%s': %+v", logOutputPath, err)
		}
	}
	tLogger.Logger = log.New(logOutput, prefix, log.LstdFlags|log.LUTC)
	return tLogger, nil
}
