// +build go1.7

package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
)

func TestGetRequestID(t *testing.T) {
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
	req, _ := http.NewRequest("GET", "/", nil)
	d := GetDuration(req.Context())
	if d != time.Duration(0) {
		t.Errorf("expected Duration to be 0, got %v", d)
	}
	w := httptest.NewRecorder()
	h := Duration(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ensure we don't round down to 0
		time.Sleep(1 * time.Millisecond)
		w.Write([]byte("hello world"))
	}))
	h.ServeHTTP(w, req)
	dur, err := time.ParseDuration(w.Header().Get("X-Request-Duration"))
	if err != nil {
		t.Fatal(err)
	}
	if dur == 0 {
		t.Errorf("got 0 duration, wanted a greater than 0 duration")
	}
	if dur > 5*time.Millisecond {
		t.Errorf("got a duration greater than 5ms: %v", dur)
	}
}
