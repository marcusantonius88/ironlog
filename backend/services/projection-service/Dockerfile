FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
WORKDIR /app/services/projection-service
RUN go build -o projection-service ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/services/projection-service/projection-service .

EXPOSE 8083
CMD ["./projection-service"]
