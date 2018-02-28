package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestRegexpHandler(t *testing.T) {
	t.Parallel()
	// GET /v1/jobs/:job-name
	route := regexp.MustCompile(`^/v1$`)

	h := new(Regexp)
	h.HandleFunc(route, []string{"GET", "POST"}, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello World!")
	})
	req := httptest.NewRequest("GET", "/v1", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Body.String() != "Hello World!" {
		t.Errorf("expected Body to equal 'Hello World!', got %v", w.Body.String())
	}

	req = httptest.NewRequest("PATCH", "/v1", nil)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 405 {
		t.Errorf("Expected PATCH request to return 405, got %d", w.Code)
	}

	req = httptest.NewRequest("POST", "/unknown", nil)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 404 {
		t.Errorf("Expected POST request to unknown route to return 404, got %d", w.Code)
	}
}

func TestHeadAllowed(t *testing.T) {
	t.Parallel()
	route := regexp.MustCompile(`^/v1$`)

	h := new(Regexp)
	h.HandleFunc(route, []string{"GET", "POST"}, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello World!")
	})
	req := httptest.NewRequest("HEAD", "/v1", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected HEAD request to return 200, got %d", w.Code)
	}
}

func TestNil(t *testing.T) {
	t.Parallel()
	route := regexp.MustCompile(`^/v1$`)

	h := new(Regexp)
	h.HandleFunc(route, nil, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello World!")
	})
	req := httptest.NewRequest("PATCH", "/v1", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected HEAD request to return 200, got %d", w.Code)
	}
}

func TestMatchingRoutes(t *testing.T) {
	t.Parallel()
	route := regexp.MustCompile(`^/v1$`)
	h := new(Regexp)
	h.HandleFunc(route, []string{"GET"}, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello Get World!")
	})
	h.HandleFunc(route, []string{"POST"}, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello Post World!")
	})
	req := httptest.NewRequest("POST", "/v1", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected POST request to return 200, got %d", w.Code)
	}
	if w.Body.String() != "Hello Post World!" {
		t.Errorf("Expected POST request to return body, got %s", w.Body)
	}
}

func TestOptionsNil(t *testing.T) {
	t.Parallel()
	route := regexp.MustCompile(`^/v1$`)

	h := new(Regexp)
	h.HandleFunc(route, nil, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
		io.WriteString(w, "Hello World!")
	})
	req := httptest.NewRequest("OPTIONS", "/v1", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected %s request to return 200, got %d", req.Method, w.Code)
	}
	if !strings.Contains(w.Header().Get("Allow"), "CONNECT") {
		t.Errorf("Expected ALLOW header to contain list of methods, got %q", w.Header().Get("Allow"))
	}
}
