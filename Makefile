all: gen build test
build:
	go build -o gones cmd/gones/main.go
gen:
	go generate ./...
test:
	go test ./...

.PHONY: all gen build test
