.PHONY: build
build:
	go mod download
	go build -o toimi-rest ./cmd/main.go

.PHONY: run
run:
	go mod download
	go build -o toimi-rest ./cmd/main.go
	./toimi-rest

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build