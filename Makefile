build:
	go build -o bin/gendiff ./cmd/gendiff

run: build
	./bin/gendiff

test:
	go test -v ./...

test-coverage:
	go test -coverpkg=./... -coverprofile=coverage.out ./...

lint:
	golangci-lint run ./...

.PHONY: build run test test-coverage lint
