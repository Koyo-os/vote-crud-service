package server

import (
	"context"

	"github.com/Koyo-os/vote-crud-service/internal/service"
	"github.com/Koyo-os/vote-crud-service/pkg/api/protobuf"
	"github.com/Koyo-os/vote-crud-service/pkg/logger"
	"github.com/Koyo-os/vote-crud-service/pkg/retrier"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	protobuf.UnimplementedVoteServiceServer
	logger  *logger.Logger
	service *service.Service
}

func NewServer(service *service.Service, logger *logger.Logger) *Server {
	return &Server{
		service: service,
		logger:  logger,
	}
}

func (s *Server) Get(req *protobuf.RequestGet, resp grpc.ServerStreamingServer[protobuf.Vote]) error {
	respChan, err := s.service.GetByPollID(req.ID)
	if err != nil {
		return err
	}

	for vote := range respChan {
		if err = retrier.Do(3, 1, func() error {
			return resp.Send(vote.ToProtobuf())
		}); err != nil {
			s.logger.Error("failed send vote",
				zap.String("vote_id", vote.ID.String()),
				zap.Error(err))
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Update(ctx context.Context, req *protobuf.RequestUpdate) (*protobuf.Response, error) {
}

func (s *Server) Delete(ctx context.Context, req *protobuf.RequestDelete) (*protobuf.Response, error) {
}

func (s *Server) Create(ctx context.Context, req *protobuf.RequestCreate) (*protobuf.Response, error) {
}
