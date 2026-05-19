FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
WORKDIR /app/services/analytics-service
RUN go build -o analytics-service ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/services/analytics-service/analytics-service .

EXPOSE 8082
CMD ["./analytics-service"]
