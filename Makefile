all: gen build test
build:
	go build cmd/gones/main.go
gen:
	go generate ./...
test:
	go test ./...

.PHONY: all gen build test
