FROM golang:1.24.1-alpine AS builder

# Install build tools
RUN apk add --no-cache build-base

# Install golang-migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -tags musl -o fin-ai ./cmd

# Final stage
FROM alpine:3.14

# Install netcat-openbsd for nc command and ca-certificates
RUN apk update && apk add --no-cache netcat-openbsd ca-certificates

WORKDIR /app

# Copy the pre-built binary and the migrate tool from the builder stage
COPY --from=builder /app/fin-ai .
COPY --from=builder /go/bin/migrate /app/migrate

# Copy the entrypoint script
COPY docker-entrypoint.sh .
RUN chmod +x docker-entrypoint.sh

# Copy migrations
COPY migrations ./migrations

ENTRYPOINT ["/app/docker-entrypoint.sh"]
CMD ["./fin-ai"]