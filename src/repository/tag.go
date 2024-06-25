package repository

import (
	"context"
	"gorm.io/gorm"
	"test-asset-fendr/src/models"
)

type TagRepositoryInterface interface {
	FindOne(ctx context.Context, query string, args ...interface{}) (models.Tag, error)
	Create(tx *gorm.DB, data *models.Tag) error
}

type TagRepository struct {
	Database *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{
		Database: db,
	}
}

func (r *TagRepository) Create(tx *gorm.DB, data *models.Tag) error {
	if err := tx.Model(models.Tag{}).Create(data).Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *TagRepository) FindOne(ctx context.Context, query string, args ...interface{}) (models.Tag, error) {
	var tag models.Tag

	if err := r.Database.WithContext(ctx).Model(models.Tag{}).Where(query, args...).First(&tag).Error; err != nil {
		return models.Tag{}, err
	}

	return tag, nil
}
