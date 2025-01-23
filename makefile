pkgdoc:
	pkgsite -open .
.PHONY: pkgdoc

format:
	go fmt ./...
.PHONY: format

test:
	go test ./... -v -cover
.PHONY: test

build:
	go build -v -o number-guessing ./cmd/main.go
.PHONY: build