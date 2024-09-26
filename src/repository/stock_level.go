package repository

import (
	"context"
	"gorm.io/gorm"
	"test-edot/src/models"
)

type StockLevelRepositoryInterface interface {
	Create(tx *gorm.DB, stockLevel *models.StockLevel) error
	FindOne(ctx context.Context, selectField, query string, args ...any) (models.StockLevel, error)
}

type StockLevelRepository struct {
	Database *gorm.DB
	Tx       *gorm.DB
}

func NewStockLevelRepository(db *gorm.DB) *StockLevelRepository {
	return &StockLevelRepository{
		Database: db,
	}
}

func (r *StockLevelRepository) Create(tx *gorm.DB, stockLevel *models.StockLevel) error {
	if err := tx.Model(models.StockLevel{}).Create(stockLevel).Error; err != nil {
		return err
	}

	return nil
}

func (r *StockLevelRepository) FindOne(ctx context.Context, selectField, query string, args ...any) (models.StockLevel, error) {
	var stockLevel models.StockLevel
	dbCon := r.Database.WithContext(ctx).Model(models.StockLevel{})

	if selectField != "*" {
		dbCon = dbCon.Select(selectField)
	}

	if err := dbCon.Where(query, args...).Take(&stockLevel).Error; err != nil {
		return models.StockLevel{}, err
	}

	return stockLevel, nil
}
