package repository

import (
	"context"
	"gorm.io/gorm"
	"test-edot/src/models"
)

type OrderRepositoryInterface interface {
	Create(tx *gorm.DB, Order *models.Order) error
	FindOne(ctx context.Context, selectField, query string, args ...any) (models.Order, error)
	Begin() *gorm.DB
}

type OrderRepository struct {
	Database *gorm.DB
	Tx       *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		Database: db,
	}
}

func (r *OrderRepository) Begin() *gorm.DB {
	return r.Database.Begin()
}

func (r *OrderRepository) Create(tx *gorm.DB, orderDetail *models.Order) error {
	if err := tx.Model(models.Order{}).Create(orderDetail).Error; err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) FindOne(ctx context.Context, selectField, query string, args ...any) (models.Order, error) {
	var Order models.Order
	dbCon := r.Database.WithContext(ctx).Model(models.Order{})

	if selectField != "*" {
		dbCon = dbCon.Select(selectField)
	}

	if err := dbCon.Where(query, args...).Take(&Order).Error; err != nil {
		return models.Order{}, err
	}

	return Order, nil
}
