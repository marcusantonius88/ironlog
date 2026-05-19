FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
WORKDIR /app/services/recommendation-service
RUN go build -o recommendation-service ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/services/recommendation-service/recommendation-service .

EXPOSE 8084
CMD ["./recommendation-service"]
