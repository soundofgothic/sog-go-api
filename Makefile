.PHONY: dev, docker

dev:
	cd development && docker-compose up -d --build

docker:
	docker build . -f docker/Dockerfile -t soundofgothic/backend:latest