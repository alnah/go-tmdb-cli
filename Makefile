BIN=bin/
DIST=dist/
SRC=$(shell find . -name "*.go")
TARGET=$(BIN)/go-tmdb-cli

ifeq (, $(shell which golangci-lint))
	$(warning "could not find golangci-lint in $(PATH), \
	run: curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh")
endif

ifeq (, $(shell which goreleaser))
	$(warning "could not find goreleaser in $(PATH), \
	run: go install github.com/goreleaser/goreleaser/v2@latest")
endif

.PHONY: fmt lint test install build clean

default: build

all: install fmt lint test benchmark build

install:
	$(info 📥 DOWNLOADING DEPENDENCIES...)
	go get -v ./...

fmt:
	$(info ✨ CHECKING CODE FORMATTING...)
	@test -z $(shell gofmt -l $(SRC)) || (gofmt -d $(SRC); exit 1)

lint:
	$(info 🔍 RUNNING LINT TOOLS...)
	golangci-lint run --config .golangci.yaml

test: install
	$(info 🧪 RUNNING TESTS...)
	go test -v ./... -cover

benchmark: install
	$(info 🚀 RUNNING BENCHMARKS...)
	go test -bench=.

build: install
	$(info 🏗️ BUILDING THE PROJECT...)
	@if [ -e "$(TARGET)" ]; then rm -rf "$(TARGET)"; fi
	@mkdir -p $(BIN)
	@go build -o $(TARGET)

release: fmt lint test benchmark
	$(info 📦 CREATING A NEW RELEASE...)
	goreleaser release

clean:
	$(info 🧹 CLEANING UP...)
	rm -rf $(BIN)
	rm -rf $(DIST)
