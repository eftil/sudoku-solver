package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

// LogLevel defines the severity of log messages
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// String returns the string representation of a LogLevel
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger handles all logging operations with structured output
type Logger struct {
	mu         sync.Mutex
	level      LogLevel
	output     io.Writer
	prefix     string
	showTime   bool
	showCaller bool
}

// Global logger instance
var globalLogger *Logger
var once sync.Once

// init initializes the global logger with default settings
func init() {
	globalLogger = &Logger{
		level:      INFO,
		output:     os.Stdout,
		prefix:     "",
		showTime:   true,
		showCaller: false,
	}
}

// GetLogger returns the global logger instance
func GetLogger() *Logger {
	return globalLogger
}

// SetLevel sets the minimum log level for the logger
func SetLevel(level LogLevel) {
	globalLogger.mu.Lock()
	defer globalLogger.mu.Unlock()
	globalLogger.level = level
}

// SetOutput sets the output destination for the logger
func SetOutput(w io.Writer) {
	globalLogger.mu.Lock()
	defer globalLogger.mu.Unlock()
	globalLogger.output = w
}

// SetPrefix sets a prefix for all log messages
func SetPrefix(prefix string) {
	globalLogger.mu.Lock()
	defer globalLogger.mu.Unlock()
	globalLogger.prefix = prefix
}

// formatMessage formats a log message with timestamp and level
func (l *Logger) formatMessage(level LogLevel, format string, args ...interface{}) string {
	l.mu.Lock()
	defer l.mu.Unlock()

	var msg string
	if l.showTime {
		msg = fmt.Sprintf("[%s] ", time.Now().Format("2006-01-02 15:04:05"))
	}

	msg += fmt.Sprintf("[%s] ", level.String())

	if l.prefix != "" {
		msg += fmt.Sprintf("[%s] ", l.prefix)
	}

	msg += fmt.Sprintf(format, args...)

	return msg
}

// log is the internal logging method
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	msg := l.formatMessage(level, format, args...)

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.output != nil {
		fmt.Fprintln(l.output, msg)
	}
}

// Debug logs a debug message
func Debug(format string, args ...interface{}) {
	globalLogger.log(DEBUG, format, args...)
}

// Info logs an info message
func Info(format string, args ...interface{}) {
	globalLogger.log(INFO, format, args...)
}

// Warn logs a warning message
func Warn(format string, args ...interface{}) {
	globalLogger.log(WARN, format, args...)
}

// Error logs an error message
func Error(format string, args ...interface{}) {
	globalLogger.log(ERROR, format, args...)
}

// DebugCell logs cell-specific debug information
func DebugCell(row, col int, format string, args ...interface{}) {
	msg := fmt.Sprintf("[Cell R%dC%d] %s", row+1, col+1, fmt.Sprintf(format, args...))
	globalLogger.log(DEBUG, "%s", msg)
}

// InfoCell logs cell-specific info
func InfoCell(row, col int, format string, args ...interface{}) {
	msg := fmt.Sprintf("[Cell R%dC%d] %s", row+1, col+1, fmt.Sprintf(format, args...))
	globalLogger.log(INFO, "%s", msg)
}

// DebugConstraint logs constraint-specific debug information
func DebugConstraint(constraintName string, format string, args ...interface{}) {
	msg := fmt.Sprintf("[%s] %s", constraintName, fmt.Sprintf(format, args...))
	globalLogger.log(DEBUG, "%s", msg)
}

// InfoConstraint logs constraint-specific info
func InfoConstraint(constraintName string, format string, args ...interface{}) {
	msg := fmt.Sprintf("[%s] %s", constraintName, fmt.Sprintf(format, args...))
	globalLogger.log(INFO, "%s", msg)
}

// SolvingStep logs a solving technique step
func SolvingStep(technique string, format string, args ...interface{}) {
	msg := fmt.Sprintf("[SOLVING: %s] %s", technique, fmt.Sprintf(format, args...))
	globalLogger.log(INFO, "%s", msg)
}

// CandidateElimination logs when candidates are eliminated
func CandidateElimination(row, col, candidate int, reason string) {
	msg := fmt.Sprintf("[Cell R%dC%d] Eliminated candidate %d - Reason: %s",
		row+1, col+1, candidate, reason)
	globalLogger.log(DEBUG, "%s", msg)
}

// CellSolved logs when a cell is solved
func CellSolved(row, col, value int, reason string) {
	msg := fmt.Sprintf("[Cell R%dC%d] Solved with value %d - Reason: %s",
		row+1, col+1, value, reason)
	globalLogger.log(INFO, "%s", msg)
}

// Fatal logs a fatal error and exits the program
func Fatal(format string, args ...interface{}) {
	msg := globalLogger.formatMessage(ERROR, format, args...)
	log.Fatal(msg)
}

// NewLogger creates a new logger instance with custom settings
func NewLogger(level LogLevel, output io.Writer, prefix string) *Logger {
	return &Logger{
		level:      level,
		output:     output,
		prefix:     prefix,
		showTime:   true,
		showCaller: false,
	}
}
