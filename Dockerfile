# Build stage
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache build-base

WORKDIR /app
COPY . .

# Static build with musl
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ravly ./main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/ravly .

ENTRYPOINT ["./ravly"]
