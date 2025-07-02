package repository

import (
	"errors"
	"fmt"

	"github.com/Koyo-os/vote-crud-service/internal/entity"
	"github.com/Koyo-os/vote-crud-service/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct{
	db *gorm.DB
	logger *logger.Logger
}

func NewRepository(db *gorm.DB, logger *logger.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func (r *Repository) Create(Vote *entity.Vote) error {
	if Vote == nil {
		return errors.New("Vote is nil")
	}
	if err := r.db.Create(Vote).Error; err != nil {
		r.logger.Error("failed to create Vote", zap.Error(err))
		return err
	}
	return nil
}

func (r *Repository) Get(id string) (*entity.Vote, error) {
	var Vote entity.Vote
	if err := r.db.Where("id = ?", id).First(&Vote).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Error("failed to get Vote by id",
			zap.String("id", id),
			zap.Error(err))

		return nil, err
	}
	
	return &Vote, nil
}

func (r *Repository) Update(id string, key string, value interface{}) error {
	updates := map[string]interface{}{key: value}
	if err := r.db.Model(&entity.Vote{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		r.logger.Error("failed to update Vote",
			zap.String("key", key),
			zap.String("id", id),
			zap.Error(err))

		return err
	}
	return nil
}

func (r *Repository) Delete(id string) error {
	if err := r.db.Where("id = ?", id).Delete(&entity.Vote{}).Error; err != nil {
		r.logger.Error("failed to delete Vote",
			zap.String("id", id),
			zap.Error(err))

		return err
	}
	return nil
}

func (r *Repository) GetMore(key string, value interface{}) ([]entity.Vote, error) {
	var Votes []entity.Vote
	query := fmt.Sprintf("%s = ?", key)
	if err := r.db.Where(query, value).Find(&Votes).Error; err != nil {
		r.logger.Error("failed to get Votes by",
			zap.String("key", key))
			
		return nil, err
	}
	return Votes, nil
}