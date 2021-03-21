DOCKER_VOLUME_NAME ?= graphql-server-go
DOCKER_NETWORK_NAME ?= graphql-server-go

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
db: ## start db service
	docker-compose up --build -d postgres
