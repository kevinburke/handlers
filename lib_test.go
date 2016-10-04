package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	uuid "github.com/satori/go.uuid"
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

func TestSetRequestID(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	u := uuid.NewV4()
	req = SetRequestID(req, u)
	fmt.Printf("%#v\n", req.Context())
	rid := req.Header.Get("X-Request-Id")
	if rid != u.String() {
		t.Errorf("expected X-Request-Id to equal %s, got %s", u.String(), rid)
	}
	val := req.Context().Value(requestID)
	v, ok := val.(uuid.UUID)
	if !ok {
		t.Fatalf("couldn't get requestID out of the request context")
	}
	if v.String() != u.String() {
		t.Errorf("expected %s (from context) to equal %s", v.String(), u.String())
	}
}

func TestGetRequestID(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	u := uuid.NewV4()
	_, ok := GetRequestID(req)
	if ok != false {
		t.Error("expected request id get to return false, got true")
	}
	req = SetRequestID(req, u)
	uid, ok := GetRequestID(req)
	if !ok {
		t.Error("expected request id get to return true, got false")
	}
	if uid.String() != u.String() {
		t.Errorf("expected %s (from context) to equal %s", uid.String(), u.String())
	}
}
