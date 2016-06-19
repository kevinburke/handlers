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