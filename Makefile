
viamflir: *.go cmd/module/*.go
	go build -tags netgo,osusergo -o viamflir cmd/module/cmd.go

test:
	go test

lint:
	gofmt -w -s .

module: viamflir
	tar czf module.tar.gz viamflir

all: module test
