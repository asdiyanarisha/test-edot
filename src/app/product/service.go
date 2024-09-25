package product

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
	AddProduct(ctx context.Context, payload dto.PayloadAddProduct, userClaim dto.UserClaimJwt) error
}

type service struct {
	Log               *zap.Logger
	UserRepository    repository.UserRepositoryInterface
	ShopRepository    repository.ShopRepositoryInterface
	ProductRepository repository.ProductRepositoryInterface
}

func NewService(f *factory.Factory) Service {
	return &service{
		Log:               f.Log,
		UserRepository:    f.UserRepository,
		ShopRepository:    f.ShopRepository,
		ProductRepository: f.ProductRepository,
	}
}

func (s *service) AddProduct(ctx context.Context, payload dto.PayloadAddProduct, userClaim dto.UserClaimJwt) error {
	shop, err := s.ShopRepository.FindOne(ctx, "id", "user_id = ? and id = ?", userClaim.UserId, payload.ShopId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return constants.ShopNotFound
		}

		s.Log.Error("error get shop", zap.Error(err), zap.Any("payload", payload))
		return err
	}

	productData, err := s.ProductRepository.FindOne(ctx, "id", "sku = ? and shop_id = ?", payload.Sku, shop.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.Log.Error("error get product", zap.Error(err), zap.Any("payload", payload))
		return err
	}

	if productData != (models.Product{}) {
		return constants.ProductAlreadyInserted
	}

	product := models.Product{
		Name:      payload.Name,
		Sku:       payload.Sku,
		Price:     payload.Price,
		ShopId:    shop.ID,
		CreatedAt: time.Now().In(util.LocationTime),
		UpdatedAt: time.Now().In(util.LocationTime),
	}

	if err := s.ProductRepository.Create(ctx, &product); err != nil {

		s.Log.Error("error insert product", zap.String("product", product.Name), zap.Error(err))
		return err
	}

	s.Log.Info("success insert product", zap.String("product", product.Name))
	return nil
}
