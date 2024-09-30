package warehouse

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
	"test-edot/constants"
	"test-edot/src/dto"
	"test-edot/src/factory"
	"test-edot/src/models"
	"test-edot/src/repository"
	"test-edot/util"
	"time"
)

type Service interface {
	AddWarehouse(ctx context.Context, payload dto.PayloadAddWarehouse, userClaim dto.UserClaimJwt) (dto.ResponseWarehouse, error)
	GetWarehouses(ctx context.Context, userClaim dto.UserClaimJwt, payload dto.ParameterQueryWarehouse) (any, error)
	ChangeStatusWarehouse(ctx context.Context, userClaim dto.UserClaimJwt, payload dto.ParameterChangeStatusWarehouse) error
	TransferProductWarehouse(ctx context.Context, userClaim dto.UserClaimJwt, fromId, toId int) error
}

type service struct {
	Log                   *zap.Logger
	UserRepository        repository.UserRepositoryInterface
	ShopRepository        repository.ShopRepositoryInterface
	ProductRepository     repository.ProductRepositoryInterface
	WarehouseRepository   repository.WarehouseRepositoryInterface
	StockLevelRepository  repository.StockLevelRepositoryInterface
	OrderDetailRepository repository.OrderDetailRepositoryInterface
}

func NewService(f *factory.Factory) Service {
	return &service{
		Log:                   f.Log,
		UserRepository:        f.UserRepository,
		ShopRepository:        f.ShopRepository,
		ProductRepository:     f.ProductRepository,
		WarehouseRepository:   f.WarehouseRepository,
		StockLevelRepository:  f.StockLevelRepository,
		OrderDetailRepository: f.OrderDetailRepository,
	}
}

func (s *service) ChangeStatusWarehouse(ctx context.Context, userClaim dto.UserClaimJwt, payload dto.ParameterChangeStatusWarehouse) error {
	warehouse, stock, err := s.InitiateDataWarehouse(ctx, payload, userClaim)
	if err != nil {
		return err
	}

	if warehouse.IsActive == payload.IsActive {
		return constants.StatusNotSamePrevious
	}

	if stock.StockCount != 0 && stock.ReservedStockCount != 0 {
		return constants.StockMustEmpty
	}

	updatedField := models.Warehouse{IsActive: payload.IsActive, UpdatedAt: time.Now().In(util.LocationTime)}
	if err := s.WarehouseRepository.Update(ctx, updatedField, "is_active,updated_at", "id = ?", payload.WarehouseId); err != nil {
		return err
	}

	return nil
}

func (s *service) InitiateDataWarehouse(ctx context.Context, payload dto.ParameterChangeStatusWarehouse, userClaim dto.UserClaimJwt) (models.Warehouse, models.StockWarehouse, error) {
	var (
		wg             sync.WaitGroup
		warehouse      models.Warehouse
		stockWarehouse models.StockWarehouse
	)

	errChan := make(chan error)
	wgDone := make(chan struct{})

	c, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	wg.Add(2)

	go func() {
		defer wg.Done()
		s.Log.Info("get warehouse")

		select {
		case <-c.Done():
			return
		default:
			res, err := s.WarehouseRepository.FindOne(c, "id,is_active", "id = ? and user_id = ?", payload.WarehouseId, userClaim.UserId)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					errChan <- constants.WarehouseNotFound
					return
				}

				s.Log.Error("error fetch warehouses", zap.Error(err), zap.Int("user_id", userClaim.UserId))
				errChan <- err
				return
			}

			warehouse = res
		}
	}()

	go func() {
		defer wg.Done()
		s.Log.Info("get stock warehouse")

		select {
		case <-c.Done():
			return
		default:
			res, err := s.StockLevelRepository.SumStockWarehouse(c, "warehouse_id = ?", payload.WarehouseId)
			if err != nil {
				s.Log.Error("error fetch warehouses", zap.Error(err), zap.Int("warehouse_id", payload.WarehouseId))
				errChan <- err
				return
			}

			stockWarehouse = res
		}
	}()

	go func() {
		wg.Wait()
		close(wgDone)
		close(errChan)
	}()

	select {
	case <-wgDone:
		break
	case err := <-errChan:
		return warehouse, models.StockWarehouse{}, err
	}

	return warehouse, stockWarehouse, nil
}

