FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
WORKDIR /app/services/workout-service
RUN go build -o workout-service ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/services/workout-service/workout-service .

EXPOSE 8080
CMD ["./workout-service"]
