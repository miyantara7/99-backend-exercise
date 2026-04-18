SHELL := /bin/zsh

.PHONY: run run-alt build

run:
	go run ./cmd/dev

run-alt:
	go run ./cmd/dev --listing-port=6400 --user-port=6401 --public-port=7400

build:
	go build -o .run/dev-runner ./cmd/dev
