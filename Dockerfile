FROM golang:1.20.7 AS build

WORKDIR /app

COPY . .

# Build the Go application
RUN CGO_ENABLED=1 GOOS=linux go build -o main cmd/server/server.go

FROM ubuntu:23.10

WORKDIR /app

COPY --from=build /app/main .
COPY --from=build /app/views views/
COPY --from=build /app/assets assets/

EXPOSE 3000

ENTRYPOINT ["/app/main"]

