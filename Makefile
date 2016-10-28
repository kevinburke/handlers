vet:
	go vet ./...

test: vet
	go test ./...

race-test: vet
	go test -race ./...

install:
	go install ./...

release: race-test
	bump_version minor lib.go
