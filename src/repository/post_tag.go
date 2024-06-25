package repository

import (
	"gorm.io/gorm"
	"test-asset-fendr/src/models"
)

type PostTagRepositoryInterface interface {
	Create(tx *gorm.DB, data *models.PostTag) error
}

type PostTagRepository struct {
	Database *gorm.DB
}

func NewPostTagRepository(db *gorm.DB) *PostTagRepository {
	return &PostTagRepository{
		Database: db,
	}
}

func (r *PostTagRepository) Create(tx *gorm.DB, data *models.PostTag) error {
	if err := tx.Model(models.PostTag{}).Create(data).Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
