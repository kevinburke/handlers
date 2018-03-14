// Package handlers implements a number of useful HTTP middlewares.
//
// The general format of the middlewares in this package is to wrap an existing
// http.Handler in another one. So if you have a ServeMux, you can simply do:
//
//     mux := http.NewServeMux()
//     h := handlers.Log(handlers.Debug(mux))
//     http.ListenAndServe(":5050", h)
//
// And wrap as many handlers as you'd like using that idiom.
package handlers

import (
	"bytes"
	"crypto/subtle"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/inconshreveable/log15"
	"github.com/kevinburke/go.uuid"
	"github.com/kevinburke/rest"
)

const Version = "0.39"

func push(w http.ResponseWriter, target string, opts *http.PushOptions) error {
	if pusher, ok := w.(http.Pusher); ok {
		return pusher.Push(target, opts)
	}
	return http.ErrNotSupported
}

// All wraps h with every handler in this file.
func All(h http.Handler, serverName string) http.Handler {
	return Duration(Log(Debug(UUID(TrailingSlashRedirect(JSON(Server(h, serverName)))))))
}

// JSON sets the Content-Type to application/json; charset=utf-8.
func JSON(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		h.ServeHTTP(w, r)
	})
}

type serverWriter struct {
	w           http.ResponseWriter
	name        string
	wroteHeader bool
}

func (s *serverWriter) WriteHeader(code int) {
	if s.wroteHeader == false {
		s.w.Header().Set("Server", s.name)
		s.wroteHeader = true
	}
	s.w.WriteHeader(code)
}

func (s *serverWriter) Write(b []byte) (int, error) {
	if s.wroteHeader == false {
		s.w.Header().Set("Server", s.name)
		s.wroteHeader = true
	}
	return s.w.Write(b)
}

func (s *serverWriter) Header() http.Header {
	return s.w.Header()
}

// Push implements the http.Pusher interface.
func (s *serverWriter) Push(target string, opts *http.PushOptions) error {
	return push(s.w, target, opts)
}

// TrailingSlashRedirect redirects any path that ends with a "/" - say,
// "/messages/" - to the stripped version, say "/messages".
func TrailingSlashRedirect(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) > 1 && strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, r.URL.Path[:len(r.URL.Path)-1], http.StatusMovedPermanently)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// Server attaches a Server header to the response.
func Server(h http.Handler, serverName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := &serverWriter{
			w:           w,
			name:        serverName,
			wroteHeader: false,
		}
		h.ServeHTTP(sw, r)
		if sw.wroteHeader == false {
			sw.w.Header().Set("Server", sw.name)
			sw.wroteHeader = true
		}
	})
}

// UUID attaches a X-Request-Id header to the request, and sets one on the
// request context, unless one already exists.
func UUID(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get("X-Request-Id")
		if rid == "" {
			r = SetRequestID(r, uuid.NewV4())
		}
		h.ServeHTTP(w, r)
	})
}

// BasicAuth protects all requests to the given handler, unless the request has
// basic auth with a username and password in the users map.
func BasicAuth(h http.Handler, realm string, users map[string]string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			rest.Unauthorized(w, r, realm)
			return
		}

		serverPass, ok := users[user]
		if !ok {
			if user == "" {
				rest.Unauthorized(w, r, realm)
			} else {
				rest.Forbidden(w, r, &rest.Error{
					Title: "Username or password are invalid. Please double check your credentials",
					ID:    "forbidden",
				})
			}
			return
		}
		if subtle.ConstantTimeCompare([]byte(pass), []byte(serverPass)) != 1 {
			rest.Forbidden(w, r, &rest.Error{
				Title:    fmt.Sprintf("Incorrect password for user %s", user),
				ID:       "incorrect_password",
				Instance: r.URL.Path,
			})
			return
		}
		h.ServeHTTP(w, r)
	})
}

