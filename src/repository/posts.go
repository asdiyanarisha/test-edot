package repository

import (
	"context"
	"gorm.io/gorm"
	"test-asset-fendr/src/models"
)

type PostRepositoryInterface interface {
	FindOne(ctx context.Context, query string, args ...interface{}) (models.Post, error)
	Create(data *models.Post) error
	GetPosts(ctx context.Context, offset, limit int) ([]models.PostWithTag, error)
	DbTransaction
}

type DbTransaction interface {
	Begin(ctx context.Context) *gorm.DB
	Rollback()
	Error() error
	Commit() error
}

type PostRepository struct {
	Database *gorm.DB
	Tx       *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{
		Database: db,
	}
}

func (r *PostRepository) Begin(ctx context.Context) *gorm.DB {
	r.Tx = r.Database.WithContext(ctx).Begin()
	return r.Tx
}

func (r *PostRepository) Rollback() {
	r.Tx.Rollback()
}

func (r *PostRepository) Commit() error {
	return r.Tx.Commit().Error
}

func (r *PostRepository) Error() error {
	return r.Tx.Error
}

func (r *PostRepository) Create(data *models.Post) error {
	if err := r.Tx.Model(models.Post{}).Create(data).Error; err != nil {
		r.Rollback()
		return err
	}

	return nil
}

func (r *PostRepository) FindOne(ctx context.Context, query string, args ...interface{}) (models.Post, error) {
	var post models.Post

	if err := r.Database.WithContext(ctx).Model(models.Post{}).Where(query, args...).First(&post).Error; err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func (r *PostRepository) GetPosts(ctx context.Context, offset, limit int) ([]models.PostWithTag, error) {
	var posts []models.PostWithTag

	if err := r.Database.WithContext(ctx).Model(models.PostWithTag{}).Preload("Tags").
		Offset(offset).
		Limit(limit).
		Find(&posts).Error; err != nil {
		return []models.PostWithTag{}, err
	}

	return posts, nil
}
