package handlers

import (
	"net/http"
	"regexp"
)

func ExampleRegexp() {
	// GET /v1/jobs/:job-name
	route := regexp.MustCompile(`^/v1/jobs/(?P<JobName>[^\s\/]+)$`)

	h := new(Regexp)
	h.HandleFunc(route, []string{"GET", "POST"}, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
}
