all: sausage

sausage: *.go
	go build -x -o sausage

sausage-test:
	go test -v -cover

sausage-clean:
	go clean -x

get-deps:
	go get -v github.com/pivotal-golang/bytefmt
	go get -v gopkg.in/alecthomas/kingpin.v2

test: sausage-test

clean: sausage-clean

.PHONY: all sausage-test sausage-clean get-deps test clean
