FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener ./cmd/main.go
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/url-shortener .
COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["./url-shortener"]
