.phony: lint format server build db-start check

#check all go files for linting issues
lint:
	golangci-lint run

#check and format all go files
format:
	gofmt -s -w .

#start the server with live reload
server:
	air

#build the server binary
build:
	go build -o bin/tahrir cmd/server/main.go

#start the database container
db-start:
	podman start tahrir-pg

#check all code quality checks
check: format lint

#download dependencies
deps:
	go mod tidy
	go mod download