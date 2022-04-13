.PHONY : all

all : run tests lint build

include .env
export

run:
	go run main.go

tests:
	go test -cover  ./...

lint:
	gofmt -w .
	
build: 
	go build main.go