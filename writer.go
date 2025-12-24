package cuslog

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// Writer interface defines how log entries are written
type Writer interface {
	Write(p []byte) (n int, err error)
	Close() error
}

// ConsoleWriter writes logs to stdout/stderr
type ConsoleWriter struct {
	mu sync.Mutex
}

// Write writes log data to console
func (w *ConsoleWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return os.Stdout.Write(p)
}

// Close closes the console writer (no-op for console)
func (w *ConsoleWriter) Close() error {
	return nil
}

// FileWriter writes logs to a file
type FileWriter struct {
	FilePath string
	file     *os.File
	mu       sync.Mutex
}

// NewFileWriter creates a new file writer
func NewFileWriter(filePath string) (*FileWriter, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return &FileWriter{
		FilePath: filePath,
		file:     file,
	}, nil
}

// Write writes log data to file
func (w *FileWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file == nil {
		// Reopen the file if it was closed
		file, err := os.OpenFile(w.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return 0, fmt.Errorf("failed to reopen log file: %w", err)
		}
		w.file = file
	}

	return w.file.Write(p)
}

// Close closes the file writer
func (w *FileWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file != nil {
		err := w.file.Close()
		w.file = nil
		return err
	}
	return nil
}

// MultiWriter writes to multiple writers simultaneously
type MultiWriter struct {
	writers []Writer
	mu      sync.Mutex
}

// NewMultiWriter creates a new multi-writer
func NewMultiWriter(writers ...Writer) *MultiWriter {
	return &MultiWriter{
		writers: writers,
	}
}

// Write writes log data to all underlying writers
func (w *MultiWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Write to all writers
	var lastErr error
	for _, writer := range w.writers {
		if _, err := writer.Write(p); err != nil {
			lastErr = err
		}
	}

	// Return the length of data and the last error if any
	if lastErr != nil {
		return len(p), lastErr
	}
	return len(p), nil
}

// Close closes all underlying writers
func (w *MultiWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	var lastErr error
	for _, writer := range w.writers {
		if err := writer.Close(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// IsConsole returns true if this MultiWriter contains only ConsoleWriters
func (w *MultiWriter) IsConsole() bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	for _, writer := range w.writers {
		if _, ok := writer.(*ConsoleWriter); !ok {
			return false
		}
	}
	return true
}

// HasConsole returns true if this MultiWriter contains at least one ConsoleWriter
func (w *MultiWriter) HasConsole() bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	for _, writer := range w.writers {
		if _, ok := writer.(*ConsoleWriter); ok {
			return true
		}
	}
	return false
}

// ioWriter adapts our Writer interface to io.Writer
type ioWriter struct {
	writer Writer
}

func (w *ioWriter) Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
}

// ToIOWriter converts our Writer to io.Writer
func ToIOWriter(w Writer) io.Writer {
	return &ioWriter{writer: w}
}
