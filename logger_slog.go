//go:build go1.21
// +build go1.21

package handlers

import (
	"log/slog"
	"os"

	"golang.org/x/term"
)

// NewLogger returns a new customizable Logger, with the same initial settings
// as the package Logger. Compared with a default log15.Logger, the 40-char
// spacing gap between the message and the first key is omitted, and timestamps
// are written with more precision.
func NewLogger() *slog.Logger {
	return NewLoggerLevel(slog.LevelInfo)
}

// NewLoggerLevel returns a Logger with our customized settings, set to log
// messages at or more critical than the given level.
func NewLoggerLevel(lvl slog.Level) *slog.Logger {
	// TODO add colorized logger
	if term.IsTerminal(int(os.Stdout.Fd())) {
		return slog.Default()
	} else {
		return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: lvl,
		}))
	}
}
