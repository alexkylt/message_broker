SERVER_OUT := "server.bin"
CLIENT_OUT := "client.bin"
# API_OUT := "api/api.pb.go"
PKG := "github.com/alexkylt/message_broker"
# Packages lists
PACKAGES=$(shell go list ./...)
GOFLAGS ?= $(GOFLAGS:)
SERVER_PKG_BUILD := "${PKG}/internal"
CLIENT_PKG_BUILD := "${PKG}/internal/pkg/client"
# PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)

.PHONY: all build_server build_client

all: build_server build_client

dep: ## Get the dependencies
	@go get -v -d ./...

build_server: ## Build the binary file for server
	@go build -i -v -o $(SERVER_OUT) $(SERVER_PKG_BUILD)

build_client: ## Build the binary file for client
	@go build -i -v -o $(CLIENT_OUT) $(CLIENT_PKG_BUILD)




















gofmt: ## Runs gofmt against all packages.
	@echo Running GOFMT

	@for package in $(TE_PACKAGES); do \
		echo "Checking "$$package; \
		files=$$(go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}' $$package); \
		if [ "$$files" ]; then \
			gofmt_output=$$(gofmt -d -s $$files 2>&1); \
			if [ "$$gofmt_output" ]; then \
				echo "$$gofmt_output"; \
				echo "gofmt failure"; \
				exit 1; \
			fi; \
		fi; \
	done
	@echo "gofmt success"; \

govet: ## Runs govet against all packages.
	@echo Running GOVET
	@go vet $(GOFLAGS) $(PACKAGES) || exit 1

goimports: ## Runs goimports against all packages.
	@echo Running GOIMPORTS
	@goimports $(GOFLAGS) $(PACKAGES) || exit 1

clean: ## Remove previous builds
	@rm $(SERVER_OUT) $(CLIENT_OUT) $(API_OUT)

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

