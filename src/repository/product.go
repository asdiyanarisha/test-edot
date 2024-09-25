package repository

import (
	"context"
	"gorm.io/gorm"
	"test-edot/src/models"
)

type ProductRepositoryInterface interface {
	Create(ctx context.Context, Product *models.Product) error
	FindOne(ctx context.Context, selectField, query string, args ...any) (models.Product, error)
}

type ProductRepository struct {
	Database *gorm.DB
	Tx       *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		Database: db,
	}
}

func (r *ProductRepository) Create(ctx context.Context, Product *models.Product) error {
	if err := r.Database.WithContext(ctx).Model(models.Product{}).Create(Product).Error; err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) FindOne(ctx context.Context, selectField, query string, args ...any) (models.Product, error) {
	var Product models.Product
	dbCon := r.Database.WithContext(ctx).Model(models.Product{})

	if selectField != "*" {
		dbCon = dbCon.Select(selectField)
	}

	if err := dbCon.Where(query, args...).Take(&Product).Error; err != nil {
		return models.Product{}, err
	}

	return Product, nil
}
