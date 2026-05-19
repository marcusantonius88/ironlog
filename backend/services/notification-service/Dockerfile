FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
WORKDIR /app/services/notification-service
RUN go build -o notification-service ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/services/notification-service/notification-service .

EXPOSE 8085
CMD ["./notification-service"]
