FROM golang:1.23-alpine AS builder

RUN apk add --no-cache alpine-sdk git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o fin-ai ./cmd

# Final stage
FROM alpine:3.18

RUN apk update && apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/fin-ai .

EXPOSE 8080

CMD ["./fin-ai"]