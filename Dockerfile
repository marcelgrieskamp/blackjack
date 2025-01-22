FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o blackjack-server cmd/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/blackjack-server .
COPY --from=builder /app/static ./static

EXPOSE 8080
CMD ["./blackjack-server"] 