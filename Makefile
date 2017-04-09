BUMP_VERSION := $(shell command -v bump_version)
STATICCHECK := $(shell command -v staticcheck)

SHELL = /bin/bash

vet:
ifndef STATICCHECK
	go get -u honnef.co/go/tools/cmd/staticcheck
endif
	go vet ./...
	staticcheck ./...

test: vet
	go test ./...

race-test: vet
	go test -race ./...

install:
	go install ./...

release: race-test
ifndef BUMP_VERSION
	go get github.com/Shyp/bump_version
endif
	bump_version minor lib.go
