# command alias
GOGET ?= go get -u -v
GOBUILD ?= go build
GOFMT ?= gofmt -s
GOTEST ?= go test

# variables
GOFILES ?= $(shell find . -name "*.go")
PKGS ?= $(shell go list ./...)
VERSION ?= $(shell git describe --tag --abbrev=0)
REVISION ?= $(shell git rev-parse --short HEAD)
# https://stackoverflow.com/questions/18136918/how-to-get-current-relative-directory-of-your-makefile
MAKEFILE_PATH ?= $(abspath $(lastword $(MAKEFILE_LIST)))
REPO_NAME ?= $(notdir $(patsubst %/,%,$(dir $(MAKEFILE_PATH))))
LDFLAGS ?= -ldflags="-s -w \
		-X 'main.version=$(VERSION)' \
		-X 'main.revision=$(REVISION)' \
	"
OUT_DIR ?= outputs
BIN_PATH ?= $(OUT_DIR)/$(REPO_NAME)

DOCKER_VOLUME_NAME ?= graphql-server-go
DOCKER_NETWORK_NAME ?= graphql-server-go
DB_OPTION ?= "--build"

# https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.DEFAULT_GOAL := help

.PHONY: init
init: ## initialize repository
	cp .env.sample .env
	docker volume create $(DOCKER_VOLUME_NAME)
	docker network create $(DOCKER_NETWORK_NAME)

.PHONY: clean
clean: ## clean up repository
	docker volume rm $(DOCKER_VOLUME_NAME)
	docker network rm $(DOCKER_NETWORK_NAME)

.PHONY: db
db: ## start db service(Usage: make db DB_OPTION="--build -d")
	docker-compose up $(DB_OPTION) postgres

.PHONY: generate
generate: ## generate codes
	gqlgen generate

.PHONY: install
install: ## install
	$(GOGET) \
		github.com/99designs/gqlgen \
		github.com/xo/xo \
		golang.org/x/tools/gopls \
		golang.org/x/lint/golint

.PHONY: tidy
tidy: ## tidy
	go mod tidy

.PHONY: format
format: ## format codes
	$(GOFMT) -w $(GOFILES)

.PHONY: lint
lint: ## lint codes
	for PKG in $(PKGS); do golint -set_exit_status $$PKG || exit 1; done;

.PHONY: vet
vet: ## vet
	for PKG in $(PKGS); do go vet $$PKG || exit 1; done;

.PHONY: test
test: ## run tests
	$(GOTEST) -cover -v ./...

.PHONY: build
build: ## build an app
	$(GOBUILD) $(LDFLAGS) -o $(BIN_PATH) .

# https://qiita.com/dtan4/items/8c417b629b6b2033a541
.PHONY: crossbuild
crossbuild: ## cross build
	mkdir -p $(OUT_DIR)/release
	for os in darwin linux windows; do \
		for arch in amd64; do \
			make build GOBUILD="GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build" BIN_PATH="$(OUT_DIR)/$$os-$$arch/$(REPO_NAME)"; \
			zip $(OUT_DIR)/release/$(REPO_NAME)_$(VERSION)_$$os-$$arch.zip $(OUT_DIR)/$$os-$$arch/$(REPO_NAME); \
		done; \
	done
	# Raspberry Pi
	make build \
		GOBUILD='GOOS=linux GOARCH=arm GOARM=7 go build' \
		BIN_PATH="$(OUT_DIR)/linux-arm/$(REPO_NAME)"
	zip $(OUT_DIR)/release/$(REPO_NAME)_$(VERSION)_linux-arm.zip $(OUT_DIR)/linux-arm/$(REPO_NAME)

.PHONY: ci
ci: install format lint vet test crossbuild build ## run ci tests
