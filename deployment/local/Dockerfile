FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o denet ./cmd/denet/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/denet .

EXPOSE 8081
CMD ["./denet"]