func (s *service) GetWarehouses(ctx context.Context, userClaim dto.UserClaimJwt, payload dto.ParameterQueryWarehouse) (any, error) {
	q := "user_id = ?"
	if payload.Status != "" {
		q += fmt.Sprintf(" AND status = %s", payload.Status)
	}

	fields := "id,name,location,is_active,created_at,updated_at"
	warehouses, err := s.WarehouseRepository.Find(ctx, fields, q, userClaim.UserId)
	if err != nil {
		s.Log.Error("error fetch warehouses", zap.Error(err), zap.Int("user_id", userClaim.UserId))
		return nil, err
	}

	return warehouses, err
}

func (s *service) TransferProductWarehouse(ctx context.Context, userClaim dto.UserClaimJwt, fromId, toId int) error {
	fromWarehouse, toWarehouse, err := s.InitialTransferProductWarehouse(ctx, userClaim, fromId, toId)
	if err != nil {
		return err
	}

	tx := s.StockLevelRepository.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		s.Log.Error("error begin transaction", zap.Error(err))
		return err
	}

	if err := s.ProcessTransferProductWarehouse(tx, fromWarehouse, toWarehouse); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Error("error commit transaction", zap.Error(err))
		return err
	}

	return nil
}

func (s *service) AddWarehouse(ctx context.Context, payload dto.PayloadAddWarehouse, userClaim dto.UserClaimJwt) (dto.ResponseWarehouse, error) {
	warehouseDt, err := s.WarehouseRepository.FindOne(ctx, "id,name", "name = ? and user_id = ?", payload.Name, userClaim.UserId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.Log.Error("error finding warehouse", zap.Error(err))
		return dto.ResponseWarehouse{}, err
	}

	if warehouseDt != (models.Warehouse{}) {
		return dto.ResponseWarehouse{}, constants.WarehouseAlreadyExisted
	}

	warehouse := models.Warehouse{
		Name:      payload.Name,
		Location:  payload.Location,
		UserId:    userClaim.UserId,
		IsActive:  true,
		CreatedAt: time.Now().In(util.LocationTime),
		UpdatedAt: time.Now().In(util.LocationTime),
	}

	if err := s.WarehouseRepository.Create(ctx, &warehouse); err != nil {
		s.Log.Error("error creating warehouse", zap.Error(err))
		return dto.ResponseWarehouse{}, err
	}

	s.Log.Info("success create warehouse", zap.Any("warehouse", warehouse))

	return dto.ResponseWarehouse{
		ID:       warehouse.ID,
		Name:     warehouse.Name,
		Location: warehouse.Location,
		UserId:   warehouse.UserId,
		IsActive: warehouse.IsActive,
	}, nil
}

func (s *service) InitialTransferProductWarehouse(ctx context.Context, userClaim dto.UserClaimJwt, fromId, toId int) (models.Warehouse, models.Warehouse, error) {
	var (
		wg            sync.WaitGroup
		fromWarehouse models.Warehouse
		toWarehouse   models.Warehouse
	)

	errChan := make(chan error)
	wgDone := make(chan struct{})

	c, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	wg.Add(2)

	go func() {
		defer wg.Done()

		select {
		case <-c.Done():
			return
		default:
			res, err := s.WarehouseRepository.FindOne(c, "id,name", "id = ? and user_id = ? and is_active = 1", fromId, userClaim.UserId)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				errChan <- err
				return
			}

			if res == (models.Warehouse{}) {
				errChan <- constants.FromWarehouseNotFound
			}

			fromWarehouse = res
		}
	}()

	go func() {
		defer wg.Done()

		select {
		case <-c.Done():
			return
		default:
			res, err := s.WarehouseRepository.FindOne(c, "id,name", "id = ? and user_id = ? and is_active = 1", toId, userClaim.UserId)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				errChan <- err
				return
			}

			if res == (models.Warehouse{}) {
				errChan <- constants.ToWarehouseNotFound
			}

			toWarehouse = res
		}
	}()

	go func() {
		wg.Wait()
		close(wgDone)
		close(errChan)
	}()

	select {
	case <-wgDone:
		break
	case err := <-errChan:
		return models.Warehouse{}, models.Warehouse{}, err
	}

	return fromWarehouse, toWarehouse, nil
}

