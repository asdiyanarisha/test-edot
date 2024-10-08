package repository

import (
	"context"
	"gorm.io/gorm"
	"strings"
	"test-edot/src/models"
)

type PostRepositoryInterface interface {
	FindOne(ctx context.Context, query string, args ...interface{}) (models.Post, error)
	FindOneWithTag(ctx context.Context, query string, args ...interface{}) (models.PostWithTag, error)
	Create(data *models.Post) error
	GetPosts(ctx context.Context, offset, limit int) ([]models.PostWithTag, error)
	UpdateOne(tx *gorm.DB, data models.Post, updatedField, query string, args ...interface{}) error
	DeleteOne(tx *gorm.DB, query string, args ...any) error
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

func (r *PostRepository) DeleteOne(tx *gorm.DB, query string, args ...any) error {
	if err := tx.Model(models.Post{}).Where(query, args...).Delete(models.Post{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *PostRepository) UpdateOne(tx *gorm.DB, data models.Post, updatedField, query string, args ...interface{}) error {
	if err := tx.Model(models.Post{}).
		Select(strings.Split(updatedField, ",")).
		Where(query, args...).Debug().Updates(data).Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *PostRepository) FindOneWithTag(ctx context.Context, query string, args ...interface{}) (models.PostWithTag, error) {
	var post models.PostWithTag

	if err := r.Database.WithContext(ctx).Model(models.PostWithTag{}).Preload("Tags").
		Where(query, args...).
		Take(&post).Error; err != nil {
		return models.PostWithTag{}, err
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
