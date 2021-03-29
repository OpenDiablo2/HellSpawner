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

.PHONY: all build setup test cover lint clean run help

## all: Default targe	t, now is build
all: build

## build: Builds the binary
build:
	@echo "Building..."
	@$(GOCMD) build -o ${NAME}

## setup: Runs mod download and generate
setup:
	@echo "Downloading tools and dependencies..."
	@git submodule update --init --recursive
	@$(GOCMD) get -u
	@$(GOCMD) mod tidy
	@$(GOCMD) get -d
	@$(GOCMD) mod download -x
	@$(GOCMD) generate -v ./...

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
	@golangci-lint run ./...
	@misspell ./...

## clean: Runs go run
clean:
	@echo "Cleaning..."
	@$(GOCMD) clean

## run: Runs go run
run: build
	@$(GOCMD) run ${TARGET}

## help: Prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
