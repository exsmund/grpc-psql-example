package main

import (
	lpb "github.com/exsmund/grpc-psql-example/proto/logger"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type loggerServer struct {
	lpb.UnimplementedLoggerRepoServer
	ch *CH
}

func (s *loggerServer) GetTail(req *lpb.GetTailRequest, stream lpb.LoggerRepo_GetTailServer) error {
	rows, err := s.ch.Read()
	if err != nil {
		return err
	}
	for _, r := range *rows {
		stream.Send(&lpb.LogRow{Ts: timestamppb.New(r.Dt), Msg: r.Msg})
	}

	return nil
}

func newServer(ch *CH) *loggerServer {
	s := &loggerServer{
		ch: ch,
	}
	return s
}
