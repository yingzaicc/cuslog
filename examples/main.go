package main

import (
	"fmt"
	"time"

	l "cuslog"
)

func main() {
	// Example 1: Using default logger
	fmt.Println("=== Example 1: Default Logger ===")
	l.Debug("This is a debug message (won't show - default level is Info)")
	l.Info("This is an info message")
	l.Warn("This is a warning message")
	l.Error("This is an error message")

	// Example 2: Creating a custom logger with console output
	fmt.Println("\n=== Example 2: Custom Console Logger ===")
	consoleConfig := l.Config{
		Level:        l.LevelDebug,
		Format:       "txt",
		Outputs:      []l.OutputConfig{{Type: "console"}},
		EnableCaller: true,
		EnableColor:  true,
	}
	consoleLogger, err := l.New(consoleConfig)
	if err != nil {
		panic(err)
	}
	consoleLogger.Debug("Debug message with caller info")
	consoleLogger.Infof("Info message: %s", "formatted output")
	consoleLogger.Error("Error message from custom logger")

	// Example 3: Creating a logger with JSON format
	fmt.Println("\n=== Example 3: JSON Format Logger ===")
	jsonConfig := l.Config{
		Level:        l.LevelInfo,
		Format:       "json",
		Outputs:      []l.OutputConfig{{Type: "console"}},
		EnableCaller: true,
		EnableColor:  false,
	}
	jsonLogger, err := l.New(jsonConfig)
	if err != nil {
		panic(err)
	}
	jsonLogger.Info("JSON formatted info message")
	jsonLogger.Warn("JSON formatted warning message")

	// Example 4: Creating a logger with file output
	fmt.Println("\n=== Example 4: File Output Logger ===")
	fileConfig := l.Config{
		Level:  l.LevelDebug,
		Format: "txt",
		Outputs: []l.OutputConfig{
			{Type: "console"},
			{Type: "file", FilePath: "app.log"},
		},
		EnableCaller: true,
		EnableColor:  true,
	}
	fileLogger, err := l.New(fileConfig)
	if err != nil {
		panic(err)
	}
	defer fileLogger.Close()

	fileLogger.Info("This message goes to both console and file")
	fileLogger.Debugf("Debug details: timestamp=%s", time.Now().Format(time.RFC3339))
	fmt.Println("Check app.log for file output")

	// Example 5: Dynamic level changing
	fmt.Println("\n=== Example 5: Dynamic Level Change ===")
	logger, _ := l.New(l.DefaultConfig())
	logger.Info("Message before level change (Info level)")
	logger.SetLevel(l.LevelDebug)
	logger.Debug("Debug message after level change (now visible)")
	logger.SetLevel(l.LevelWarn)
	logger.Debug("Debug message won't show (level raised to Warn)")
	logger.Warn("Warning message will show")

	// Example 6: Dynamic format changing
	fmt.Println("\n=== Example 6: Dynamic Format Change ===")
	logger2, _ := l.New(l.DefaultConfig())
	logger2.Info("Text format message")
	logger2.SetFormat("json")
	logger2.Info("JSON format message")

	// Example 7: Standard library compatibility
	fmt.Println("\n=== Example 7: Standard Library Compatibility ===")
	stdLogger := l.NewStdLogger(consoleLogger)
	stdLogger.Print("Using stdlib Print method")
	stdLogger.Printf("Using stdlib Printf: %s\n", "formatted")
	stdLogger.Println("Using stdlib Println method")

	// Example 8: Using different log levels
	fmt.Println("\n=== Example 8: Different Log Levels ===")
	levelLogger, _ := l.New(l.DefaultConfig())
	levelLogger.Debug("Detailed debugging information")
	levelLogger.Info("General information")
	levelLogger.Warn("Warning message")
	levelLogger.Error("Error occurred")

	// Example 9: Concurrent logging
	fmt.Println("\n=== Example 9: Concurrent Logging ===")
	concurrentLogger, _ := l.New(l.DefaultConfig())
	for i := 0; i < 5; i++ {
		go func(id int) {
			concurrentLogger.Infof("Goroutine %d is logging", id)
		}(i)
	}
	time.Sleep(100 * time.Millisecond) // Wait for goroutines to complete

	// Example 10: Panic and Fatal
	fmt.Println("\n=== Example 10: Panic Handling ===")
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
		}
	}()

	// Note: We wrap this in a function to prevent the panic from stopping the entire program
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Panic was caught by defer")
			}
		}()
		panicLogger, _ := l.New(l.Config{
			Level:        l.LevelInfo,
			Format:       "txt",
			Outputs:      []l.OutputConfig{{Type: "console"}},
			EnableCaller: true,
			EnableColor:  false,
		})
		panicLogger.Panic("This is a panic message")
	}()

	fmt.Println("\n=== All Examples Completed ===")
}
