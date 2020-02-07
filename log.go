package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// LOG_LEVEL stores a string to integer map of available levels
var LOG_LEVEL = map[string]int{"DEBUG": 10,
	"INFO":  20,
	"ERROR": 30,
	"FATAL": 40}

// LOG_STRING reverses LOG_LEVEL
var LOG_STRING = map[int]string{10: "DEBUG",
	20: "INFO",
	30: "ERROR",
	40: "FATAL"}

// LevelLog implements leveled logging
type LevelLog struct {
	LOG_LEVEL   int
	Logger      *os.File
	syncChannel chan string
	done        chan bool
}

// LevelLogger provides the interface available to consumers of the logging API
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

	// used for logging functions
	BuildLoggedFunction(int, func(...interface{}) interface{}) func(...interface{}) interface{}
	Close()
}

// Debugf uses 'formatString' to print 'args' at the 'DEBUG' level
func (l *LevelLog) Debugf(formatString string, args ...interface{}) {
	l.LogLevelf(LOG_LEVEL["DEBUG"], formatString, args...)
}

// Infof uses 'formatString' to print 'args' at the 'INFO' level
func (l *LevelLog) Infof(formatString string, args ...interface{}) {
	l.LogLevelf(LOG_LEVEL["INFO"], formatString, args...)
}

// Errorf uses 'formatString' to print 'args' at the 'ERROR' level
func (l *LevelLog) Errorf(formatString string, args ...interface{}) {
	l.LogLevelf(LOG_LEVEL["ERROR"], formatString, args...)
}

// Fatalf uses 'formatString' to print 'args' at the 'FATAL' level
func (l *LevelLog) Fatalf(formatString string, args ...interface{}) {
	l.LogLevelf(LOG_LEVEL["FATAL"], formatString, args...)
}

// LogLevelf prints all statements with a precedence greater than or equal to that set by the logger
// any call with a level greater than LOG_LEVEL["FATAL"] will result in a call to os.Exit(1)
func (l *LevelLog) LogLevelf(level int, formatString string, args ...interface{}) {
	if l.LOG_LEVEL <= level {
		message := fmt.Sprintf("%s %-10s", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("[%s]", LOG_STRING[level]))
		postFix := fmt.Sprintf(formatString, args...)
		message = fmt.Sprintf("%s : %s", message, postFix)
		if !strings.HasSuffix(formatString, "\n") {
			message = message + "\n"
		}
		l.syncChannel <- message
	}
}

// Debug logs the given string at the DEBUG level
func (l *LevelLog) Debug(message string) {
	l.LogLevelf(LOG_LEVEL["DEBUG"], "%s", message)
}

// Info logs the given string at the INFO level
func (l *LevelLog) Info(message string) {
	l.LogLevelf(LOG_LEVEL["INFO"], "%s", message)
}

// Error logs the given string at the ERROR level
func (l *LevelLog) Error(message string) {
	l.LogLevelf(LOG_LEVEL["ERROR"], "%s", message)
}

// Fatal logs the given string at the FATAL level
func (l *LevelLog) Fatal(message string) {
	l.LogLevelf(LOG_LEVEL["FATAL"], "%s", message)
}

// LogLevel logs the given string at the specified level
func (l *LevelLog) LogLevel(level int, message string) {
	l.LogLevelf(level, "%s", message)
}

// BuildLoggedFunction returns a function that executes the given function after logging it at the specified level
func (l *LevelLog) BuildLoggedFunction(level int, f func(...interface{}) interface{}) func(...interface{}) interface{} {
	return func(args ...interface{}) interface{} {
		l.LogLevelf(level, "%+v called: %+v", f, args)
		return f(args...)
	}
}

// reads from the sync channel to write information to the logger
// the reader expects only strings on the channel
func (l *LevelLog) startLogger() {
	for message := range l.syncChannel {
		fmt.Fprintf(l.Logger, "%s", message)
	}
	l.Logger.Sync()
	l.done <- true
}

// Close ensures the syncronization elements are cleaned
func (l *LevelLog) Close() {
	close(l.syncChannel)
	<-l.done
}

// initializes the logger
func InitLogger(logLevelString string, logOutputPath string, prefix string) (LevelLogger, error) {
	logLevel, ok := LOG_LEVEL[logLevelString]
	if !ok {
		log.Printf("Error for log level: %+v\nValid Error Levels: \n", ok)
		for word := range LOG_LEVEL {
			log.Printf("\t%s\n", word)
		}
		return nil, fmt.Errorf("Unrecognized Log Level: %s", logLevelString)
	}
	tLogger := new(LevelLog)
	tLogger.LOG_LEVEL = logLevel
	logOutput := os.Stderr
	var err error
	if logOutputPath != "stderr" && logOutputPath != "STDERR" {
		logOutput, err = os.OpenFile(logOutputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return nil, fmt.Errorf("Failed to open log '%s': %+v", logOutputPath, err)
		}
	}
	tLogger.Logger = logOutput
	tLogger.syncChannel = make(chan string, 1000)
	tLogger.done = make(chan bool, 1)
	go tLogger.startLogger()
	return tLogger, nil
}
