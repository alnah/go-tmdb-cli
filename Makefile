BIN=./bin
SRC=$(shell find . -name "*.go")
TARGET=$(BIN)/go-tmdb-cli

ifeq (, $(shell which golangci-lint))
$(warning "could not find golangci-lint in $(PATH), run: curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh")
endif

.PHONY: fmt lint test install build clean

default: build

all: install fmt lint test benchmark build

install:
	$(info ******************** downloading dependencies ***************)
	go get -v ./...

fmt:
	$(info ******************** checking formatting ********************)
	@test -z $(shell gofmt -l $(SRC)) || (gofmt -d $(SRC); exit 1)

lint:
	$(info ******************** running lint tools *********************)
	golangci-lint run --config .golangci.yaml

test: install
	$(info ******************** running tests **************************)
	go test -v ./... -cover

benchmark: install
	$(info ******************** running benchmarks *********************)
	go test -bench=.

build: install
	$(info ******************** building go-tmdb-cli *******************)
	@if [ -e "$(TARGET)" ]; then rm -rf "$(TARGET)"; fi
	@mkdir -p $(BIN)
	@go build -o $(TARGET)
	
clean:
	rm -rf $(BIN)