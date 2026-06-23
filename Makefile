lint:
	golangci-lint run

format:
	gofmt -s -w .

server:
	air

build:
	go build -o bin/tahrir cmd/server/main.go

db-start:
	podman start tahrir-pg

check: format lint