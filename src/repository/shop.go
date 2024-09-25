package repository

import (
	"context"
	"gorm.io/gorm"
	"test-edot/src/models"
)

type ShopRepositoryInterface interface {
	Create(ctx context.Context, Shop *models.Shop) error
	FindOne(ctx context.Context, selectField, query string, args ...any) (models.Shop, error)
}

type ShopRepository struct {
	Database *gorm.DB
	Tx       *gorm.DB
}

func NewShopRepository(db *gorm.DB) *ShopRepository {
	return &ShopRepository{
		Database: db,
	}
}

func (r *ShopRepository) Create(ctx context.Context, Shop *models.Shop) error {
	if err := r.Database.WithContext(ctx).Model(models.Shop{}).Create(Shop).Error; err != nil {
		return err
	}

	return nil
}

func (r *ShopRepository) FindOne(ctx context.Context, selectField, query string, args ...any) (models.Shop, error) {
	var Shop models.Shop
	dbCon := r.Database.WithContext(ctx).Model(models.Shop{})

	if selectField != "*" {
		dbCon = dbCon.Select(selectField)
	}

	if err := dbCon.Where(query, args...).Take(&Shop).Error; err != nil {
		return models.Shop{}, err
	}

	return Shop, nil
}
