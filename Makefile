COMPOSE := docker-compose -f docker-compose.yml

.PHONY: down up build clean test

down:
	${COMPOSE} down

up:
	${COMPOSE} up -d

build:
	${COMPOSE} up -d --build

clean:
	${COMPOSE} down -v --rmi all
	
test:
	go test -v ./...