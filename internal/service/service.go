package service

import (
	"context"

	"github.com/Koyo-os/vote-crud-service/internal/entity"
)

type (
	Repository interface {
		Create(vote *entity.Vote) error
		Update(string, string, interface{}) error
		Delete(string) error
		Get(string) (*entity.Vote, error)
		GetBy(string, interface{}) ([]entity.Vote, error)
		GetByPollID(context.Context, string) (chan entity.Vote, error)
	}

	Service struct {
		repo Repository
	}
)

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(vote *entity.Vote) error {
	return s.repo.Create(vote)
}

func (s *Service) Update(id, key string, value interface{}) error {
	return s.repo.Update(id, key, value)
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *Service) Get(id string) (*entity.Vote, error) {
	return s.repo.Get(id)
}

func (s *Service) GetBy(id string, value interface{}) ([]entity.Vote, error) {
	return s.repo.GetBy(id, value)
}

func (s *Service) GetByPollID(ctx context.Context, id string) (chan entity.Vote, error) {
	return s.repo.GetByPollID(ctx, id)
}
