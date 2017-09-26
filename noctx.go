// +build !go1.7

package handlers

import (
	"net/http"
	"time"

	"github.com/aristanetworks/goarista/monotime"
	"github.com/satori/go.uuid"
)

// SetRequestID sets the given UUID on the request and returns the modified
// HTTP request.
func SetRequestID(r *http.Request, u uuid.UUID) *http.Request {
	r.Header.Set("X-Request-Id", u.String())
	return r
}

// GetRequestID returns a UUID (if it exists on r) or false if none could
// be found.
func GetRequestID(r *http.Request) (uuid.UUID, bool) {
	rid := r.Header.Get("X-Request-Id")
	if rid != "" {
		u, err := uuid.FromString(rid)
		if err == nil {
			return u, true
		}
	}
	return uuid.UUID{}, false
}

type startWriter struct {
	w http.ResponseWriter
	// use this for durations
	monoStart   uint64
	wroteHeader bool
}

func (s *startWriter) duration() string {
	d := (monotime.Since(s.monoStart) / (100 * time.Microsecond)) * (100 * time.Microsecond)
	return d.String()
}

func (s *startWriter) WriteHeader(code int) {
	if s.wroteHeader == false {
		s.w.Header().Set("X-Request-Duration", s.duration())
		s.wroteHeader = true
	}
	s.w.WriteHeader(code)
}

func (s *startWriter) Write(b []byte) (int, error) {
	// Some chunked encoding transfers won't ever call WriteHeader(), so set
	// the header here.
	if s.wroteHeader == false {
		s.w.Header().Set("X-Request-Duration", s.duration())
		s.wroteHeader = true
	}
	return s.w.Write(b)
}

func (s *startWriter) Header() http.Header {
	return s.w.Header()
}

// Duration sets a X-Request-Duration header on the response. This header
// should go outside of any others to accurately capture the request duration.
func Duration(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := &startWriter{
			w:           w,
			monoStart:   monotime.Now(),
			wroteHeader: false,
		}
		h.ServeHTTP(sw, r)
	})
}
