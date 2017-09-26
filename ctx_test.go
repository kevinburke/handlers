package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
)

func TestGetRequestID(t *testing.T) {
	t.Parallel()
	req, _ := http.NewRequest("GET", "/", nil)
	u := uuid.NewV4()
	_, ok := GetRequestID(req.Context())
	if ok != false {
		t.Error("expected request id get to return false, got true")
	}
	req = SetRequestID(req, u)
	uid, ok := GetRequestID(req.Context())
	if !ok {
		t.Error("expected request id get to return true, got false")
	}
	if uid.String() != u.String() {
		t.Errorf("expected %s (from context) to equal %s", uid.String(), u.String())
	}
}

func TestGetDuration(t *testing.T) {
	t.Parallel()
	req, _ := http.NewRequest("GET", "/", nil)
	d := GetDuration(req.Context())
	if d != time.Duration(0) {
		t.Errorf("expected Duration to be 0, got %v", d)
	}
	w := httptest.NewRecorder()
	h := Duration(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ensure we don't round down to 0
		time.Sleep(1 * time.Millisecond)
		d = GetDuration(r.Context())
		if d == 0 {
			t.Errorf("got 0 duration, wanted a greater than 0 duration")
		}
		if d > 50*time.Millisecond {
			t.Errorf("got a duration greater than 50ms: %v", d)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("hello world"))
	}))
	h.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Errorf("expected Code to be 400, got %d", w.Code)
	}
	dur, err := time.ParseDuration(w.Header().Get("X-Request-Duration"))
	if err != nil {
		t.Fatal(err)
	}
	if dur == 0 {
		t.Errorf("got 0 duration, wanted a greater than 0 duration")
	}
	if dur > 50*time.Millisecond {
		t.Errorf("got a duration greater than 50ms: %v", dur)
	}
}

func TestSetRequestID(t *testing.T) {
	t.Parallel()
	req, _ := http.NewRequest("GET", "/", nil)
	u := uuid.NewV4()
	req = SetRequestID(req, u)
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

func TestWithTimeout(t *testing.T) {
	t.Parallel()
	req, _ := http.NewRequest("GET", "/", nil)
	h := WithTimeout(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deadline, ok := r.Context().Deadline()
		if !ok {
			t.Error("expected Deadline() to be ok, got not ok")
		}
		if time.Until(deadline) > 10*time.Millisecond {
			t.Errorf("too big of a deadline: %v", deadline)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("hello world"))
	}), 10*time.Millisecond)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Errorf("expected Code to be 400, got %d", w.Code)
	}
}
