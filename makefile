.PHONY: build test

default: test

build:
	go generate ./...

test: build
	go clean -testcache
	go test ./attributes ./logging ./tracing ./metrics -coverprofile=coverage.out
	go tool cover -html=coverage.out