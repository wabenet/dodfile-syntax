all: clean build

.PHONY: clean
clean:
	rm -f dodfile-syntax

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: lint
lint:
	golangci-lint run --enable-all

.PHONY: build
build:
	go build -o dodfile-syntax .
