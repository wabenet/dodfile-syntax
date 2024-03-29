all: clean build

.PHONY: clean
clean:
	rm -f dodfile-syntax

.PHONY: update
update:
	go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all | xargs --no-run-if-empty go get
	go mod tidy

.PHONY: lint
lint:
	golangci-lint run --enable-all

.PHONY: test
test:
	go test -cover -race ./...

.PHONY: build
build:
	go build -o dodfile-syntax .

.PHONY: image
image:
	docker build -t wabenet/dodfile-syntax .
