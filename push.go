// +build go1.8

package handlers

import (
	"net/http"
)

func push(w http.ResponseWriter, target string, opts *http.PushOptions) error {
	if pusher, ok := w.(http.Pusher); ok {
		return pusher.Push(target, opts)
	}
	return http.ErrNotSupported
}

// Push implements the http.Pusher interface.
func (s *startWriter) Push(target string, opts *http.PushOptions) error {
	return push(s.w, target, opts)
}

// Push implements the http.Pusher interface.
func (l *responseLogger) Push(target string, opts *http.PushOptions) error {
	return push(l.w, target, opts)
}

// Push implements the http.Pusher interface.
func (s *serverWriter) Push(target string, opts *http.PushOptions) error {
	return push(s.w, target, opts)
}
