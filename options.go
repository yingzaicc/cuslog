package cuslog

import (
	"io"
)

const (
	FmtEmptySeparator = ""
)

type Level uint8

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelPanic
	LevelFatal
)

type options struct {
	output        io.Writer
	level         Level
	stdLevel      Level
	formmatter    Formatter
	disableCaller bool
}

type Option func(*options)
