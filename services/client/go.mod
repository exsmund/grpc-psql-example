module github.com/exsmund/grpc-psql-example/services/client

go 1.18

require (
	github.com/exsmund/grpc-psql-example/proto v0.0.0-unpublished
	google.golang.org/grpc v1.45.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/sys v0.0.0-20200323222414-85ca7c5b95cd // indirect
	golang.org/x/text v0.3.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)

replace github.com/exsmund/grpc-psql-example/proto/user v0.0.0-unpublished => ../../proto/user

replace github.com/exsmund/grpc-psql-example/proto v0.0.0-unpublished => ../../proto
