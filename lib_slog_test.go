//go:build go1.21
// +build go1.21

package handlers

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLog(t *testing.T) {
	var buf bytes.Buffer
	h := slog.NewTextHandler(&buf, nil)
	logger := slog.New(h)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		AppendLog(r, "via", "test", "user", 7)
		AppendLog(r, "more", true)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(msg))
	})
	httpHandler := WithLogger(mux, logger)
	w := httptest.NewRecorder()
	w.Header().Set("User-Agent", "kevinburke/handlers")
	r := httptest.NewRequest("GET", "/", nil)
	httpHandler.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Errorf("expected 200 back, got %d", w.Code)
	}
	if w.Body.String() != msg {
		t.Errorf("bad response body")
	}
	if !strings.HasSuffix(buf.String(), "via=test user=7 more=true\n") {
		t.Errorf("did not log additional data to log: %q", buf.String())
	}
}
