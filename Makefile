# Some things this makefile could make use of:
#
# - test coverage target(s)
# - profiler target(s)
#

BIN            = kit-overwatch
OUTPUT_DIR     = build
TMP_DIR       := .tmp
RELEASE_VER   := $(shell git rev-parse --short HEAD)

TEST_PACKAGES      := $(shell find . -name '*_test.go' -exec dirname {} \; | uniq)
TEST_UNIT_PACKAGES := $(shell find . -name '*_unit_test.go' -exec dirname {} \; | uniq)
TEST_INT_PACKAGES  := $(shell find . -name '*_int_test.go' -exec dirname {} \; | uniq)

.PHONY: help
.DEFAULT_GOAL := help

run: ## Run via godep (without building)
	go run *.go

all: test build docker ## Test, build and docker image build
	
fmt: ## Run go fmt for all files in the project (excluding vendor)
	go fmt $(shell go list ./... | grep -v /vendor/)

test: test/fmt test/unit test/integration ## Perform both unit and integration tests

test/fmt: ## Check if all files (excluding vendor) conform to fmt
	test -z $(shell echo $(shell go fmt $(shell go list ./... | grep -v /vendor/)) | tr -d "[:space:]")

test/unit: ## Perform unit tests
	go test -v -cover -tags unit $(TEST_UNIT_PACKAGES)

test/integration: ## Perform integration tests
	go test -v -cover -tags integration $(TEST_INT_PACKAGES)

test/race: ## Perform unit and integration tests and enable the race detector
	go test -v -tags "unit integration" -race $(TEST_PACKAGES)

build: clean build/linux build/darwin ## Build for linux and darwin (save to OUTPUT_DIR/BIN)

build/linux: clean/linux ## Build for linux (save to OUTPUT_DIR/BIN)
	GOOS=linux go build -a -installsuffix cgo -ldflags "-X main.version=$(RELEASE_VER)" -o $(OUTPUT_DIR)/$(BIN)-linux .

build/darwin: clean/darwin ## Build for darwin (save to OUTPUT_DIR/BIN)
	GOOS=darwin go build -a -installsuffix cgo -ldflags "-X main.version=$(RELEASE_VER)" -o $(OUTPUT_DIR)/$(BIN)-darwin .

docker: build/linux ## Build local docker image
	docker build -t $(BIN):$(RELEASE_VER) .

jet: ## Run `jet steps`
	jet steps

clean: clean/darwin clean/linux ## Remove all build artifacts

clean/darwin: ## Remove darwin build artifacts
	$(RM) $(OUTPUT_DIR)/$(BIN)-darwin

clean/linux: ## Remove linux build artifacts
	$(RM) $(OUTPUT_DIR)/$(BIN)-linux

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
