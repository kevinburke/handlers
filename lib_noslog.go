//go:build !go1.21
// +build !go1.21

package handlers

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/inconshreveable/log15"
	"github.com/kevinburke/rest"
)

// Logger is a logger configured to avoid the 40-char spacing gap between the
// message and the first key, and to write timestamps with full nanosecond
// precision.
var Logger log15.Logger

func init() {
	Logger = NewLogger()
	rest.Logger = Logger
}

type logHandler struct {
	h http.Handler
	l log15.Logger
}

// WithLogger logs information about HTTP requests and responses to the
// provided Logger, including a detailed timestamp, the user agent, the
// response time, the number of bytes written, and more. Any errors writing log
// information to the Logger are ignored.
func WithLogger(h http.Handler, logger log15.Logger) http.Handler {
	return &logHandler{h, logger}
}

func writeLog(l log15.Logger, r *http.Request, u url.URL, t time.Time, status int, size int) {
	user, _, _ := r.BasicAuth()
	args := []interface{}{
		"method", r.Method,
		"path", r.URL.RequestURI(),
		"time", strconv.FormatInt(timeSinceMs(t), 10),
		"bytes", strconv.Itoa(size),
		"status", strconv.Itoa(status),
		// Set X-Forwarded-For to pass through headers from a proxy.
		"remote_addr", getRemoteIP(r),
		"host", r.Host,
		"user_agent", r.UserAgent(),
	}
	if user != "" {
		args = append(args, "user", user)
	}
	if id := r.Header.Get("X-Request-Id"); id != "" {
		args = append(args, "request_id", id)
	}
	holder := r.Context().Value(extraLog).(*logHolder)
	args = append(args, holder.logs...)
	l.Info("", args...)
}
