SHELL = /bin/bash -o pipefail

BUMP_VERSION := $(GOPATH)/bin/bump_version
MEGACHECK := $(GOPATH)/bin/megacheck

BAZEL_VERSION := 0.9.0
BAZEL_DEB := bazel_$(BAZEL_VERSION)_amd64.deb

vet: | $(MEGACHECK)
	go vet ./...
	$(MEGACHECK) --ignore='github.com/kevinburke/handlers/*.go:S1002' ./...

$(MEGACHECK):
	go get honnef.co/go/tools/cmd/megacheck

test: vet
	bazel test --test_output=errors //...

race-test: vet
	bazel test --features=race --test_output=errors //...

install:
	go install ./...

$(BUMP_VERSION):
	go get github.com/Shyp/bump_version

release: race-test | $(BUMP_VERSION)
	bump_version minor lib.go

install-travis:
	wget "https://storage.googleapis.com/bazel-apt/pool/jdk1.8/b/bazel/$(BAZEL_DEB)"
	sudo dpkg --force-all -i $(BAZEL_DEB)
	sudo apt-get install moreutils -y

ci:
	bazel --batch --host_jvm_args=-Dbazel.DigestFunction=SHA1 test \
		--experimental_repository_cache="$$HOME/.bzrepos" \
		--spawn_strategy=remote \
		--test_output=errors \
		--strategy=Javac=remote \
		--noshow_progress \
		--noshow_loading_progress \
		--features=race //... 2>&1 | ts '[%Y-%m-%d %H:%M:%.S]'
