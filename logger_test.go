package cuslog

import (
	"bytes"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestLevelString tests the Level.String() method
func TestLevelString(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{LevelDebug, "DEBUG"},
		{LevelInfo, "INFO"},
		{LevelWarn, "WARN"},
		{LevelError, "ERROR"},
		{LevelFatal, "FATAL"},
		{LevelPanic, "PANIC"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.level.String(); got != tt.expected {
				t.Errorf("Level.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestLevelIsEnabled tests the Level.IsEnabled() method
func TestLevelIsEnabled(t *testing.T) {
	tests := []struct {
		name     string
		level    Level
		minLevel Level
		expected bool
	}{
		{"Debug enabled when min is Debug", LevelDebug, LevelDebug, true},
		{"Debug disabled when min is Info", LevelDebug, LevelInfo, false},
		{"Info enabled when min is Debug", LevelInfo, LevelDebug, true},
		{"Info enabled when min is Info", LevelInfo, LevelInfo, true},
		{"Error enabled when min is Info", LevelError, LevelInfo, true},
		{"Warn disabled when min is Error", LevelWarn, LevelError, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.level.IsEnabled(tt.minLevel); got != tt.expected {
				t.Errorf("Level.IsEnabled() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestDefaultConfig tests the default configuration
func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Level != LevelInfo {
		t.Errorf("DefaultConfig().Level = %v, want %v", config.Level, LevelInfo)
	}

	if config.Format != "txt" {
		t.Errorf("DefaultConfig().Format = %v, want %v", config.Format, "txt")
	}

	if !config.EnableCaller {
		t.Errorf("DefaultConfig().EnableCaller = %v, want true", config.EnableCaller)
	}

	if !config.EnableColor {
		t.Errorf("DefaultConfig().EnableColor = %v, want true", config.EnableColor)
	}

	if len(config.Outputs) != 1 {
		t.Errorf("len(DefaultConfig().Outputs) = %v, want 1", len(config.Outputs))
	}

	if config.Outputs[0].Type != "console" {
		t.Errorf("DefaultConfig().Outputs[0].Type = %v, want %v", config.Outputs[0].Type, "console")
	}
}

// TestNewLogger tests creating a new logger
func TestNewLogger(t *testing.T) {
	config := DefaultConfig()
	logger, err := New(config)

	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if logger == nil {
		t.Fatal("New() returned nil logger")
	}

	if logger.config.Level != config.Level {
		t.Errorf("logger.config.Level = %v, want %v", logger.config.Level, config.Level)
	}
}

// TestTxtFormatter tests the text formatter
func TestTxtFormatter(t *testing.T) {
	formatter := &TxtFormatter{EnableColor: false}

	entry := &Entry{
		Time:    time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
		Level:   LevelInfo,
		Message: "test message",
		File:    "test.go",
		Line:    42,
	}

	output, err := formatter.Format(entry)
	if err != nil {
		t.Fatalf("Format() error = %v", err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "2023-10-01 12:00:00") {
		t.Errorf("Format() output does not contain timestamp")
	}

	if !strings.Contains(outputStr, "[INFO]") {
		t.Errorf("Format() output does not contain level")
	}

	if !strings.Contains(outputStr, "test.go:42") {
		t.Errorf("Format() output does not contain caller info")
	}

	if !strings.Contains(outputStr, "test message") {
		t.Errorf("Format() output does not contain message")
	}
}

// TestJsonFormatter tests the JSON formatter
func TestJsonFormatter(t *testing.T) {
	formatter := &JsonFormatter{}

	entry := &Entry{
		Time:    time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
		Level:   LevelInfo,
		Message: "test message",
		File:    "test.go",
		Line:    42,
	}

	output, err := formatter.Format(entry)
	if err != nil {
		t.Fatalf("Format() error = %v", err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "\"time\"") {
		t.Errorf("Format() output does not contain time field")
	}

	if !strings.Contains(outputStr, "\"level\"") {
		t.Errorf("Format() output does not contain level field")
	}

	if !strings.Contains(outputStr, "\"msg\"") {
		t.Errorf("Format() output does not contain message field")
	}

	if !strings.Contains(outputStr, "test.go") {
		t.Errorf("Format() output does not contain file field")
	}
}

// TestConsoleWriter tests the console writer
func TestConsoleWriter(t *testing.T) {
	writer := &ConsoleWriter{}

	// Redirect stdout to capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	testData := []byte("test log message\n")
	_, err := writer.Write(testData)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "test log message") {
		t.Errorf("ConsoleWriter.Write() output = %v, want to contain 'test log message'", output)
	}

	// Test Close
	if err := writer.Close(); err != nil {
		t.Errorf("Close() error = %v, want nil", err)
	}
}

// TestFileWriter tests the file writer
func TestFileWriter(t *testing.T) {
	tmpFile := "test_log.txt"
	defer os.Remove(tmpFile)

	writer, err := NewFileWriter(tmpFile)
	if err != nil {
		t.Fatalf("NewFileWriter() error = %v", err)
	}

	testData := []byte("test log message\n")
	_, err = writer.Write(testData)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Close the writer
	if err := writer.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Read the file to verify content
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if !strings.Contains(string(content), "test log message") {
		t.Errorf("FileWriter.Write() content = %v, want to contain 'test log message'", string(content))
	}
}

// TestLoggerLevels tests logging at different levels
func TestLoggerLevels(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer

	config := Config{
		Level:        LevelInfo,
		Format:       "txt",
		EnableCaller: false,
		EnableColor:  false,
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Replace the writer with a buffer writer
	logger.writer = &bufferWriter{buf: &buf}

	// Debug should not be logged (below minimum level)
	logger.Debug("debug message")
	if buf.Len() > 0 {
		t.Error("Debug() should not log when Level is Info")
	}

	// Info should be logged
	logger.Info("info message")
	if !strings.Contains(buf.String(), "info message") {
		t.Error("Info() should log 'info message'")
	}

	buf.Reset()

	// Warn should be logged
	logger.Warn("warn message")
	if !strings.Contains(buf.String(), "warn message") {
		t.Error("Warn() should log 'warn message'")
	}

	buf.Reset()

	// Error should be logged
	logger.Error("error message")
	if !strings.Contains(buf.String(), "error message") {
		t.Error("Error() should log 'error message'")
	}
}

// TestLoggerConcurrent tests concurrent logging
func TestLoggerConcurrent(t *testing.T) {
	config := DefaultConfig()
	logger, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer logger.Close()

	var wg sync.WaitGroup
	numGoroutines := 100
	messagesPerGoroutine := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < messagesPerGoroutine; j++ {
				logger.Infof("Goroutine %d, message %d", id, j)
			}
		}(i)
	}

	wg.Wait()
	// If we reach here without panic or deadlock, the test passes
}

// TestSetLevel tests dynamically changing the log level
func TestSetLevel(t *testing.T) {
	var buf bytes.Buffer

	config := Config{
		Level:        LevelInfo,
		Format:       "txt",
		EnableCaller: false,
		EnableColor:  false,
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	logger.writer = &bufferWriter{buf: &buf}

	// Debug should not be logged initially
	logger.Debug("debug message")
	if buf.Len() > 0 {
		t.Error("Debug() should not log when Level is Info")
	}

	// Change level to Debug
	logger.SetLevel(LevelDebug)

	// Now Debug should be logged
	logger.Debug("debug message")
	if !strings.Contains(buf.String(), "debug message") {
		t.Error("Debug() should log after SetLevel(LevelDebug)")
	}
}

// TestSetFormat tests dynamically changing the format
func TestSetFormat(t *testing.T) {
	config := DefaultConfig()
	logger, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Test switching to JSON format
	logger.SetFormat("json")

	var buf bytes.Buffer
	logger.writer = &bufferWriter{buf: &buf}

	logger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "{") || !strings.Contains(output, "}") {
		t.Error("SetFormat(json) should produce JSON output")
	}
}

// TestStdlibCompatibility tests standard library compatibility
func TestStdlibCompatibility(t *testing.T) {
	var buf bytes.Buffer

	// Create a custom logger for testing
	config := Config{
		Level:        LevelInfo,
		Format:       "txt",
		EnableCaller: false,
		EnableColor:  false,
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	logger.writer = &bufferWriter{buf: &buf}

	stdLogger := NewStdLogger(logger)

	// Test Print
	stdLogger.Print("test print")
	if !strings.Contains(buf.String(), "test print") {
		t.Error("Print() should log message")
	}

	buf.Reset()

	// Test Printf
	stdLogger.Printf("test %s", "printf")
	if !strings.Contains(buf.String(), "test printf") {
		t.Error("Printf() should log formatted message")
	}
}

// bufferWriter is a test helper that writes to a bytes.Buffer
type bufferWriter struct {
	buf *bytes.Buffer
	mu  sync.Mutex
}

func (w *bufferWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.Write(p)
}

func (w *bufferWriter) Close() error {
	return nil
}
