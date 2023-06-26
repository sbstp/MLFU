.PHONY: build
build:
	go build -o dist/mlfu -ldflags "-s -w" main.go