// Debug prints debugging information about the request to stdout if the
// DEBUG_HTTP_TRAFFIC environment variable is set to true.
func Debug(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("DEBUG_HTTP_TRAFFIC") == "true" {
			// You need to write the entire thing in one Write, otherwise the
			// output will be jumbled with other requests.
			b := new(bytes.Buffer)
			bits, err := httputil.DumpRequest(r, true)
			if err != nil {
				_, _ = b.WriteString(err.Error())
			} else {
				if w.Header().Get("Content-Encoding") == "gzip" {
					_, _ = b.WriteString("[binary data omitted]")
				} else {
					_, _ = b.Write(bits)
				}
			}
			res := httptest.NewRecorder()
			h.ServeHTTP(res, r)

			_, _ = b.WriteString(fmt.Sprintf("HTTP/1.1 %d\r\n", res.Code))
			_ = res.HeaderMap.Write(b)
			for k, v := range res.HeaderMap {
				w.Header()[k] = v
			}
			w.WriteHeader(res.Code)
			_, _ = b.WriteString("\r\n")
			writer := io.MultiWriter(w, b)
			_, _ = res.Body.WriteTo(writer)
			if w.Header().Get("Content-Encoding") == "gzip" {
				os.Stderr.WriteString("[binary data omitted]")
			} else {
				_, _ = b.WriteTo(os.Stderr)
			}
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

// responseLogger is wrapper of http.ResponseWriter that keeps track of its HTTP
// status code and body size
type responseLogger struct {
	w      http.ResponseWriter
	status int
	size   int
}

func (l *responseLogger) Header() http.Header {
	return l.w.Header()
}

func (l *responseLogger) Write(b []byte) (int, error) {
	if l.status == 0 {
		// The status will be StatusOK if WriteHeader has not been called yet
		l.status = http.StatusOK
	}
	size, err := l.w.Write(b)
	l.size += size
	return size, err
}

func (l *responseLogger) WriteHeader(s int) {
	l.w.WriteHeader(s)
	l.status = s
}

func (l *responseLogger) Status() int {
	if l.status == 0 {
		return http.StatusOK // default status
	}
	return l.status
}

func (l *responseLogger) Size() int {
	return l.size
}

func (l *responseLogger) Flush() {
	f, ok := l.w.(http.Flusher)
	if ok {
		f.Flush()
	}
}

// Push implements the http.Pusher interface.
func (l *responseLogger) Push(target string, opts *http.PushOptions) error {
	return push(l.w, target, opts)
}

type hijackLogger struct {
	responseLogger
}

type hijackCloseNotifier struct {
	loggingResponseWriter
	http.Hijacker
	http.CloseNotifier
}

type closeNotifyWriter struct {
	loggingResponseWriter
	http.CloseNotifier
}

func makeLogger(w http.ResponseWriter) loggingResponseWriter {
	var logger loggingResponseWriter = &responseLogger{w: w}
	if _, ok := w.(http.Hijacker); ok {
		logger = &hijackLogger{responseLogger{w: w}}
	}
	h, ok1 := logger.(http.Hijacker)
	c, ok2 := w.(http.CloseNotifier)
	if ok1 && ok2 {
		return hijackCloseNotifier{logger, h, c}
	}
	if ok2 {
		return &closeNotifyWriter{logger, c}
	}
	return logger
}

type loggingResponseWriter interface {
	http.ResponseWriter
	http.Flusher
	Status() int
	Size() int
}

type logHandler struct {
	h http.Handler
	l log.Logger
}

func (l logHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	logWriter := makeLogger(w)
	u := *r.URL
	l.h.ServeHTTP(logWriter, r)
	writeLog(l.l, r, u, t, logWriter.Status(), logWriter.Size())
}

func getRemoteIP(r *http.Request) string {
	fwd := r.Header.Get("X-Forwarded-For")
	if fwd == "" {
		return r.RemoteAddr
	}
	return strings.Split(fwd, ",")[0]
}

// Return the time since the given time, in ms.
func timeSinceMs(t time.Time) int64 {
	// Add 500 microseconds so we round up or down to the nearest MS.
	ns := time.Since(t).Nanoseconds() + 500*int64(time.Microsecond)
	return ns / int64(time.Millisecond)
}

func writeLog(l log.Logger, r *http.Request, u url.URL, t time.Time, status int, size int) {
	user, _, _ := r.BasicAuth()
	args := []interface{}{
		"method", r.Method,
		"path", r.URL.RequestURI(),
		"time", strconv.FormatInt(timeSinceMs(t), 10),
		"bytes", strconv.Itoa(size),
		"status", strconv.Itoa(status),
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
	l.Info("", args...)
}

// Log serves the http request and writes information about the
// request/response using the default Logger (handlers.Logger). Any errors
// writing to the Logger are ignored.
func Log(h http.Handler) http.Handler {
	return WithLogger(h, Logger)
}

// WithLogger logs information about HTTP requests and responses to the
// provided Logger, including a detailed timestamp, the user agent, the
// response time, the number of bytes written, and more. Any errors writing log
// information to the Logger are ignored.
func WithLogger(h http.Handler, logger log.Logger) http.Handler {
	return &logHandler{h, logger}
}

// RedirectProto redirects requests with an "X-Forwarded-Proto: http" header to
// their HTTPS equivalent.
func RedirectProto(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Forwarded-Proto") == "http" {
			r2 := new(http.Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Scheme = "https"
			r2.URL.Host = r.Host
			http.Redirect(w, r2, r2.URL.String(), http.StatusFound)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// STS sets a Strict-Transport-Security header on the response.
func STS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; preload")
		h.ServeHTTP(w, r)
	})
}
