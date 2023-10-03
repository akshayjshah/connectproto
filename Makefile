# See https://tech.davis-hansson.com/p/make/
SHELL := bash
.DELETE_ON_ERROR:
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := all
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-print-directory
BIN := .tmp/bin
export PATH := $(BIN):$(PATH)
export GOBIN := $(abspath $(BIN))

.PHONY: help
help: ## Describe useful make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "%-30s %s\n", $$1, $$2}'

.PHONY: all
all: ## Build, test, and lint (default)
	$(MAKE) test
	$(MAKE) lint

.PHONY: test
test: build ## Run unit tests
	go test -vet=off -race -cover ./...

.PHONY: build
build: ## Build all packages
	go build ./...

.PHONY: lint
lint: $(BIN)/gofmt $(BIN)/staticcheck ## Lint Go
	test -z "$$(gofmt -s -l . | tee /dev/stderr)"
	go vet ./...
	staticcheck ./...

.PHONY: lintfix
lintfix: $(BIN)/gofmt ## Automatically fix some lint errors
	gofmt -s -w .

.PHONY: upgrade
upgrade: ## Upgrade dependencies
	go get -u -t ./... && go mod tidy -v

.PHONY: clean
clean: ## Remove intermediate artifacts
	rm -rf .tmp

$(BIN)/gofmt:
	@mkdir -p $(@D)
	go build -o $(@) cmd/gofmt

$(BIN)/staticcheck:
	@mkdir -p $(@D)
	go install honnef.co/go/tools/cmd/staticcheck@latest
