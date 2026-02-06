SHELL := /bin/sh

APP := prettycat
CMD := ./cmd/prettycat
BIN_DIR := bin
BIN := $(BIN_DIR)/$(APP)
ARGS ?=

.PHONY: help build run test testv fmt fmt-check tidy check clean install

help: ## Show available commands
	@awk 'BEGIN {FS = ":.*## "}; /^[a-zA-Z0-9_-]+:.*## / {printf "\033[36m%-12s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the CLI binary into ./bin
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN) $(CMD)

run: ## Run the CLI (use ARGS="file.md --no-color")
	go run $(CMD) $(ARGS)

test: ## Run test suite
	go test ./...

testv: ## Run test suite (verbose)
	go test -v ./...

fmt: ## Format all Go files
	gofmt -w $$(find . -type f -name '*.go' -not -path './.cache/*')

fmt-check: ## Check formatting without writing changes
	@test -z "$$(gofmt -l $$(find . -type f -name '*.go' -not -path './.cache/*'))" || \
		(echo "Go files need formatting. Run: make fmt" && exit 1)

tidy: ## Sync and clean module dependencies
	go mod tidy

check: fmt-check test build ## Run quality gate (format check, tests, build)

install: ## Install prettycat into GOPATH/bin
	go install $(CMD)

clean: ## Remove build artifacts
	rm -rf $(BIN_DIR)
