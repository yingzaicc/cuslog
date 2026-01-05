package cuslog

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Logger represents the main logger instance
type Logger struct {
	config    Config
	formatter Formatter
	writer    Writer
	mu        sync.Mutex // for concurrent safety
}

// New creates a new logger instance with the given configuration
func New(config Config) (*Logger, error) {
	// Create writer based on configuration
	writer, err := createWriter(config.Outputs)
	if err != nil {
		return nil, fmt.Errorf("failed to create writer: %w", err)
	}

	// Determine if we should enable color based on writer type
	// Only enable color if the writer is exclusively console output
	enableColor := config.EnableColor && shouldEnableColor(writer)

	// Create formatter
	formatter := GetFormatter(config.Format, enableColor)

	return &Logger{
		config:    config,
		formatter: formatter,
		writer:    writer,
	}, nil
}

// shouldEnableColor returns true if color should be enabled for the given writer
func shouldEnableColor(writer Writer) bool {
	switch w := writer.(type) {
	case *ConsoleWriter:
		return true
	case *MultiWriter:
		return w.IsConsole()
	default:
		return false
	}
}

// createWriter creates appropriate writer(s) based on output configuration
func createWriter(outputs []OutputConfig) (Writer, error) {
	if len(outputs) == 0 {
		// Default to console output
		return &ConsoleWriter{}, nil
	}

	var writers []Writer
	for _, output := range outputs {
		switch strings.ToLower(output.Type) {
		case "console":
			writers = append(writers, &ConsoleWriter{})
		case "file":
			if output.FilePath == "" {
				return nil, fmt.Errorf("file path is required for file output")
			}
			fileWriter, err := NewFileWriter(output.FilePath)
			if err != nil {
				return nil, fmt.Errorf("failed to create file writer for %s: %w", output.FilePath, err)
			}
			writers = append(writers, fileWriter)
		default:
			return nil, fmt.Errorf("unsupported output type: %s", output.Type)
		}
	}

	if len(writers) == 1 {
		return writers[0], nil
	}
	return NewMultiWriter(writers...), nil
}

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.config.Level = level
}

// SetFormat sets the output format
func (l *Logger) SetFormat(format string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.config.Format = format
	// Re-evaluate color support based on writer type
	enableColor := l.config.EnableColor && shouldEnableColor(l.writer)
	l.formatter = GetFormatter(format, enableColor)
}

// log writes a log entry with the specified level and message
func (l *Logger) log(level Level, args ...interface{}) {
	// Check if the level is enabled
	if !level.IsEnabled(l.config.Level) {
		return
	}

	// Create log entry
	entry := &Entry{
		Time:    time.Now(),
		Level:   level,
		Message: fmt.Sprint(args...),
		Data:    make(map[string]interface{}),
	}

	// Add caller info if enabled
	if l.config.EnableCaller {
		file, line, _ := l.getCallerInfo(3) // Skip: log/logf + actual log method + getCallerInfo
		entry.File = file
		entry.Line = line
	}

	// Format the entry
	formatted, err := l.formatter.Format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting log: %v\n", err)
		return
	}

	// Write the formatted entry with mutex lock
	l.mu.Lock()
	defer l.mu.Unlock()

	_, err = l.writer.Write(formatted)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing log: %v\n", err)
	}
}

// logf writes a formatted log entry
func (l *Logger) logf(level Level, format string, args ...interface{}) {
	// Check if the level is enabled
	if !level.IsEnabled(l.config.Level) {
		return
	}

	// Create log entry
	entry := &Entry{
		Time:    time.Now(),
		Level:   level,
		Message: fmt.Sprintf(format, args...),
		Data:    make(map[string]interface{}),
	}

	// Add caller info if enabled
	if l.config.EnableCaller {
		file, line, _ := l.getCallerInfo(3) // Skip: log/logf + actual log method + getCallerInfo
		entry.File = file
		entry.Line = line
	}

	// Format the entry
	formatted, err := l.formatter.Format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting log: %v\n", err)
		return
	}

	// Write the formatted entry with mutex lock
	l.mu.Lock()
	defer l.mu.Unlock()

	_, err = l.writer.Write(formatted)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing log: %v\n", err)
	}
}

