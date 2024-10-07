# Setting up workspace
WORKSPACE ?= $(shell pwd)
include ${WORKSPACE}/.env
export $(shell sed 's/=.*//' .env)
export WORKSPACE

compose-run:
	docker compose up --build
run-local:
	go run ${WORKSPACE}/bin/cybertask
lint:
	golangci-lint run ./...
test:
	go test -v -race ./...
swag:
	swag init -g ./internal/controller/handler/routes.go 