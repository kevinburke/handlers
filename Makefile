BUMP_VERSION := $(GOPATH)/bin/bump_version
MEGACHECK := $(GOPATH)/bin/megacheck

SHELL = /bin/bash -o pipefail

vet: | $(MEGACHECK)
	go vet ./...
	$(MEGACHECK) --ignore='github.com/kevinburke/handlers/*.go:S1002' ./...

$(MEGACHECK):
	go get honnef.co/go/tools/cmd/megacheck

test: vet
	bazel test --test_output=errors //...

race-test: vet
	go test -race ./...

install:
	go install ./...

$(BUMP_VERSION):
	go get github.com/Shyp/bump_version

release: race-test | $(BUMP_VERSION)
	bump_version minor lib.go

ci:
	bazel --batch --host_jvm_args=-Dbazel.DigestFunction=SHA1 test \
		--experimental_repository_cache="$$HOME/.bzrepos" \
		--spawn_strategy=remote \
		--remote_rest_cache=https://remote.rest.stackmachine.com/cache \
		--test_output=errors \
		--strategy=Javac=remote \
		--noshow_progress \
		--noshow_loading_progress \
		--features=race //... 2>&1 | ts '[%Y-%m-%d %H:%M:%.S]'
