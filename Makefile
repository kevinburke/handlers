BUMP_VERSION := $(GOPATH)/bin/bump_version
MEGACHECK := $(GOPATH)/bin/megacheck

SHELL = /bin/bash

vet: | $(MEGACHECK)
	go vet ./...
	$(MEGACHECK) ./...

$(MEGACHECK):
	go get honnef.co/go/tools/cmd/megacheck

test: vet
	go test ./...

race-test: vet
	go test -race ./...

install:
	go install ./...

$(BUMP_VERSION):
	go get github.com/Shyp/bump_version

release: race-test | $(BUMP_VERSION)
	bump_version minor lib.go
