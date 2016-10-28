package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
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
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
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
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	ts := testServer(false)
	Server(ts, "foo bar").ServeHTTP(w, req)
	if w.Header().Get("Server") != "foo bar" {
		t.Errorf("expected server header \"foo bar\", got %s", w.Header().Get("Server"))
	}
}

func TestServerOverwritesInnerHeaders(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
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
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
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
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	ts := TrailingSlashRedirect(testServer(false))
	ts.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("expected Code to be 200, got %d", w.Code)
	}
	req, _ = http.NewRequest("GET", "/trailingslash//////", nil)
	w = httptest.NewRecorder()
	ts.ServeHTTP(w, req)
	if w.Code != 301 {
		t.Errorf("expected Code to be 301, got %d", w.Code)
	}
	location := w.Header().Get("Location")
	if location != "/trailingslash/" {
		t.Errorf("expected Location header to be /trailingslash/, got %s", location)
	}
	req, _ = http.NewRequest("GET", "/trailingslash/", nil)
	w = httptest.NewRecorder()
	ts.ServeHTTP(w, req)
	location = w.Header().Get("Location")
	if location != "/trailingslash" {
		t.Errorf("expected Location header to be /trailingslash, got %s", location)
	}
}
