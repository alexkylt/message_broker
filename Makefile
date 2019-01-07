SERVER_BIN := "server.bin"
CLIENT_BIN := "client.bin"

BIN_DIR := $(GOPATH)/bin
# BUILD_DIR=${GOPATH}/src/github.com/${GITHUB_USERNAME}
PKG := "github.com/alexkylt/message_broker"
# Packages lists
PACKAGES=$(shell go list ./...)
GOFLAGS ?= $(GOFLAGS:)
SERVER_PKG_BUILD := "${PKG}"
CLIENT_PKG_BUILD := "${PKG}/internal/pkg/client"
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

.PHONY: all build_server build_client

all: build_server build_client docker_network docker_build_server docker_build_client docker_build_psql docker_run_psql docker_run_server docker_run_client

gobuild: ## 
	

dep: ## Get the dependencies
	@go get -v -d ./...

govet: ## Runs govet against all packages.
	@echo Running GOVET
	@go vet ./... || exit 1

goimports: ## Runs goimports against all packages.
	@echo Running goimports

	@for package in $(PACKAGES); do \
		echo "Checking "$$package; \
		files=$$(go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}' $$package); \
		if [ "$$files" ]; then \
			goimports_output=$$(goimports -l -d $$files 2>&1); \
			if [ "$$goimports_output" ]; then \
				echo "$$goimports_output"; \
				echo "goimports failure"; \
				exit 1; \
			fi; \
		fi; \
	done
	@echo "goimports success"; \
# TODO FIXME!!!!
golint: ## Runs golint against all packages.
	@echo Running GOLINT
	@golint ./... || exit 1

build_server: ## Build the binary file for server
	@go build -i -v -o $(BIN_DIR)/$(SERVER_BIN) $(SERVER_PKG_BUILD)

build_client: ## Build the binary file for client
	@go build -i -v -o $(BIN_DIR)/$(CLIENT_BIN) $(CLIENT_PKG_BUILD)

docker_network: ##
	@if [ $(shell docker network ls --format '{{.Name}}'| grep $(NETWORK)| wc -l) -eq 0 ]; then \
		docker network create $(NETWORK); \
	fi

docker_build_server : docker_network ## Build default docker image
	@docker build -t "$(DOCKER_BUILD_SERVER)" -f Dockerfile.server .

docker_build_client : docker_network ## Build default docker image
	@docker build -t "$(DOCKER_BUILD_CLIENT)" -f Dockerfile.client .

docker_build_psql : docker_network ## Build default docker image
	@docker build -t "$(DOCKER_BUILD_PSQL)" -f Dockerfile.psql .

docker_run_psql: ## Run default docker image
	@if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/$(DOCKER_IMAGE_PSQL)$$ | wc -l) -eq 0 ]; then \
		echo starting $(DOCKER_IMAGE_PSQL); \
		docker run --network $(NETWORK) --name=$(DOCKER_IMAGE_PSQL) -d -p 5433:5432 $(DOCKER_BUILD_PSQL); \
	fi
#   		--network-alias $(DOCKER_IMAGE_SERVER).$(NETWORK)
docker_run_server: ## Run default docker image
	@if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/$(DOCKER_IMAGE_SERVER)$$ | wc -l) -eq 0 ]; then \
		echo starting $(DOCKER_IMAGE_SERVER); \
		docker run --network $(NETWORK) --name=$(DOCKER_IMAGE_SERVER) -d -p $(SERVER_PORT):$(SERVER_PORT) $(DOCKER_BUILD_SERVER) --port=$(SERVER_PORT) --mode=$(STORAGE_MODE) > /dev/null; \
	fi

docker_run_client: ## Run default docker image
	@if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/$(DOCKER_IMAGE_CLIENT)$$ | wc -l) -eq 0 ]; then \
		echo starting $(DOCKER_IMAGE_CLIENT); \
		docker run --rm --network $(NETWORK) --name=$(DOCKER_IMAGE_CLIENT) -it $(DOCKER_BUILD_CLIENT) --port=$(SERVER_PORT) \
		--host=$(SERVER_HOST); \
	fi

clean: ## Remove previous builds
	@rm $(SERVER_OUT) $(CLIENT_OUT) $(API_OUT)

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

