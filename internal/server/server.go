package server

import (
	"context"

	"github.com/Koyo-os/vote-crud-service/internal/entity"
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
	if err := s.service.Update(req.ID, req.Key, req.Value); err != nil {
		s.logger.Error("failed update vote",
			zap.String("id", req.ID),
			zap.String("key", req.Key),
			zap.Error(err))

		return &protobuf.Response{
			Error: err.Error(),
			Ok:    false,
		}, err
	}

	return &protobuf.Response{
		Error: "",
		Ok:    true,
	}, nil
}

func (s *Server) Delete(ctx context.Context, req *protobuf.RequestDelete) (*protobuf.Response, error) {
	if err := s.service.Delete(req.ID); err != nil {
		s.logger.Error("failed delete vote",
			zap.String("id", req.ID),
			zap.Error(err))

		return &protobuf.Response{
			Error: err.Error(),
			Ok:    false,
		}, err
	}

	return &protobuf.Response{
		Error: "",
		Ok:    true,
	}, nil
}

func (s *Server) Create(ctx context.Context, req *protobuf.RequestCreate) (*protobuf.Response, error) {
	vote, err := entity.ToEntityVote(req.Vote)
	if err != nil{
		return nil, err
	}
	
	if err = s.service.Create(vote); err != nil {
		s.logger.Error("failed create vote",
			zap.String("id", req.Vote.ID),
			zap.Error(err))

		return &protobuf.Response{
			Error: err.Error(),
			Ok:    false,
		}, err
	}

	return &protobuf.Response{
		Error: "",
		Ok:    true,
	}, nil
}
