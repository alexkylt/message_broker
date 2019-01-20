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
DOCKERFILE_PSQL := "Dockerfile.psql"
DOCKERFILE_SERVER := "Dockerfile.server"
DOCKERFILE_CLIENT := "Dockerfile.client"
DOCKERFILE_BUILDER := "Dockerfile.build"

.PHONY: all build_server build_client

all: clean build_docker_builder check docker_network build_server build_client docker_server docker_client docker_psql run_psql run_server run_client

check: build_docker_builder goimports govet golint



init:
	@if [ -d "$(BINARIES)" ]; then  \
		if [ -L "$(BINARIES)" ]; then \
			rm -f $(BINARIES); \
			mkdir $(BINARIES); \
		else \
			rm -r -f $(BINARIES); \
			mkdir $(BINARIES); \
		fi \
	else \
		mkdir $(BINARIES); \
	fi \
	

build_docker_builder: init ## Build builder image
	@docker build -t "$(DOCKER_BUILD_BUILDER)" -f $(DOCKERFILE_BUILDER) .

goimports: build_docker_builder ## Run goimports
	@echo Running goimports

	@for package in $(PACKAGES); do \
		echo "Checking "$$package; \
		files=$$(docker run $(DOCKER_BUILD_BUILDER) list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}' $$package); \
		if [ "$$files" ]; then \
			echo $$files; \
			goimports_output=$$(docker run --entrypoint="goimports" $(DOCKER_BUILD_BUILDER) -l -d $$files 2>&1); \
			if [ "$$goimports_output" ]; then \
				echo "$$goimports_output"; \
				echo "goimports failure"; \
				exit 1; \
			fi; \
		fi; \
	done
	@echo "goimports success"; \

govet: build_docker_builder ## Run govet
	@echo Running GOVET
	@docker run -e CGO_ENABLED=0 $(DOCKER_BUILD_BUILDER) vet ./... || exit 1
	@echo "GOVET success";

golint: build_docker_builder ## Run golint
	@echo Running GOLINT
	@docker run --entrypoint="golint" $(DOCKER_BUILD_BUILDER) -set_exit_status ./... || exit 1
	
	@echo "GOLINT success";

docker_network: ## Create docker network
	@echo "START CREATE DOCKER NETWORK" 
	@if [ $(shell docker network ls --format '{{.Name}}'| grep $(NETWORK)| wc -l) -eq 0 ]; then \
		docker network create $(NETWORK); \
		echo finish; \
	else \
		for i in $(shell docker network inspect -f '{{range .Containers}}{{.Name}} {{end}}' $(NETWORK));\
			do \
				docker network disconnect -f $(NETWORK) $$i; \
			done; \
		docker network rm $(NETWORK); \
		docker network create $(NETWORK); \
		echo finish; \
	fi
	@echo "END CREATE DOCKER NETWOK" 

build_server: init build_docker_builder ## Build server 
	@echo "START BUILD SERVER" 	
	@docker run -e CGO_ENABLED=0 -v $(BINARIES):$(BIN_DIR) $(DOCKER_BUILD_BUILDER) build -i -v -o $(BIN_DIR)/$(SERVER_BIN) $(SERVER_PKG_BUILD)
	@echo "END BUILD SERVER" 

build_client: init build_docker_builder ## Build client 
	@echo "START BUILD CLIENT" 	
	@docker run -e CGO_ENABLED=0 -v $(BINARIES):$(BIN_DIR) $(DOCKER_BUILD_BUILDER) build -i -v -o $(BIN_DIR)/$(CLIENT_BIN) $(CLIENT_PKG_BUILD)
	@echo "END BUILD CLIENT" 

build: init build_docker_builder build_server build_client

run_psql: docker_network ## Run PSQL
	@if [ $(STORAGE_MODE) = "db" ]; then \
		if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/$(DOCKER_IMAGE_PSQL)$$ | wc -l) -eq 0 ]; then \
			echo "START POSTGRES DATABASE" $(DOCKER_IMAGE_PSQL); \
			docker run --network $(NETWORK) --name=$(DOCKER_IMAGE_PSQL) -d -p 5433:5432 $(DOCKER_BUILD_PSQL); \
			echo finish; \
		else \
			docker container rm -f $(DOCKER_IMAGE_PSQL); \
			docker run --network $(NETWORK) --name=$(DOCKER_IMAGE_PSQL) -d -p 5433:5432 $(DOCKER_BUILD_PSQL); \
			echo finish; \
		fi \
	fi

docker_server: docker_network build_server ## Build server image
	@echo "START BUILD SERVER DOCKER" 	
	@docker build -t "$(DOCKER_BUILD_SERVER)" -f $(DOCKERFILE_SERVER) .
	@echo "END BUILD SERVER DOCKER" 	

docker_client: docker_network build_client ## Build client image
	@echo "START BUILD CLIENT DOCKER" 
	@docker build -t "$(DOCKER_BUILD_CLIENT)" -f $(DOCKERFILE_CLIENT) .
	@echo "END BUILD CLIENT DOCKER" 

docker_psql: docker_network ## Build PSQL image
	@if [ $(STORAGE_MODE) = "db" ]; then \
		echo "START BUILD PSQL DOCKER" ; \
		docker build -t "$(DOCKER_BUILD_PSQL)" -f $(DOCKERFILE_PSQL) . ; \
		echo "END BUILD PSQL DOCKER" ; \
	fi

run_server: docker_network run_psql ## Run server
	@if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/$(DOCKER_IMAGE_SERVER)$$ | wc -l) -eq 0 ]; then \
		echo starting $(DOCKER_IMAGE_SERVER); \
		docker run --network $(NETWORK) --name=$(DOCKER_IMAGE_SERVER) -d -p $(SERVER_PORT):$(SERVER_PORT) $(DOCKER_BUILD_SERVER) --port=$(SERVER_PORT) --mode=$(STORAGE_MODE) > /dev/null; \
		echo finish; \
	else \
		echo starting $(DOCKER_IMAGE_SERVER) plus; \
		docker container rm -f $(DOCKER_IMAGE_SERVER); \
		docker run --network $(NETWORK) --name=$(DOCKER_IMAGE_SERVER) -d -p $(SERVER_PORT):$(SERVER_PORT) $(DOCKER_BUILD_SERVER) --port=$(SERVER_PORT) --mode=$(STORAGE_MODE) > /dev/null; \
		echo finish; \
	fi

run_client: run_server ## Run client
	@if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/$(DOCKER_IMAGE_CLIENT)$$ | wc -l) -eq 0 ]; then \
		echo starting $(DOCKER_IMAGE_CLIENT); \
		docker run --rm --network $(NETWORK) --name=$(DOCKER_IMAGE_CLIENT) -it $(DOCKER_BUILD_CLIENT) --port=$(SERVER_PORT) \
		--host=$(SERVER_HOST); \
	fi

clean: ## Remove previous builds
	@rm -r -f $(BINARIES)
	@if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/$(DOCKER_IMAGE_SERVER)$$ | wc -l) -eq 1 ]; then \
		docker container rm -f $(DOCKER_IMAGE_SERVER); \
	fi
	@if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/$(DOCKER_IMAGE_PSQL)$$ | wc -l) -eq 1 ]; then \
		docker container rm -f $(DOCKER_IMAGE_PSQL); \
	fi

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
