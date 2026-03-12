# Simple Makefile for the leaderboard system

build:
	go build -o bin/server cmd/server/main.go

run: build
	./bin/server

test:
	go test ./...

clean:
	rm -rf bin/

install:
	go mod tidy

docker-build:
	docker build -t leaderboard-system .

docker-run:
	docker run -p 8080:8080 leaderboard-system

.PHONY: build run test clean install docker-build docker-run