//go:build go1.21
// +build go1.21

package handlers

import (
	"context"
	"fmt"
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

type wrappedHandler struct {
	next slog.Handler
}

func (w *wrappedHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return w.next.Enabled(ctx, level)
}

func wrap(s string, color int) string {
	if color == 0 {
		return s
	}
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, s)
}

func (w *wrappedHandler) Handle(ctx context.Context, r slog.Record) error {
	color := 0
	switch r.Level {
	case slog.LevelError:
		color = 31
	case slog.LevelWarn:
		color = 33
	case slog.LevelInfo:
		color = 32
	case slog.LevelDebug:
		color = 36
	default:
		color = 35
	}

	r2 := slog.NewRecord(r.Time, r.Level, wrap(r.Message, color), r.PC)

	r.Attrs(func(a slog.Attr) bool {
		a.Key = wrap(a.Key, color)
		r2.AddAttrs(a)
		return true
	})
	return w.next.Handle(ctx, r2)
}

func (w *wrappedHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return w.next.WithAttrs(attrs)
}

func (w *wrappedHandler) WithGroup(name string) slog.Handler {
	return w.next.WithGroup(name)
}

func NewColorHandler() slog.Handler {
	h := slog.Default().Handler()
	return &wrappedHandler{next: h}
}
