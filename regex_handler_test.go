package handlers

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestRegexpHandler(t *testing.T) {
	t.Parallel()
	// GET /v1/jobs/:job-name
	route := regexp.MustCompile(`^/v1$`)

	h := new(Regexp)
	h.HandleFunc(route, []string{"GET", "POST"}, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	req, _ := http.NewRequest("GET", "/v1", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Body.String() != "Hello World!" {
		t.Errorf("expected Body to equal 'Hello World!', got %v", w.Body.String())
	}

	req, _ = http.NewRequest("PATCH", "/v1", nil)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 405 {
		t.Errorf("Expected PATCH request to return 405, got %d", w.Code)
	}

	req, _ = http.NewRequest("POST", "/unknown", nil)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 404 {
		t.Errorf("Expected POST request to unknown route to return 404, got %d", w.Code)
	}
}

func TestHeadAllowed(t *testing.T) {
	t.Parallel()
	// GET /v1/jobs/:job-name
	route := regexp.MustCompile(`^/v1$`)

	h := new(Regexp)
	h.HandleFunc(route, []string{"GET", "POST"}, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	req, _ := http.NewRequest("HEAD", "/v1", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected HEAD request to return 200, got %d", w.Code)
	}
}
