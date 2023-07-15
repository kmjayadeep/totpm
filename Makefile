all: build test

build:
	go build -o bin/server cmd/server/server.go
	go build -o bin/totp cmd/cli/cli.go

test:
	go test -v ./...

css:
	@tailwindcss -i ./assets/app.css -o ./assets/dist/app.css 

run: css
	@go run cmd/server/server.go


fmt:
	@go fmt ./...

watch:
	@ag -l | entr -r go run cmd/server/server.go

clean:
	go clean
	rm ./bin/server


# Cli commands
run-cli:
	@go run cmd/cli/cli.go
