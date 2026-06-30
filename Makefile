.phony: lint format server build db-start check deps db-shell docker-build docker-run db-create test test-verbose test-one

include .env
export

DB_DSN=host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable

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

#apply all pending migrations
migrate-up:
	goose -dir internal/db/migrations postgres "$(DB_DSN)" up

#roll back the last migration
migrate-down:
	goose -dir internal/db/migrations postgres "$(DB_DSN)" down

#show migration status
migrate-status:
	goose -dir internal/db/migrations postgres "$(DB_DSN)" status

#create the database container (first time only)
db-create:
	podman run -d \
		--name tahrir-pg \
		-e POSTGRES_USER=tahrir \
		-e POSTGRES_PASSWORD=tahrir \
		-e POSTGRES_DB=tahrir_go \
		-p 5432:5432 \
		postgres:16

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


db-shell:
	-podman exec -it tahrir-pg psql -U tahrir -d tahrir_go || echo "Failed to connect to database"

#run all tests
test:
	go test ./...

#run all tests with verbose output
test-verbose:
	go test ./... -v

#run a specific test by name, e.g. make test-one NAME=TestGetPersonByID
test-one:
	go test ./... -run $(NAME) -v
