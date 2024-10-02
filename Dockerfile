FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /nimbus ./cmd/nimbus

FROM alpine:latest

COPY --from=builder /nimbus /nimbus 

EXPOSE 8000

ENTRYPOINT ["/nimbus"]