//go:build ignore
// +build ignore

package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/NYTimes/gziphandler"
)

// Not sure this test is super useful and this means that we need to import the
// nytimes library everywhere, so ignore this test for now.

func TestDebugGzip(t *testing.T) {
	envFunc = func(s string) string {
		if s == "DEBUG_HTTP_SERVER_TRAFFIC" {
			return "true"
		}
		return ""
	}
	defer func() { envFunc = os.Getenv }()
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := strings.Repeat("a", 2000)
		io.WriteString(w, resp)
	})
	gzipper := gziphandler.GzipHandler(h)
	r := httptest.NewRequest("GET", "/", nil)
	debugBuf := new(bytes.Buffer)
	w := httptest.NewRecorder()
	r.Header.Set("Accept-Encoding", "gzip")
	// Log(Debug(gzipper)).ServeHTTP(w, r)
	DebugWriter(gzipper, debugBuf).ServeHTTP(w, r)
	if w.Body.Len() != 39 {
		t.Errorf("expected body length to be 39, got %d", w.Body.Len())
	}
	if w.Header().Get("Content-Encoding") != "gzip" {
		t.Errorf("bad value for content encoding")
	}
	io.Copy(os.Stderr, debugBuf)
	io.WriteString(os.Stderr, "\n\n")
}
