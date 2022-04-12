FROM golang:1.18
WORKDIR /go/src/github.com/exsmund/grpc-psql-example

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY ./proto/. ./proto
COPY ./packages/. ./packages
WORKDIR /go/src/github.com/exsmund/grpc-psql-example/services/logger
COPY ./services/logger/. .
RUN go build -o app .

ENTRYPOINT [ "./app" ]
