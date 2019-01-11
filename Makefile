SERVER_BIN := "server.bin"
CLIENT_BIN := "client.bin"

CURRENT_DIR := $(shell pwd)
BIN_DIR := $(GOPATH)/bin
PKG := "github.com/alexkylt/message_broker"
# Packages lists
PACKAGES=$(shell go list ./...)
GOFLAGS ?= $(GOFLAGS:)
SERVER_PKG_BUILD := "${PKG}/cmd/server"
CLIENT_PKG_BUILD := "${PKG}/cmd/client"
BINARIES := "cmds"
# PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)

DOCKER_TAG = v5
DOCKER_IMAGE_SERVER = server
DOCKER_IMAGE_CLIENT = client
DOCKER_IMAGE_PSQL = postgres
DOCKER_BUILD_SERVER := $(DOCKER_IMAGE_SERVER):$(DOCKER_TAG)
DOCKER_BUILD_CLIENT := $(DOCKER_IMAGE_CLIENT):$(DOCKER_TAG)
DOCKER_BUILD_PSQL := $(DOCKER_IMAGE_PSQL):$(DOCKER_TAG)
NETWORK := dockernet

STORAGE_MODE ?= "db"
SERVER_PORT ?= 9090
SERVER_HOST ?= $(DOCKER_IMAGE_SERVER)

DOCKER_BUILD_BUILDER := builder
BINARIES := $(CURRENT_DIR)/cmds
# GOARCH = amd64
goarch = $(shell go env GOARCH)
goos = $(shell go env GOOS)

.PHONY: all build_server build_client

all: build_docker_builder check docker_network build_server build_client docker_server docker_client docker_psql docker_run_psql docker_run_server docker_run_client
# golint 
check: build_docker_builder goimports govet



init:
	mkdir $(BINARIES)

build_docker_builder: init ## Build default docker image
	@docker build -t "$(DOCKER_BUILD_BUILDER)" -f Dockerfile.build .
	# @if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/$(DOCKER_BUILD_BUILDER)$$ | wc -l) -eq 0 ]; then \
	# fi

goimports: ## build_docker_builder Runs goimports against all packages.
	@echo Running goimports

	@for package in $(PACKAGES); do \
		echo "Checking "$$package; \
		files=$$(docker run $(DOCKER_BUILD_BUILDER) list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}' $$package); \
		if [ "$$files" ]; then \
			goimports_output=$$(docker run --entrypoint="goimports" $(DOCKER_BUILD_BUILDER) -l -d $$files 2>&1); \
			if [ "$$goimports_output" ]; then \
				echo "$$goimports_output"; \
				echo "goimports failure"; \
				exit 1; \
			fi; \
		fi; \
	done
	@echo "goimports success"; \

govet: ## build_docker_builder Runs govet against all packages.
	@echo Running GOVET
	@docker run -e CGO_ENABLED=0 $(DOCKER_BUILD_BUILDER) vet ./... || exit 1
	@echo "GOVET success";

# TODO FIXME!!!!
# golint: ## build_docker_builder Runs golint against all packages.
# 	@echo Running GOLINT
# 	@docker run --entrypoint="golint" $(DOCKER_BUILD_BUILDER) -set_exit_status ./... || exit 1
# 	
# 	@echo "GOLINT success";

docker_network: ##
	@if [ $(shell docker network ls --format '{{.Name}}'| grep $(NETWORK)| wc -l) -eq 0 ]; then \
		@docker network create $(NETWORK); \
	fi

# -e GOOS=$(goos) -e GOARCH=$(goarch) -e CGO_ENABLED=0
#--user $(id -u):$(id -g)1
build_server: ## build_docker_builder Build the binary file for server	
	@docker run -v $(BINARIES):$(BIN_DIR) $(DOCKER_BUILD_BUILDER) build -i -v -o $(BIN_DIR)/$(SERVER_BIN) $(SERVER_PKG_BUILD)

build_client: ## build_docker_builder Build the binary file for server
	@docker run -u $(id -u):$(id -g) -v $(BINARIES):$(BIN_DIR) $(DOCKER_BUILD_BUILDER) build -i -v -o $(BIN_DIR)/$(CLIENT_BIN) $(CLIENT_PKG_BUILD)

docker_run_psql: docker_network ## Run default docker image
	@if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/$(DOCKER_IMAGE_PSQL)$$ | wc -l) -eq 0 ]; then \
		echo starting $(DOCKER_IMAGE_PSQL); \
		docker run --network $(NETWORK) --name=$(DOCKER_IMAGE_PSQL) -d -p 5433:5432 $(DOCKER_BUILD_PSQL); \
	fi

docker_server: docker_network build_server ## Build default docker image
	@docker build -t "$(DOCKER_BUILD_SERVER)" -f Dockerfile.server .

docker_client: docker_network build_client ## Build default docker image
	@docker build -t "$(DOCKER_BUILD_CLIENT)" -f Dockerfile.client .

docker_psql: docker_network ## Build default docker image
	@docker build -t "$(DOCKER_BUILD_PSQL)" -f Dockerfile.psql .

docker_run_server: docker_network ## Run default docker image
	@if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/$(DOCKER_IMAGE_SERVER)$$ | wc -l) -eq 0 ]; then \
		echo starting $(DOCKER_IMAGE_SERVER); \
		docker run --network $(NETWORK) --name=$(DOCKER_IMAGE_SERVER) -d -p $(SERVER_PORT):$(SERVER_PORT) $(DOCKER_BUILD_SERVER) --port=$(SERVER_PORT) --mode=$(STORAGE_MODE) > /dev/null; \
	fi

docker_run_client: docker_network ## Run default docker image
	@if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/$(DOCKER_IMAGE_CLIENT)$$ | wc -l) -eq 0 ]; then \
		echo starting $(DOCKER_IMAGE_CLIENT); \
		docker run --rm --network $(NETWORK) --name=$(DOCKER_IMAGE_CLIENT) -it $(DOCKER_BUILD_CLIENT) --port=$(SERVER_PORT) \
		--host=$(SERVER_HOST); \
	fi

# @docker stop $(shell docker ps -a -q)
# @docker rm $(shell docker ps -a -q)
# @docker image prune --all
clean: ## Remove previous builds
	@rm -r -f $(BINARIES)

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

