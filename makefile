.PHONY: test

default: test

test:
	go clean -testcache
	go test ./attributes ./logging ./tracing ./metrics -v covermode=count -coverprofile=coverage.out
	go tool cover -html=coverage.out