module github.com/exsmund/grpc-psql-example/services/logger

go 1.18

require (
	github.com/ClickHouse/clickhouse-go/v2 v2.0.12
	github.com/confluentinc/confluent-kafka-go v1.8.2
	github.com/exsmund/grpc-psql-example/proto v0.0.0-unpublished
	google.golang.org/grpc v1.45.0
	google.golang.org/protobuf v1.28.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/paulmach/orb v0.5.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.14 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	go.opentelemetry.io/otel v1.6.3 // indirect
	go.opentelemetry.io/otel/trace v1.6.3 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sys v0.0.0-20211110154304-99a53858aa08 // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
)

replace github.com/exsmund/grpc-psql-example/proto/logger v0.0.0-unpublished => ../../proto/logger

replace github.com/exsmund/grpc-psql-example/proto v0.0.0-unpublished => ../../proto
