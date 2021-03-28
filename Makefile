NAME=HellSpawner
GOCMD=LC_ALL=C go
TIMEOUT=5

# go tools
export PATH := ./bin:$(PATH)
export GO111MODULE := on
export GOPROXY = https://proxy.golang.org,direct

# go source files
SRC = $(shell find . -type f -name "*.go")
# The name of the executable (default is current directory name)
TARGET := $(shell echo $${PWD-`pwd`})

.PHONY: all build setup fmt test cover lint clean run help

## all: Default target, now is build
all: build

## build: Builds the binary
build: fmt
	@echo "Building..."
	@$(GOCMD) build -o ${NAME}

## setup: Runs mod download and generate
setup:
	@echo "Downloading tools and dependencies..."
	@$(GOCMD) mod download -x
	@$(GOCMD) generate -v ./...

## fmt: Runs go goimports and gofmt
fmt: setup
	@echo "Checking the imports..."
	@$(GOCMD)imports -w ${SRC}
	@echo "Formatting the go files..."
	@$(GOCMD)fmt -w -s ${SRC}

## test: Runs the tests with coverage
test:
	@echo "Running tests..."
	@$(GOCMD) test -failfast -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt ./... -run . -timeout $(TIMEOUT)m

## cover: Runs all tests and opens the coverage report in the browser
cover: test
	@$(GOCMD) tool cover -html=coverage.txt

## lint: Runs golangci-lint (configuration at .golangci.yml) and misspell
lint: setup
	@echo "Running linters..."
	@golangci-lint run --max-issues-per-linter 0 --max-same-issues 0 ./...
	@misspell ./...

## run: Runs go run
run: build
	@$(GOCMD) run ${TARGET}

## help: Prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
