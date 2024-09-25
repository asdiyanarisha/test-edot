package warehouse

import (
	"context"
	"errors"
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
