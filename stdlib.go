package cuslog

import (
	"fmt"
	"io"
)

// StdLogger is a compatibility layer for Go's standard library log package
type StdLogger struct {
	logger *Logger
}

// NewStdLogger creates a new standard library compatible logger
func NewStdLogger(logger *Logger) *StdLogger {
	return &StdLogger{logger: logger}
}

// Print implements the standard library log.Print
func (l *StdLogger) Print(v ...interface{}) {
	l.logger.Info(v...)
}

// Printf implements the standard library log.Printf
func (l *StdLogger) Printf(format string, v ...interface{}) {
	l.logger.Infof(format, v...)
}

// Println implements the standard library log.Println
func (l *StdLogger) Println(v ...interface{}) {
	l.logger.Info(fmt.Sprint(v...))
}

// Fatal implements the standard library log.Fatal
func (l *StdLogger) Fatal(v ...interface{}) {
	l.logger.Fatal(v...)
}

// Fatalf implements the standard library log.Fatalf
func (l *StdLogger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatalf(format, v...)
}

// Fatalln implements the standard library log.Fatalln
func (l *StdLogger) Fatalln(v ...interface{}) {
	l.logger.Fatal(fmt.Sprint(v...))
}

// Panic implements the standard library log.Panic
func (l *StdLogger) Panic(v ...interface{}) {
	l.logger.Panic(v...)
}

// Panicf implements the standard library log.Panicf
func (l *StdLogger) Panicf(format string, v ...interface{}) {
	l.logger.Panicf(format, v...)
}

// Panicln implements the standard library log.Panicln
func (l *StdLogger) Panicln(v ...interface{}) {
	l.logger.Panic(fmt.Sprint(v...))
}

// Flags implements the standard library log.Flags (no-op for compatibility)
func (l *StdLogger) Flags() int {
	return 0
}

// SetFlags implements the standard library log.SetFlags (no-op for compatibility)
func (l *StdLogger) SetFlags(flag int) {}

// Prefix implements the standard library log.Prefix (no-op for compatibility)
func (l *StdLogger) Prefix() string {
	return ""
}

// SetPrefix implements the standard library log.SetPrefix (no-op for compatibility)
func (l *StdLogger) SetPrefix(prefix string) {}

// Output implements the standard library log.Output (compatibility shim)
func (l *StdLogger) Output(calldepth int, s string) error {
	l.logger.Info(s)
	return nil
}

// Writer implements the standard library log.Writer
func (l *StdLogger) Writer() io.Writer {
	return ToIOWriter(l.logger.writer)
}

// Standard library compatibility functions using the default logger
// Note: Fatal, Panic, and their formatted versions are already defined in logger.go
// as global functions, so we only provide the Print variants here

// Print writes to the standard logger using Info level
func Print(v ...interface{}) {
	std.Info(v...)
}

// Printf writes a formatted message to the standard logger using Info level
func Printf(format string, v ...interface{}) {
	std.Infof(format, v...)
}

// Println writes a line to the standard logger using Info level
func Println(v ...interface{}) {
	std.Info(fmt.Sprint(v...))
}

// SetOutput sets the output destination for the standard logger
// This is a compatibility function that recreates the logger with new output
func SetOutput(w io.Writer) {
	// Note: This is a simplified implementation
	// For full compatibility, you would need to wrap io.Writer to our Writer interface
	// For now, we just update the writer if it's a compatible type
	if wrapper, ok := w.(*ioWriter); ok {
		std.writer = wrapper.writer
	}
}

// Flags returns the flags for the standard logger (compatibility, always 0)
func Flags() int {
	return 0
}

// SetFlags sets the flags for the standard logger (no-op for compatibility)
func SetFlags(flag int) {
	// No-op for compatibility with standard library
}

// Prefix returns the prefix for the standard logger (compatibility, always empty)
func Prefix() string {
	return ""
}

// SetPrefix sets the prefix for the standard logger (no-op for compatibility)
func SetPrefix(prefix string) {
	// No-op for compatibility with standard library
}

// Output writes a message with the specified call depth
func Output(calldepth int, s string) error {
	std.Info(s)
	return nil
}
