# Project variables
APP_NAME=sl-edu-service
DOCKER_USER=lakmaltharindu
DOCKER_IMAGE=$(DOCKER_USER)/$(APP_NAME)

# Default target (runs if you just type `make`)
.DEFAULT_GOAL := help

## Run the application locally
run:
	go run cmd/server/main.go

## Run tests
test:
	go test ./... -v

## Build Docker image
docker-build:
	docker build -t $(DOCKER_IMAGE):latest .

## Push Docker image to Docker Hub
docker-push: docker-build
	docker push $(DOCKER_IMAGE):latest

## Remove temporary files, binaries, Docker images
clean:
	go clean
	docker rmi $(DOCKER_IMAGE):latest || true

## Show available commands
help:
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?##' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  make %-15s %s\n", $$1, $$2}'
