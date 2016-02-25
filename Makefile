all: sausage

sausage: *.go
	go build -x -o sausage