func (s *service) ProcessTransferProductWarehouse(tx *gorm.DB, fromWarehouse, toWarehouse models.Warehouse) error {
	stockLevelFrom, err := s.StockLevelRepository.FindTx(tx, "updated_at asc", "warehouse_id = ?", fromWarehouse.ID)
	if err != nil {
		return err
	}

	for _, slF := range stockLevelFrom {
		stockDest, err := s.StockLevelRepository.FindOneTx(tx, "updated_at asc", "warehouse_id = ? and product_id = ?", toWarehouse.ID, slF.ProductId)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if stockDest == (models.StockLevelProduct{}) {
			stockLevel := models.StockLevel{
				ProductId:     slF.ProductId,
				WarehouseId:   toWarehouse.ID,
				Stock:         slF.Stock,
				ReservedStock: slF.ReservedStock,
				CreatedAt:     time.Now().In(util.LocationTime),
				UpdatedAt:     time.Now().In(util.LocationTime),
			}
			if err := s.StockLevelRepository.Create(tx, &stockLevel); err != nil {
				s.Log.Error("error creating stock level", zap.Error(err))
				return err
			}

			if err := s.EmptyWarehouse(tx, slF); err != nil {
				return err
			}

			s.Log.Info("success process stock level", zap.Any("stockLevelId", stockLevel.ID))
		} else {
			updatedStockLevel := models.StockLevel{
				Stock:         slF.Stock + stockDest.Stock,
				ReservedStock: slF.ReservedStock + stockDest.ReservedStock,
				UpdatedAt:     time.Now().In(util.LocationTime),
			}

			if err := s.StockLevelRepository.UpdateOneTx(tx, &updatedStockLevel, "stock,reserved_stock", "warehouse_id = ? and product_id = ?", stockDest.WarehouseId, slF.ProductId); err != nil {
				return err
			}

			if err := s.HandlingReservedStock(tx, slF.ID, stockDest.ID); err != nil {
				return err
			}

			if err := s.EmptyWarehouse(tx, slF); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *service) EmptyWarehouse(tx *gorm.DB, stockLevel models.StockLevelProduct) error {
	stockLevelUpdate := models.StockLevel{Stock: 0, ReservedStock: 0, UpdatedAt: time.Now().In(util.LocationTime)}
	if err := s.StockLevelRepository.UpdateOneTx(tx, &stockLevelUpdate, "stock,reserved_stock", "warehouse_id = ?", stockLevel.WarehouseId); err != nil {
		return err
	}

	s.Log.Info("success empty warehouse", zap.Int("warehouseId", stockLevel.WarehouseId))
	return nil
}

func (s *service) HandlingReservedStock(tx *gorm.DB, stockIdFrom, stockIdDest int) error {
	orderDetail := models.OrderDetail{StockId: stockIdDest, UpdatedAt: time.Now().In(util.LocationTime)}
	if err := s.OrderDetailRepository.UpdateOneTx(tx, &orderDetail, "stock_id,updated_at", "stock_id = ? and expired_at > ?", stockIdFrom, time.Now().Format("2006-01-02 15:04:05")); err != nil {
		return err
	}

	s.Log.Info("success update stock order", zap.Int("stockId", stockIdFrom))
	return nil
}
