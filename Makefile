VERSION := $(shell git rev-parse --short HEAD)
CHANGES := $(shell git status --porcelain)

.PHONY: dev docker docker-push

dev:
	cd development && docker-compose up -d --build

changes:
	@echo $(CHANGES)

docker:
	docker build . -f docker/Dockerfile -t soundofgothic/api:$(VERSION)

docker-push: docker
	if [ -n "$(CHANGES)" ]; then echo "You have uncommited changes, aborting"; exit 1; fi
	docker push soundofgothic/api:$(VERSION)

version:
	@echo soundofgothic/api:$(VERSION)

list:
	@echo "dev - run development environment"
	@echo "docker - build docker image"
	@echo "docker-push - push docker image to registry"
	@echo "version - print current version"
	@echo "list - print this list"