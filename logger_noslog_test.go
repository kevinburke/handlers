//go:build !go1.21
// +build !go1.21

package handlers

import (
	"testing"

	"github.com/inconshreveable/log15"
	"github.com/mattn/go-colorable"
)

func TestLog15Handler(t *testing.T) {
	lg := log15.New()
	lh := log15.StreamHandler(colorable.NewColorableStdout(), termFormat())
	lg.SetHandler(lh)
	lg.Info("test message", "foo", "bar", "one", 2, "two", []int{3, 4, 5})
	lg.Warn("warn message", "foo", "bar", "one", 2, "two", []int{3, 4, 5})
	lg.Error("error message", "foo", "bar", "one", 2, "two", []int{3, 4, 5})
	lg.Debug("error message", "foo", "bar", "one", 2, "two", []int{3, 4, 5})
}
