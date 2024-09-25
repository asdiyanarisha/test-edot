package shop

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
	CreateShop(ctx context.Context, user dto.UserClaimJwt, payload dto.PayloadCreateShop) error
}

type service struct {
	Log            *zap.Logger
	UserRepository repository.UserRepositoryInterface
	ShopRepository repository.ShopRepositoryInterface
}

func NewService(f *factory.Factory) Service {
	return &service{
		Log:            f.Log,
		UserRepository: f.UserRepository,
		ShopRepository: f.ShopRepository,
	}
}

func (s *service) CreateShop(ctx context.Context, user dto.UserClaimJwt, payload dto.PayloadCreateShop) error {
	shop, err := s.ShopRepository.FindOne(ctx, "name", "name = ? and user_id = ?", payload.Name, user.UserId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if shop != (models.Shop{}) {
		return constants.ShopAlreadyInserted
	}

	shopData := models.Shop{
		Name:      payload.Name,
		Location:  payload.Location,
		UserId:    user.UserId,
		CreatedAt: time.Now().In(util.LocationTime),
		UpdatedAt: time.Now().In(util.LocationTime),
	}

	if err := s.ShopRepository.Create(ctx, &shopData); err != nil {
		return err
	}

	s.Log.Info("shop created", zap.String("name", shopData.Name), zap.String("location", shopData.Location))

	return nil
}
