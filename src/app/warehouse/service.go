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
	AddWarehouse(ctx context.Context, payload dto.PayloadAddWarehouse, userClaim dto.UserClaimJwt) error
	GetWarehouses(ctx context.Context, userClaim dto.UserClaimJwt, payload dto.ParameterQueryWarehouse) (any, error)
	ChangeStatusWarehouse(ctx context.Context, userClaim dto.UserClaimJwt, payload dto.ParameterChangeStatusWarehouse) error
}

type service struct {
	Log                  *zap.Logger
	UserRepository       repository.UserRepositoryInterface
	ShopRepository       repository.ShopRepositoryInterface
	ProductRepository    repository.ProductRepositoryInterface
	WarehouseRepository  repository.WarehouseRepositoryInterface
	StockLevelRepository repository.StockLevelRepositoryInterface
}

func NewService(f *factory.Factory) Service {
	return &service{
		Log:                  f.Log,
		UserRepository:       f.UserRepository,
		ShopRepository:       f.ShopRepository,
		ProductRepository:    f.ProductRepository,
		WarehouseRepository:  f.WarehouseRepository,
		StockLevelRepository: f.StockLevelRepository,
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

func (s *service) AddWarehouse(ctx context.Context, payload dto.PayloadAddWarehouse, userClaim dto.UserClaimJwt) error {
	warehouseDt, err := s.WarehouseRepository.FindOne(ctx, "id,name", "name = ? and user_id = ?", payload.Name, userClaim.UserId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.Log.Error("error finding warehouse", zap.Error(err))
		return err
	}

	if warehouseDt != (models.Warehouse{}) {
		return constants.WarehouseAlreadyExisted
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
		return err
	}

	s.Log.Info("success create warehouse", zap.Any("warehouse", warehouse))

	return nil
}
