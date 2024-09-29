package repository

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"test-edot/src/models"
)

type OrderDetailRepositoryInterface interface {
	Create(tx *gorm.DB, OrderDetail *models.OrderDetail) error
	FindOne(ctx context.Context, selectField, query string, args ...any) (models.OrderDetail, error)
	Begin() *gorm.DB
	UpdateOneTx(tx *gorm.DB, updateOrderDetail *models.OrderDetail, selectFields, query string, args ...interface{}) error
	FindTx(tx *gorm.DB, selectField, query string, args ...any) ([]models.OrderDetail, error)
}

type OrderDetailRepository struct {
	Database *gorm.DB
	Tx       *gorm.DB
}

func NewOrderDetailRepository(db *gorm.DB) *OrderDetailRepository {
	return &OrderDetailRepository{
		Database: db,
	}
}

func (r *OrderDetailRepository) Begin() *gorm.DB {
	return r.Database.Begin()
}

func (r *OrderDetailRepository) Create(tx *gorm.DB, OrderDetail *models.OrderDetail) error {
	if err := tx.Model(models.OrderDetail{}).Create(OrderDetail).Error; err != nil {
		return err
	}

	return nil
}

func (r *OrderDetailRepository) UpdateOneTx(tx *gorm.DB, updateOrderDetail *models.OrderDetail, selectFields, query string, args ...interface{}) error {
	dbConn := tx.Model(models.OrderDetail{})

	if selectFields != "*" {
		dbConn = dbConn.Select(strings.Split(selectFields, ","))
	}

	if err := dbConn.Where(query, args...).Debug().Updates(updateOrderDetail).Error; err != nil {
		return err
	}

	return nil
}

func (r *OrderDetailRepository) FindOne(ctx context.Context, selectField, query string, args ...any) (models.OrderDetail, error) {
	var OrderDetail models.OrderDetail
	dbCon := r.Database.WithContext(ctx).Model(models.OrderDetail{})

	if selectField != "*" {
		dbCon = dbCon.Select(selectField)
	}

	if err := dbCon.Where(query, args...).Take(&OrderDetail).Error; err != nil {
		return models.OrderDetail{}, err
	}

	return OrderDetail, nil
}

func (r *OrderDetailRepository) FindTx(tx *gorm.DB, selectField, query string, args ...any) ([]models.OrderDetail, error) {
	var orderDetails []models.OrderDetail
	db := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(models.OrderDetail{})

	if selectField != "*" {
		db = db.Select(selectField)
	}

	if err := db.Where(query, args...).Take(&orderDetails).Error; err != nil {
		return []models.OrderDetail{}, err
	}

	return orderDetails, nil
}
