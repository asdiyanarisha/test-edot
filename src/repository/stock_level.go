package repository

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"test-edot/src/models"
)

type StockLevelRepositoryInterface interface {
	Create(tx *gorm.DB, stockLevel *models.StockLevel) error
	FindOne(ctx context.Context, selectField, query string, args ...any) (models.StockLevel, error)
	FindOneTx(tx *gorm.DB, order, query string, args ...interface{}) (models.StockLevelProduct, error)
	FindTx(tx *gorm.DB, order, query string, args ...interface{}) ([]models.StockLevelProduct, error)
	UpdateOneTx(tx *gorm.DB, updateStockLevel *models.StockLevel, selectFields, query string, args ...interface{}) error
	SumStockWarehouse(ctx context.Context, query string, args ...any) (models.StockWarehouse, error)
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

func (r *StockLevelRepository) FindOneTx(tx *gorm.DB, order, query string, args ...interface{}) (models.StockLevelProduct, error) {
	var transaction models.StockLevelProduct
	db := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(models.StockLevelProduct{})
	if order != "" {
		db = db.Order(order)
	}

	if err := db.Preload("Product", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,name,price")
	}).Where(query, args...).Debug().Take(&transaction).Error; err != nil {
		return models.StockLevelProduct{}, err
	}

	return transaction, nil
}

func (r *StockLevelRepository) FindTx(tx *gorm.DB, order, query string, args ...interface{}) ([]models.StockLevelProduct, error) {
	var stocks []models.StockLevelProduct
	db := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(models.StockLevelProduct{})
	if order != "" {
		db = db.Order(order)
	}

	if err := db.Preload("Product", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,name,price")
	}).Where(query, args...).Find(&stocks).Error; err != nil {
		return []models.StockLevelProduct{}, err
	}

	return stocks, nil
}

func (r *StockLevelRepository) UpdateOneTx(tx *gorm.DB, updateStockLevel *models.StockLevel, selectFields, query string, args ...interface{}) error {
	dbConn := tx.Model(models.StockLevel{})

	if selectFields != "*" {
		dbConn = dbConn.Select(strings.Split(selectFields, ","))
	}

	if err := dbConn.Where(query, args...).Debug().Updates(updateStockLevel).Error; err != nil {
		return err
	}

	return nil
}

func (r *StockLevelRepository) SumStockWarehouse(ctx context.Context, query string, args ...any) (models.StockWarehouse, error) {
	var res models.StockWarehouse

	if err := r.Database.WithContext(ctx).Model(models.StockLevel{}).
		Select("sum(stock) as stock_count, sum(reserved_stock) as reserved_stock_count").Where(query, args...).
		Take(&res).Error; err != nil {
		return models.StockWarehouse{}, err
	}

	return res, nil
}
