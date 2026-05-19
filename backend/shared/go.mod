module github.com/marcus/ironlog/backend/shared

go 1.21

require (
	github.com/google/uuid v1.5.0
	github.com/lib/pq v1.10.9
	github.com/segmentio/kafka-go v0.4.46
	github.com/go-redis/redis/v8 v8.11.5
	go.opentelemetry.io/otel v1.21.0
	go.opentelemetry.io/otel/exporters/prometheus v0.43.0
	go.opentelemetry.io/otel/metric v1.21.0
	go.opentelemetry.io/otel/sdk/metric v1.21.0
	go.opentelemetry.io/otel/semconv v1.21.0
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/klauspost/compress v1.17.4 // indirect
	github.com/pierrec/lz4/v4 v4.1.18 // indirect
	github.com/prometheus/client_golang v1.18.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.45.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	go.opentelemetry.io/otel/sdk v1.21.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
)
