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

type testServer bool

var msg = ("{\"message\": \"Hello World\"}")

func (ts testServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg))
}

func TestJSON(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	ts := testServer(false)
	JSON(ts).ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != msg {
		t.Errorf("expected \"%s\", got %s", msg, w.Body.String())
	}
	if w.Header().Get("Content-Type") != "application/json; charset=utf-8" {
		t.Errorf("expected content-type \"application/json\", got %s", w.Header().Get("Content-Type"))
	}
}

func TestServer(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	ts := testServer(false)
	Server(ts, "foo bar").ServeHTTP(w, req)
	if w.Header().Get("Server") != "foo bar" {
		t.Errorf("expected server header \"foo bar\", got %s", w.Header().Get("Server"))
	}
}

func TestServerOverwritesInnerHeaders(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	sh := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Server", "Inner header!")
		testServer(false).ServeHTTP(w, r)
	})
	Server(sh, "foo bar").ServeHTTP(w, req)
	if w.Header().Get("Server") != "foo bar" {
		t.Errorf("expected server header \"foo bar\", got %s", w.Header().Get("Server"))
	}
}

func TestAll(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	ts := testServer(false)
	All(ts, "foo bar").ServeHTTP(w, req)
	if w.Header().Get("Server") != "foo bar" {
		t.Errorf("expected server header \"foo bar\", got %s", w.Header().Get("Server"))
	}
	if w.Header().Get("Content-Type") != "application/json; charset=utf-8" {
		t.Errorf("expected content-type \"application/json\", got %s", w.Header().Get("Content-Type"))
	}
}

func TestTrailingSlashRedirect(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	ts := TrailingSlashRedirect(testServer(false))
	ts.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("expected Code to be 200, got %d", w.Code)
	}
	req = httptest.NewRequest("GET", "/trailingslash//////", nil)
	w = httptest.NewRecorder()
	ts.ServeHTTP(w, req)
	if w.Code != 301 {
		t.Errorf("expected Code to be 301, got %d", w.Code)
	}
	location := w.Header().Get("Location")
	if location != "/trailingslash/" {
		t.Errorf("expected Location header to be /trailingslash/, got %s", location)
	}
	req = httptest.NewRequest("GET", "/trailingslash/", nil)
	w = httptest.NewRecorder()
	ts.ServeHTTP(w, req)
	location = w.Header().Get("Location")
	if location != "/trailingslash" {
		t.Errorf("expected Location header to be /trailingslash, got %s", location)
	}
}

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
}
