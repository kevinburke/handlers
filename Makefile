SHELL = /bin/bash -o pipefail

BUMP_VERSION := $(GOPATH)/bin/bump_version

vet:
	go vet ./...
	staticcheck ./...

test: vet
	go test -timeout=10s ./...

install-ci:
	GO111MODULE=on go install honnef.co/go/tools/cmd/staticcheck@latest
	GO111MODULE=on go install github.com/kevinburke/goget@latest
	goget -https github.com/gofrs/uuid
	goget -https github.com/inconshreveable/log15
	goget -https github.com/kevinburke/rest
	goget -https github.com/mattn/go-colorable
	goget -https github.com/mattn/go-isatty
	goget -https golang.org/x/term
	goget -https golang.org/x/sys

ci: install-ci vet race-test

race-test: vet
	go test -race -timeout=10s ./...

install:
	go install ./...

$(BUMP_VERSION):
	go get github.com/kevinburke/bump_version

release: race-test | $(BUMP_VERSION)
	$(BUMP_VERSION) --tag-prefix=v minor lib.go
