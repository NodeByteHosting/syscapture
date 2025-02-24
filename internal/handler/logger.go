package handler

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Level represents the logging level
type Level uint32

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

// Convert the Level to a string representation
func (l Level) String() string {
	switch l {
	case TraceLevel:
		return "TRACE"
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	case PanicLevel:
		return "PANIC"
	default:
		return fmt.Sprintf("LEVEL(%d)", l)
	}
}

// ParseLevel converts a string level to Level type
func ParseLevel(level string) (Level, error) {
	switch strings.ToLower(level) {
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	case "trace":
		return TraceLevel, nil
	}
	return InfoLevel, fmt.Errorf("invalid log level: %s", level)
}

// Logger interface defines the logging methods
type Logger interface {
	Trace(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
	Panic(msg string, args ...interface{})
	SetLevel(level Level)
	SetOutput(output io.Writer)
	SetFormatter(formatter Formatter)
	WithFields(fields Fields) Logger
}

// Fields represents structured log fields
type Fields map[string]interface{}

// Formatter interface for formatting log messages
type Formatter interface {
	Format(entry *Entry) string
}

// Entry represents a log entry
type Entry struct {
	Level      Level
	Time       time.Time
	Message    string
	Fields     Fields
	Caller     *runtime.Frame
	CallerInfo string
}

// TextFormatter implements Formatter interface
type TextFormatter struct {
	FullTimestamp    bool
	TimestampFormat  string
	DisableTimestamp bool
	DisableColors    bool
	DisableCaller    bool
}

// Format formats the log entry
func (f *TextFormatter) Format(entry *Entry) string {
	var sb strings.Builder

	// Add timestamp
	if !f.DisableTimestamp {
		timestamp := entry.Time.Format(f.TimestampFormat)
		sb.WriteString(fmt.Sprintf("[%s] ", timestamp))
	}

	// Add level with color
	if !f.DisableColors {
		sb.WriteString(f.colorize(entry.Level))
	}
	sb.WriteString(fmt.Sprintf("[%-5s]", entry.Level))
	if !f.DisableColors {
		sb.WriteString("\033[0m")
	}

	// Add caller info
	if !f.DisableCaller && entry.CallerInfo != "" {
		sb.WriteString(fmt.Sprintf(" [%s]", entry.CallerInfo))
	}

	// Add message
	sb.WriteString(fmt.Sprintf(" %s", entry.Message))

	// Add fields if any
	if len(entry.Fields) > 0 {
		sb.WriteString(" {")
		first := true
		for k, v := range entry.Fields {
			if !first {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("%s=%v", k, v))
			first = false
		}
		sb.WriteString("}")
	}

	sb.WriteString("\n")
	return sb.String()
}

// colorize returns ANSI color codes for different log levels
func (f *TextFormatter) colorize(level Level) string {
	switch level {
	case TraceLevel:
		return "\033[37m" // White
	case DebugLevel:
		return "\033[36m" // Cyan
	case InfoLevel:
		return "\033[34m" // Blue
	case WarnLevel:
		return "\033[33m" // Yellow
	case ErrorLevel:
		return "\033[31m" // Red
	case FatalLevel:
		return "\033[35m" // Magenta
	case PanicLevel:
		return "\033[41m" // Red background
	default:
		return "\033[0m" // Reset
	}
}

// SysCaptureLogger implements the Logger interface
type SysCaptureLogger struct {
	mu        sync.Mutex
	level     Level
	output    io.Writer
	formatter Formatter
	fields    Fields
}

// NewSysCaptureLogger creates a new logger instance
func NewSysCaptureLogger() *SysCaptureLogger {
	return &SysCaptureLogger{
		level:  InfoLevel,
		output: os.Stdout,
		formatter: &TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		},
		fields: make(Fields),
	}
}

func (l *SysCaptureLogger) log(level Level, msg string, args ...interface{}) {
	if level > l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Get caller information
	pc, file, line, ok := runtime.Caller(2)
	var callerInfo string
	if ok {
		callerInfo = fmt.Sprintf("%s:%d", path.Base(file), line)
	}

	// Create log entry
	entry := &Entry{
		Level:      level,
		Time:       time.Now(),
		Message:    fmt.Sprintf(msg, args...),
		Fields:     l.fields,
		CallerInfo: callerInfo,
	}

	if ok {
		entry.Caller = &runtime.Frame{
			PC:   pc,
			File: file,
			Line: line,
		}
	}

	// Format and write the log entry
	formatted := l.formatter.Format(entry)
	fmt.Fprint(l.output, formatted)

	// Handle panic and fatal levels
	if level <= FatalLevel {
		os.Exit(1)
	}
	if level <= PanicLevel {
		panic(entry.Message)
	}
}

func (l *SysCaptureLogger) Trace(msg string, args ...interface{}) {
	l.log(TraceLevel, msg, args...)
}

func (l *SysCaptureLogger) Debug(msg string, args ...interface{}) {
	l.log(DebugLevel, msg, args...)
}

func (l *SysCaptureLogger) Info(msg string, args ...interface{}) {
	l.log(InfoLevel, msg, args...)
}

func (l *SysCaptureLogger) Warn(msg string, args ...interface{}) {
	l.log(WarnLevel, msg, args...)
}

func (l *SysCaptureLogger) Error(msg string, args ...interface{}) {
	l.log(ErrorLevel, msg, args...)
}

func (l *SysCaptureLogger) Fatal(msg string, args ...interface{}) {
	l.log(FatalLevel, msg, args...)
}

func (l *SysCaptureLogger) Panic(msg string, args ...interface{}) {
	l.log(PanicLevel, msg, args...)
}

func (l *SysCaptureLogger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

func (l *SysCaptureLogger) SetOutput(output io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = output
}

func (l *SysCaptureLogger) SetFormatter(formatter Formatter) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.formatter = formatter
}

func (l *SysCaptureLogger) WithFields(fields Fields) Logger {
	newLogger := &SysCaptureLogger{
		level:     l.level,
		output:    l.output,
		formatter: l.formatter,
		fields:    make(Fields, len(l.fields)+len(fields)),
	}

	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	// Add new fields
	for k, v := range fields {
		newLogger.fields[k] = v
	}

	return newLogger
}
