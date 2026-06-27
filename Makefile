.phony: lint format server build db-start check deps

#check all go files for linting issues
lint:
	golangci-lint run

#check and format all go files
format:
	gofmt -s -w .

#start the server with live reload
server:
	air

#build the project
build:
	go build ./...

#start the database container
db-start:
	podman start tahrir-pg

#check all code quality checks
checks: format lint

#download dependencies
deps:
	go mod tidy
	go mod download

docker-build:
	docker build -t tahrir-go .

docker-run:
	docker run --network host --env-file .env tahrir-go