//go:build go1.21
// +build go1.21

package handlers

import (
	"log/slog"
	"testing"
)

func TestColorHandler(t *testing.T) {
	ch := NewColorHandler()
	l := slog.New(ch)
	l.Info("test message", "foo", "bar", "one", 2, "two", []int{3, 4, 5})
	l.Warn("warn message", "foo", "bar", "one", 2, "two", []int{3, 4, 5})
	l.Error("error message", "foo", "bar", "one", 2, "two", []int{3, 4, 5})
	l.Debug("error message", "foo", "bar", "one", 2, "two", []int{3, 4, 5})
}
