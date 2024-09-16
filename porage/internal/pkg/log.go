package pkg

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

type PorageLoggerLevel int

const (
	LogLevelDebug PorageLoggerLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// PorageLogger is the logger package for Porage.
type PorageLogger struct {
	level     PorageLoggerLevel
	output    *os.File
	withColor bool

	mutex *sync.Mutex
}

var Logger *PorageLogger = NewPorageLogger()

// =============================================================================
// Initialization functions
// =============================================================================

// NewPorageLogger creates a new PorageLogger with default settings.
func NewPorageLogger() *PorageLogger {
	l := &PorageLogger{}

	logLevelString := os.Getenv("PORAGE_LOG_LEVEL")
	if logLevelString == "" {
		logLevelString = "INFO"
	} else {
		logLevelString = strings.ToUpper(logLevelString)
	}

	l.level = l.LogLevelFromString(logLevelString)
	l.output = os.Stdout
	l.withColor = true
	l.mutex = &sync.Mutex{}

	return l
}

// LogLevelFromString converts a string to a PorageLoggerLevel. If the string is not a valid log level, it returns LogLevelInfo.
// The valid log levels are: DEBUG, INFO, WARN, ERROR. The comparison is case-insensitive.
func (l *PorageLogger) LogLevelFromString(levelString string) PorageLoggerLevel {
	levelString = strings.ToUpper(levelString)
	switch levelString {
	case "DEBUG":
		return LogLevelDebug
	case "INFO":
		return LogLevelInfo
	case "WARN":
		return LogLevelWarn
	case "ERROR":
		return LogLevelError
	}
	return LogLevelInfo
}

// SetLevel sets the level of the logger.
func (l *PorageLogger) SetLevel(level PorageLoggerLevel) {
	l.level = level
}

// SetOutput sets the output of the logger. The valid outputs are: "stdout", "stderr" and the path to a file.
// The comparison is case-sensitive.
func (l *PorageLogger) SetOutput(output string) error {
	var err error
	switch output {
	case "stdout":
		l.output = os.Stdout
	case "stderr":
		l.output = os.Stderr
	default:
		l.output, err = os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	}
	return err
}

// SetWithColor sets whether the logger should print with color.
func (l *PorageLogger) SetWithColor(withColor bool) {
	l.withColor = withColor
}

// =============================================================================
// Private functions
// =============================================================================

func (l *PorageLogger) printWithLevel(level string, content string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// Get the function position in the stack
	pc, _, _, _ := runtime.Caller(2)
	f := runtime.FuncForPC(pc)
	funcName := f.Name()
	timeString := time.Now().Format("2006-01-02 15:04:05")
	outputString := fmt.Sprintf("\r[%s %s %s] %s", timeString, level, funcName, content)
	outputString = strings.TrimSuffix(outputString, "\n")
	fmt.Fprintln(l.output, outputString)
}

func (l *PorageLogger) logLevelToString(level PorageLoggerLevel) string {
	if l.withColor {
		switch level {
		case LogLevelDebug:
			return color.HiBlueString("DEBUG")
		case LogLevelInfo:
			return color.GreenString("INFO")
		case LogLevelWarn:
			return color.YellowString("WARN")
		case LogLevelError:
			return color.RedString("ERROR")
		}
	} else {
		switch level {
		case LogLevelDebug:
			return "DEBUG"
		case LogLevelInfo:
			return "INFO"
		case LogLevelWarn:
			return "WARN"
		case LogLevelError:
			return "ERROR"
		}
	}
	return ""
}

// =============================================================================
// Print functions
// =============================================================================

func (l *PorageLogger) Debugf(format string, args ...interface{}) {
	if l.level <= LogLevelDebug {
		levelString := l.logLevelToString(LogLevelDebug)
		l.printWithLevel(levelString, fmt.Sprintf(format, args...))
	}
}

func (l *PorageLogger) Infof(format string, args ...interface{}) {
	if l.level <= LogLevelInfo {
		levelString := l.logLevelToString(LogLevelInfo)
		l.printWithLevel(levelString, fmt.Sprintf(format, args...))
	}
}

func (l *PorageLogger) Warningf(format string, args ...interface{}) {
	if l.level <= LogLevelWarn {
		levelString := l.logLevelToString(LogLevelWarn)
		l.printWithLevel(levelString, fmt.Sprintf(format, args...))
	}
}

func (l *PorageLogger) Errorf(format string, args ...interface{}) {
	if l.level <= LogLevelError {
		levelString := l.logLevelToString(LogLevelError)
		l.printWithLevel(levelString, fmt.Sprintf(format, args...))
	}
}

func (l *PorageLogger) Fatalf(format string, args ...interface{}) {
	levelString := l.logLevelToString(LogLevelError)
	l.printWithLevel(levelString, fmt.Sprintf(format, args...))
	panic("Pora panics due to a fatal error")
}
