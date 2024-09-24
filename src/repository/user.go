package repository

import (
	"context"
	"gorm.io/gorm"
	"test-edot/src/models"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *models.User) error
	FindOne(ctx context.Context, selectField, query string, args ...any) (models.User, error)
}

type UserRepository struct {
	Database *gorm.DB
	Tx       *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		Database: db,
	}
}

func (r UserRepository) Create(ctx context.Context, user *models.User) error {
	if err := r.Database.WithContext(ctx).Model(models.User{}).Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (r UserRepository) FindOne(ctx context.Context, selectField, query string, args ...any) (models.User, error) {
	var user models.User
	dbCon := r.Database.WithContext(ctx).Model(models.User{})

	if selectField != "*" {
		dbCon = dbCon.Select(selectField)
	}

	if err := dbCon.Where(query, args...).Take(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}
