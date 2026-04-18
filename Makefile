.PHONY: run build tidy sqlc test

run:
	go run ./cmd/mahora

build:
	go build -o bin/mahora ./cmd/mahora

tidy:
	go mod tidy

sqlc:
	sqlc generate -f internal/db/sqlc.yaml

test:
	go test ./...

lint:
	golangci-lint run ./...
