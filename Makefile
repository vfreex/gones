all: gen build test
build: gen
	go build cmd/gones/gones.go
gen:
	go generate ./...
test: build
	go test ./...
deps:
	go get golang.org/x/tools/cmd/stringer
	go get fyne.io/fyne
	go get ./...

.PHONY: all gen build test deps
