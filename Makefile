all: build test

build:
	go build -o bin/server cmd/server/server.go

test:
	go test -v ./...

run:
	@go run cmd/server/server.go

fmt:
	@go fmt ./...

clean:
	go clean
	rm ./bin/server
