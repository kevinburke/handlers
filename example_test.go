package handlers_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"

	"github.com/kevinburke/handlers"
)

func Example() {
	mux := http.NewServeMux()
	h := handlers.Duration(handlers.Server(mux, "custom-server"))
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	fmt.Println(w.Header().Get("Server"))
	// Output: custom-server
}

func ExampleRegexp() {
	// GET /v1/jobs/:job-name
	route := regexp.MustCompile(`^/v1/jobs/(?P<JobName>[^\s\/]+)$`)

	h := new(handlers.Regexp)
	h.HandleFunc(route, []string{"GET", "POST"}, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello World!")
	})
}
