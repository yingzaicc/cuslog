package cuslog

import "time"

// Entry represents a log entry with metadata
type Entry struct {
	Time    time.Time              // Log timestamp
	Level   Level                  // Log level
	Message string                 // Log message
	File    string                 // Caller file name
	Line    int                    // Caller line number
	Data    map[string]interface{} // Additional structured data
}
