package cuslog

// Config represents the logger configuration
type Config struct {
	Level        Level          // Log level (Debug/Info/Warn/Error/Fatal/Panic)
	Format       string         // Output format ("json"/"txt")
	Outputs      []OutputConfig // Output target configurations
	EnableCaller bool           // Whether to enable caller file path and line number
	EnableColor  bool           // Whether to enable console color output (txt format only)
}

// OutputConfig represents output target configuration
type OutputConfig struct {
	Type     string // Output type ("console"/"file"/"tcp"/"udp")
	FilePath string // File path (for file type only)
	Addr     string // Network address (for tcp/udp type only)
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		Level:        LevelInfo,
		Format:       "txt",
		Outputs:      []OutputConfig{{Type: "console"}},
		EnableCaller: true,
		EnableColor:  true,
	}
}
