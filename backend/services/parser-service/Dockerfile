FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
WORKDIR /app/services/parser-service
RUN go build -o parser-service ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/services/parser-service/parser-service .

EXPOSE 8081
CMD ["./parser-service"]
