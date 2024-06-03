.PHONY: complete logging metrics tracing help

default: help

help:
	@echo "Usage: make [command]"
	@echo "Commands:"
	@echo "  - complete   executes an example will all three (logging, metrics, and tracing) enabled."
	@echo "  - logging    demonstrates how to utilize the logging package."
	@echo "  - metrics    demonstrates how to utilize the metrics package."
	@echo "  - tracing    demonstrates how to utilize the tracing package."

complete:
	go run ./complete/main.go

logging:
	go run ./logging/main.go

metrics:
	go run ./metrics/main.go

tracing:
	go run ./tracing/main.go