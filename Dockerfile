# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /nimbus ./cmd/nimbus

# Production stage
FROM alpine:latest

COPY --from=builder /nimbus /nimbus
EXPOSE 8080
ENTRYPOINT ["/nimbus"]