// getCallerInfo gets the caller's file and line number
func (l *Logger) getCallerInfo(skip int) (string, int, string) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", 0, ""
	}

	// Extract just the filename, not the full path
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' || file[i] == '\\' {
			file = file[i+1:]
			break
		}
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return file, line, ""
	}

	funcName := fn.Name()
	// Extract just the function name, not the full path
	for i := len(funcName) - 1; i > 0; i-- {
		if funcName[i] == '.' {
			funcName = funcName[i+1:]
			break
		}
	}

	return file, line, funcName
}

// Debug logs a debug message
func (l *Logger) Debug(args ...interface{}) {
	l.log(LevelDebug, args...)
}

// Info logs an info message
func (l *Logger) Info(args ...interface{}) {
	l.log(LevelInfo, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(args ...interface{}) {
	l.log(LevelWarn, args...)
}

// Error logs an error message
func (l *Logger) Error(args ...interface{}) {
	l.log(LevelError, args...)
}

// Fatal logs a fatal message and exits the program
func (l *Logger) Fatal(args ...interface{}) {
	l.log(LevelFatal, args...)
	os.Exit(1)
}

// Panic logs a panic message and panics
func (l *Logger) Panic(args ...interface{}) {
	msg := fmt.Sprint(args...)
	l.log(LevelPanic, msg)
	panic(msg)
}

// Debugf logs a debug message with formatting
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logf(LevelDebug, format, args...)
}

// Infof logs an info message with formatting
func (l *Logger) Infof(format string, args ...interface{}) {
	l.logf(LevelInfo, format, args...)
}

// Warnf logs a warning message with formatting
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logf(LevelWarn, format, args...)
}

// Errorf logs an error message with formatting
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logf(LevelError, format, args...)
}

// Fatalf logs a fatal message with formatting and exits the program
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logf(LevelFatal, format, args...)
	os.Exit(1)
}

// Panicf logs a panic message with formatting and panics
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.logf(LevelPanic, format, args...)
	panic(fmt.Sprintf(format, args...))
}

// Close closes the logger and releases resources
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.writer.Close()
}

// Default global logger instance
var std = NewDefaultLogger()

// NewDefaultLogger creates a logger with default configuration
func NewDefaultLogger() *Logger {
	logger, err := New(DefaultConfig())
	if err != nil {
		// Fallback to minimal logger if creation fails
		return &Logger{
			config:    DefaultConfig(),
			formatter: &TxtFormatter{EnableColor: true},
			writer:    &ConsoleWriter{},
		}
	}
	return logger
}

// SetLevel sets the minimum log level for the default logger
func SetLevel(level Level) {
	std.SetLevel(level)
}

// SetFormat sets the output format for the default logger
func SetFormat(format string) {
	std.SetFormat(format)
}

// Debug logs a debug message using the default logger
func Debug(args ...interface{}) {
	std.Debug(args...)
}

// Info logs an info message using the default logger
func Info(args ...interface{}) {
	std.Info(args...)
}

// Warn logs a warning message using the default logger
func Warn(args ...interface{}) {
	std.Warn(args...)
}

// Error logs an error message using the default logger
func Error(args ...interface{}) {
	std.Error(args...)
}

// Fatal logs a fatal message using the default logger and exits
func Fatal(args ...interface{}) {
	std.Fatal(args...)
}

// Panic logs a panic message using the default logger and panics
func Panic(args ...interface{}) {
	std.Panic(args...)
}

// Debugf logs a formatted debug message using the default logger
func Debugf(format string, args ...interface{}) {
	std.Debugf(format, args...)
}

// Infof logs a formatted info message using the default logger
func Infof(format string, args ...interface{}) {
	std.Infof(format, args...)
}

// Warnf logs a formatted warning message using the default logger
func Warnf(format string, args ...interface{}) {
	std.Warnf(format, args...)
}

// Errorf logs a formatted error message using the default logger
func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}

// Fatalf logs a formatted fatal message using the default logger and exits
func Fatalf(format string, args ...interface{}) {
	std.Fatalf(format, args...)
}

// Panicf logs a formatted panic message using the default logger and panics
func Panicf(format string, args ...interface{}) {
	std.Panicf(format, args...)
}
