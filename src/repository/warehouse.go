package repository

import (
	"context"
	"gorm.io/gorm"
	"strings"
	"test-edot/src/models"
)

type WarehouseRepositoryInterface interface {
	Create(ctx context.Context, Warehouse *models.Warehouse) error
	FindOne(ctx context.Context, selectField, query string, args ...any) (models.Warehouse, error)
	Find(ctx context.Context, selectField, query string, args ...any) ([]models.Warehouse, error)
	Update(ctx context.Context, updatedField models.Warehouse, selectFields, query string, args ...any) error
}

type WarehouseRepository struct {
	Database *gorm.DB
	Tx       *gorm.DB
}

func NewWarehouseRepository(db *gorm.DB) *WarehouseRepository {
	return &WarehouseRepository{
		Database: db,
	}
}

func (r *WarehouseRepository) Create(ctx context.Context, Warehouse *models.Warehouse) error {
	if err := r.Database.WithContext(ctx).Model(models.Warehouse{}).Create(Warehouse).Error; err != nil {
		return err
	}

	return nil
}

func (r *WarehouseRepository) FindOne(ctx context.Context, selectField, query string, args ...any) (models.Warehouse, error) {
	var Warehouse models.Warehouse
	dbCon := r.Database.WithContext(ctx).Model(models.Warehouse{})

	if selectField != "*" {
		dbCon = dbCon.Select(selectField)
	}

	if err := dbCon.Where(query, args...).Take(&Warehouse).Error; err != nil {
		return models.Warehouse{}, err
	}

	return Warehouse, nil
}

func (r *WarehouseRepository) Find(ctx context.Context, selectField, query string, args ...any) ([]models.Warehouse, error) {
	var warehouses []models.Warehouse
	dbCon := r.Database.WithContext(ctx).Model(models.Warehouse{})

	if selectField != "*" {
		dbCon = dbCon.Select(selectField)
	}

	if err := dbCon.Where(query, args...).Find(&warehouses).Error; err != nil {
		return []models.Warehouse{}, err
	}

	return warehouses, nil
}

func (r *WarehouseRepository) Update(ctx context.Context, updatedField models.Warehouse, selectFields, query string, args ...any) error {
	dbConn := r.Database.WithContext(ctx).Model(models.Warehouse{})

	if selectFields != "*" {
		dbConn = dbConn.Select(strings.Split(selectFields, ","))
	}

	if err := dbConn.Where(query, args...).Updates(&updatedField).Error; err != nil {
		return err
	}

	return nil
}
