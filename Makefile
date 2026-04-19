include .env

.PHONY: help up down build test migrate-force

up:
	docker compose up -d

down:
	docker compose down

build:
	docker compose up -d --build

test:
	go test -v ./...

logs:
	docker compose logs -f

setup:
	cp .env.example .env