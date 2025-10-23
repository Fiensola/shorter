# Makefile
.PHONY: run build test fmt vet clean migrate-up migrate-down

run:
	go run cmd/api/main.go

build:
	go build -o bin/shorty cmd/api/main.go

test:
	go test -v -cover ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf bin/

migrate-up:
	migrate -path migrations -database "postgres://user:user@localhost:5432/shorter?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://user:user@localhost:5432/shorter?sslmode=disable" down