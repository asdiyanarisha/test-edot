package repository

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"test-edot/src/models"
)

type OrderRepositoryInterface interface {
	Create(tx *gorm.DB, Order *models.Order) error
	FindOne(ctx context.Context, selectField, query string, args ...any) (models.Order, error)
	Begin() *gorm.DB
	FindOneTx(tx *gorm.DB, fields, query string, args ...interface{}) (models.Order, error)
	UpdateOneTx(tx *gorm.DB, updateOrder *models.Order, selectFields, query string, args ...interface{}) error
	FindAll(ctx context.Context, selectField, query string, args ...any) ([]models.Order, error)
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

func (r *OrderRepository) FindAll(ctx context.Context, selectField, query string, args ...any) ([]models.Order, error) {
	var orders []models.Order
	dbCon := r.Database.WithContext(ctx).Model(models.Order{})

	if selectField != "*" {
		dbCon = dbCon.Select(selectField)
	}

	if err := dbCon.Where(query, args...).Find(&orders).Error; err != nil {
		return []models.Order{}, err
	}

	return orders, nil
}

func (r *OrderRepository) FindOneTx(tx *gorm.DB, fields, query string, args ...interface{}) (models.Order, error) {
	var order models.Order
	db := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(models.Order{})

	if fields != "*" {
		db = db.Select(fields)
	}

	if err := db.Where(query, args...).Take(&order).Error; err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (r *OrderRepository) UpdateOneTx(tx *gorm.DB, updateOrder *models.Order, selectFields, query string, args ...interface{}) error {
	dbConn := tx.Model(models.Order{})

	if selectFields != "*" {
		dbConn = dbConn.Select(strings.Split(selectFields, ","))
	}

	if err := dbConn.Where(query, args...).Debug().Updates(updateOrder).Error; err != nil {
		return err
	}

	return nil
}
