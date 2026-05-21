### 0.49.0 - 2026-05-21

Update `github.com/inconshreveable/log15/v3` to v3.2.1 and
`golang.org/x/term` to v0.43.0.

Add Dependabot updates for Go modules and GitHub Actions, with a seven-day
cooldown for all update types.

Drop GOPATH-era CI setup and simplify the project Makefile and README.

Run the codebase through current Go formatting, removing legacy build tag
comments that are now covered by `//go:build` lines.

### 0.48.0 - 2026-03-05

Update `github.com/gofrs/uuid/v5` and other dependencies.

### 0.47.0 - 2025-07-18

Update to `github.com/inconshreveable/log15/v3`.

Update GitHub Actions to recent Go versions.

### 0.46 - 2023-11-07

Use `log/slog` for Go 1.21 and newer.

### 0.45 - 2023-09-24

Include the date in logger output.

### 0.44 - 2023-05-29

Use more precision for non-stdout log timestamps.

Update GitHub Actions to test against the latest two Go versions.

### 0.43 - 2022-11-21

Use `golang.org/x/term` instead of `github.com/mattn/go-isatty`.

Update GitHub Actions to test against more recent Go versions.

Run `go fmt` across the library.

### 0.42 - 2022-04-16

Update the GitHub Actions CI environment.

Ignore the New York Times host in `IgnoreLogging` tests.

### 0.41 - 2021-07-22

Change UUID handling from `github.com/kevinburke/go.uuid` to
`github.com/gofrs/uuid`.

Change the behavior of `handlers.AppendLog`.

### 0.40 - 2021-05-07

Add `AppendLog` support for appending to log lines during request processing.

Use the newer `Clone` function for request context handling.

Switch CI from Travis CI to GitHub Actions.

Add `DEBUG_HTTP_SERVER_TRAFFIC` and tests for the `Debug` handler.

Add `HandleString` and `HandleStringFunc` helpers for regex handlers.

Fix staticcheck and lint issues.

### 0.39 - 2018-03-14

Remove Bazel.

Write headers even if `Write` or `WriteHeader` is never called.

### 0.38 - 2018-03-06

If `WriteHeader` is never called, `Log` will log status=200 to the log,
instead of status=0.

### 0.37 - 2018-02-28

Support OPTIONS queries with a nil route.

### 0.36 - 2018-01-25

Point the UUID library at a fork.

### 0.35 - 2017-12-31

If two paths are declared that both match a given path, we'll try both of them
before giving up and returning a HTTP 405. This is slightly slower (we have to
try every route before giving up), but does the right thing.

### 0.34 - 2017-12-29

Pass nil to specify you want to handle all HTTP methods in the router.

### 0.33 - 2017-12-15

Make copies of all HTTP requests before passing them further down the chain. Per
the documentation, this is the correct way to handle this situation.

### 0.32 - 2017-11-16

Fix an error in the duration information reported by `X-Request-Duration`.
