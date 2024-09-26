package repository

import (
	"context"
	"gorm.io/gorm"
	"test-edot/src/models"
)

type ProductRepositoryInterface interface {
	Create(tx *gorm.DB, Product *models.Product) error
	FindOne(ctx context.Context, selectField, query string, args ...any) (models.Product, error)
	Begin() *gorm.DB
	GetProductDetail(ctx context.Context, offset, limit int, selectField, query string, args ...any) ([]models.ProductDetail, error)
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

func (r *ProductRepository) Begin() *gorm.DB {
	return r.Database.Begin()
}

func (r *ProductRepository) Create(tx *gorm.DB, Product *models.Product) error {
	if err := tx.Model(models.Product{}).Create(Product).Error; err != nil {
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

func (r *ProductRepository) GetProductDetail(ctx context.Context, offset, limit int, selectField, query string, args ...any) ([]models.ProductDetail, error) {
	var products []models.ProductDetail
	err := r.Database.WithContext(ctx).Model(models.ProductDetail{}).
		Preload("Shop", func(db *gorm.DB) *gorm.DB {
			return db.Order("id,name")
		}).
		Preload("Stock", func(db *gorm.DB) *gorm.DB {
			return db.Order("id,stock")
		}).
		Offset(offset).Limit(limit).
		Select(selectField).Where(query, args...).Debug().Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}
