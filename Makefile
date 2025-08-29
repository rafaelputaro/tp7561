SHELL := /bin/bash

default: docker-image

deps:
	cd src && go mod tidy
	cd src && go mod vendor
.PHONY: deps

docker-compose-dev.yaml: compose-generator.py config.ini
	python3 compose-generator.py
.PHONY: docker-compose-dev.yaml

docker-image: deps
	docker build -f ./src/server/workers/filter/Dockerfile -t "peer:latest" .
.PHONY: docker-image

docker-compose-up: docker-compose-dev.yaml docker-image
	docker compose -f docker-compose-dev.yaml down --remove-orphans || true
	docker compose -f docker-compose-dev.yaml up -d --build --force-recreate
.PHONY: docker-compose-up

docker-compose-logs:
	docker compose -f docker-compose-dev.yaml logs -f
.PHONY: docker-compose-logs

docker-compose-down:
	docker compose -f docker-compose-dev.yaml stop -t 1
	docker compose -f docker-compose-dev.yaml down
.PHONY: docker-compose-down