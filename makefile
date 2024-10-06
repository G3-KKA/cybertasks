# Setting up workspace
WORKSPACE ?= $(shell pwd)
include ${WORKSPACE}/.env
export $(shell sed 's/=.*//' .env)
export WORKSPACE

lint:
	golangci-lint run ./...
test:
	go test -v -race ./...
swag:
	swag init -g ./internal/controller/handler/routes.go 