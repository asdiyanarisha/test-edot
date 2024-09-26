package warehouse

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
}

type service struct {
	Log                 *zap.Logger
	UserRepository      repository.UserRepositoryInterface
	ShopRepository      repository.ShopRepositoryInterface
	ProductRepository   repository.ProductRepositoryInterface
	WarehouseRepository repository.WarehouseRepositoryInterface
}

func NewService(f *factory.Factory) Service {
	return &service{
		Log:                 f.Log,
		UserRepository:      f.UserRepository,
		ShopRepository:      f.ShopRepository,
		ProductRepository:   f.ProductRepository,
		WarehouseRepository: f.WarehouseRepository,
	}
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
