module github.com/exsmund/grpc-psql-example/services/server

go 1.18

require (
	github.com/exsmund/grpc-psql-example/proto v0.0.0-unpublished
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/lib/pq v1.10.4
	google.golang.org/grpc v1.45.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.18.1 // indirect
	golang.org/x/net v0.0.0-20210428140749-89ef3d95e781 // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)

replace github.com/exsmund/grpc-psql-example/proto/user v0.0.0-unpublished => ../../proto/user

replace github.com/exsmund/grpc-psql-example/proto v0.0.0-unpublished => ../../proto
