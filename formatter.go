package cuslog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Formatter interface defines how log entries are formatted
type Formatter interface {
	Format(entry *Entry) ([]byte, error)
}

// TxtFormatter formats log entries as plain text
type TxtFormatter struct {
	EnableColor bool
}

// Format formats a log entry as plain text
func (f *TxtFormatter) Format(entry *Entry) ([]byte, error) {
	var buf bytes.Buffer

	// Format time: 2023-10-01 12:00:00
	buf.WriteString(entry.Time.Format("2006-01-02 15:04:05"))
	buf.WriteString(" ")

	// Format level with optional color
	level := entry.Level.String()
	if f.EnableColor {
		level = f.colorizeLevel(entry.Level, level)
	}
	buf.WriteString(fmt.Sprintf("[%s]", level))
	buf.WriteString(" ")

	// Add caller info if available
	if entry.File != "" {
		buf.WriteString(fmt.Sprintf("%s:%d: ", entry.File, entry.Line))
	}

	// Add message
	buf.WriteString(entry.Message)

	// Add structured data if available
	if len(entry.Data) > 0 {
		buf.WriteString(" ")
		data, err := json.Marshal(entry.Data)
		if err == nil {
			buf.WriteString(string(data))
		}
	}

	buf.WriteString("\n")
	return buf.Bytes(), nil
}

// colorizeLevel adds ANSI color codes to the level string
func (f *TxtFormatter) colorizeLevel(level Level, levelStr string) string {
	var colorCode int
	switch level {
	case LevelDebug:
		colorCode = 36 // Cyan
	case LevelInfo:
		colorCode = 32 // Green
	case LevelWarn:
		colorCode = 33 // Yellow
	case LevelError:
		colorCode = 31 // Red
	case LevelFatal:
		colorCode = 35 // Magenta
	case LevelPanic:
		colorCode = 35 // Magenta
	default:
		return levelStr
	}
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", colorCode, levelStr)
}

// JsonFormatter formats log entries as JSON
type JsonFormatter struct{}

// Format formats a log entry as JSON
func (f *JsonFormatter) Format(entry *Entry) ([]byte, error) {
	// Create an ordered struct to preserve field order
	type jsonLog struct {
		Time    string                 `json:"time"`
		Level   string                 `json:"level"`
		File    string                 `json:"file,omitempty"`
		Line    int                    `json:"line,omitempty"`
		Message string                 `json:"msg"`
		Data    map[string]interface{} `json:"-"`
	}

	jsonEntry := jsonLog{
		Time:    entry.Time.Format(time.RFC3339),
		Level:   entry.Level.String(),
		Message: entry.Message,
	}

	// Add caller info if available
	if entry.File != "" {
		jsonEntry.File = entry.File
		jsonEntry.Line = entry.Line
	}

	// Add structured data by merging with main fields
	if len(entry.Data) > 0 {
		// Create a map with all fields in correct order
		fieldsMap := make(map[string]interface{})
		fieldsMap["time"] = jsonEntry.Time
		fieldsMap["level"] = jsonEntry.Level
		if jsonEntry.File != "" {
			fieldsMap["file"] = jsonEntry.File
			fieldsMap["line"] = jsonEntry.Line
		}
		fieldsMap["msg"] = jsonEntry.Message

		// Add custom data fields
		for k, v := range entry.Data {
			fieldsMap[k] = v
		}

		return json.Marshal(fieldsMap)
	}

	// Marshal to JSON (without custom data)
	return json.Marshal(jsonEntry)
}

// GetFormatter returns the appropriate formatter based on format string
func GetFormatter(format string, enableColor bool) Formatter {
	switch strings.ToLower(format) {
	case "json":
		return &JsonFormatter{}
	case "txt":
		return &TxtFormatter{EnableColor: enableColor}
	default:
		return &TxtFormatter{EnableColor: enableColor}
	}
}
