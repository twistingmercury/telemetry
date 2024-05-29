.PHONY: test

default: test

test:
	go clean -testcache
	go test ./logging ./metrics ./tracing -v -coverprofile=coverage.out
	go tool cover -html=coverage.out
