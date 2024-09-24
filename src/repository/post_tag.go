package repository

import (
	"gorm.io/gorm"
	"test-edot/src/models"
)

type PostTagRepositoryInterface interface {
	Create(tx *gorm.DB, data *models.PostTag) error
	DeleteOne(tx *gorm.DB, query string, args ...any) error
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

func (r *PostTagRepository) DeleteOne(tx *gorm.DB, query string, args ...any) error {
	if err := tx.Model(models.PostTag{}).Where(query, args...).Delete(models.PostTag{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
