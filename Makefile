# Makefile
run:
	go run cmd/api/main.go

build:
	go build -o bin/shorty cmd/api/main.go

test:
	go test -v ./...

fmt:
	go fmt ./...

clean:
	rm -rf bin/