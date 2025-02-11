package handler

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
)

// LogLevel represents the level of logging
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// Logger is the interface for logging
type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	SetLevel(level LogLevel)
	SetOutput(output *os.File)
	SetFormatter(formatter Formatter)
}

// Formatter is the interface for log message formatting
type Formatter interface {
	Format(level LogLevel, msg string) string
}

// TextFormatter is the custom implementation of the Formatter interface
type TextFormatter struct {
	FullTimestamp    bool
	TimestampFormat  string
	DisableTimestamp bool
	DisableColors    bool
}

// Format formats the log message
func (f *TextFormatter) Format(level LogLevel, msg string) string {
	timestamp := ""
	if !f.DisableTimestamp {
		if f.FullTimestamp {
			timestamp = time.Now().Format(f.TimestampFormat)
		} else {
			timestamp = time.Now().Format(time.RFC3339)
		}
	}

	color := f.getColor(level)
	resetColor := "\033[0m"

	if f.DisableColors {
		color = ""
		resetColor = ""
	}

	_, file, line, ok := runtime.Caller(3)
	if ok {
		file = fmt.Sprintf("%s:%d", path.Base(file), line)
	} else {
		file = "unknown"
	}

	return fmt.Sprintf("%s[%s] [SysCapture:%s]%s: %s\n", color, timestamp, file, resetColor, msg)
}

// getColor returns the color code for the given log level
func (f *TextFormatter) getColor(level LogLevel) string {
	switch level {
	case DEBUG:
		return "\033[36m" // Cyan
	case INFO:
		return "\033[34m" // Blue
	case WARN:
		return "\033[33m" // Yellow
	case ERROR:
		return "\033[31m" // Red
	default:
		return "\033[0m" // Reset
	}
}

// SysCaptureLogger is the implementation of the Logger interface
type SysCaptureLogger struct {
	mu        sync.Mutex
	level     LogLevel
	output    *os.File
	formatter Formatter
}

// NewSysCaptureLogger creates a new SysCaptureLogger
func NewSysCaptureLogger() *SysCaptureLogger {
	return &SysCaptureLogger{
		level:     INFO,
		output:    os.Stdout,
		formatter: &TextFormatter{FullTimestamp: true, TimestampFormat: time.RFC3339},
	}
}

// SetLevel sets the logging level
func (l *SysCaptureLogger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// SetOutput sets the output destination for the logger
func (l *SysCaptureLogger) SetOutput(output *os.File) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = output
}

// SetFormatter sets the formatter for the logger
func (l *SysCaptureLogger) SetFormatter(formatter Formatter) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.formatter = formatter
}

// log writes a log message with the given level
func (l *SysCaptureLogger) log(level LogLevel, msg string) {
	if level < l.level {
		return
	}
	formattedMsg := l.formatter.Format(level, msg)
	fmt.Fprint(l.output, formattedMsg)
}

// Debug logs a debug message
func (l *SysCaptureLogger) Debug(msg string) {
	l.log(DEBUG, msg)
}

// Info logs an info message
func (l *SysCaptureLogger) Info(msg string) {
	l.log(INFO, msg)
}

// Warn logs a warning message
func (l *SysCaptureLogger) Warn(msg string) {
	l.log(WARN, msg)
}

// Error logs an error message
func (l *SysCaptureLogger) Error(msg string) {
	l.log(ERROR, msg)
}
