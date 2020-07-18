GO=go

.PHONY: build setup run fmt

build:
	cd ./src && ${GO} build -o ./main main.go

run:
	cd ./src && ./main

fmt:
	go fmt ./src && goimports -w ./src
