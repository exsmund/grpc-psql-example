package main

import (
	"context"

	lpb "github.com/exsmund/grpc-psql-example/proto/logger"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type loggerServer struct {
	lpb.UnimplementedLoggerRepoServer
	ch *CH
}

func (s *loggerServer) GetTail(ctx context.Context, req *lpb.GetTailRequest) (*lpb.Log, error) {
	rows, err := s.ch.Read()
	if err != nil {
		return nil, err
	}
	result := make([]*lpb.LogRow, len(*rows))
	for i, r := range *rows {
		result[i] = &lpb.LogRow{Ts: timestamppb.New(r.Dt), Msg: r.Msg}
	}

	return &lpb.Log{Logs: result}, nil
}

func newServer(ch *CH) *loggerServer {
	s := &loggerServer{
		ch: ch,
	}
	return s
}
