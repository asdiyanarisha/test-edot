package product

import (
	"context"
	"errors"
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
	AddProduct(ctx context.Context, payload dto.PayloadAddProduct, userClaim dto.UserClaimJwt) error
}

type service struct {
	Log                  *zap.Logger
	UserRepository       repository.UserRepositoryInterface
	ShopRepository       repository.ShopRepositoryInterface
	ProductRepository    repository.ProductRepositoryInterface
	StockLevelRepository repository.StockLevelRepositoryInterface
	WarehouseRepository  repository.WarehouseRepositoryInterface
}

func NewService(f *factory.Factory) Service {
	return &service{
		Log:                  f.Log,
		UserRepository:       f.UserRepository,
		ShopRepository:       f.ShopRepository,
		ProductRepository:    f.ProductRepository,
		StockLevelRepository: f.StockLevelRepository,
		WarehouseRepository:  f.WarehouseRepository,
	}
}

func (s *service) AddProduct(ctx context.Context, payload dto.PayloadAddProduct, userClaim dto.UserClaimJwt) error {
	if err := s.ValidateAddProduct(ctx, payload, userClaim); err != nil {
		return err
	}

	if err := s.CreateProduct(payload, payload.ShopId); err != nil {
		return err
	}

	return nil
}

func (s *service) ValidateAddProduct(ctx context.Context, payload dto.PayloadAddProduct, userClaim dto.UserClaimJwt) error {
	var (
		wg sync.WaitGroup
	)

	errChan := make(chan error)
	wgDone := make(chan struct{})

	c, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	wg.Add(3)

	go func() {
		defer wg.Done()

		select {
		case <-c.Done():
			return
		default:
			res, err := s.ShopRepository.FindOne(ctx, "id", "user_id = ? and id = ?", userClaim.UserId, payload.ShopId)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				s.Log.Error("error get shop", zap.Error(err), zap.Any("payload", payload))
				errChan <- err
				return
			}

			if res == (models.Shop{}) {
				errChan <- constants.ShopNotFound
				return
			}

			return
		}
	}()

	go func() {
		defer wg.Done()

		select {
		case <-c.Done():
			return

		default:
			productData, err := s.ProductRepository.FindOne(ctx, "id", "sku = ? and shop_id = ?", payload.Sku, payload.ShopId)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				s.Log.Error("error get product", zap.Error(err), zap.Any("payload", payload))
				errChan <- err
				return
			}

			if productData != (models.Product{}) {
				errChan <- constants.ProductAlreadyInserted
				return
			}
		}
	}()

	go func() {
		defer wg.Done()

		select {
		case <-c.Done():
			return
		default:
			warehouseData, err := s.WarehouseRepository.FindOne(ctx, "id", "user_id = ? and id = ?", userClaim.UserId, payload.WarehouseId)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				s.Log.Error("error get warehouse", zap.Error(err), zap.Any("payload", payload))
				errChan <- err
				return
			}

			if warehouseData == (models.Warehouse{}) {
				errChan <- constants.WarehouseNotFound
				return
			}
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
		return err
	}

	return nil
}

func (s *service) CreateProduct(payload dto.PayloadAddProduct, shopId int) error {
	tx := s.ProductRepository.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		s.Log.Error("error begin transaction", zap.Error(err))
		return err
	}

	product := models.Product{
		Name:      payload.Name,
		Sku:       payload.Sku,
		Price:     payload.Price,
		ShopId:    shopId,
		CreatedAt: time.Now().In(util.LocationTime),
		UpdatedAt: time.Now().In(util.LocationTime),
	}
	if err := s.ProductRepository.Create(tx, &product); err != nil {
		tx.Rollback()
		s.Log.Error("error insert product", zap.String("product", product.Name), zap.Error(err))
		return err
	}

	stockLevel := models.StockLevel{
		ProductId:     product.Id,
		WarehouseId:   payload.WarehouseId,
		Stock:         payload.Qty,
		ReservedStock: 0,
		CreatedAt:     time.Now().In(util.LocationTime),
		UpdatedAt:     time.Now().In(util.LocationTime),
	}

	if err := s.StockLevelRepository.Create(tx, &stockLevel); err != nil {
		tx.Rollback()
		s.Log.Error("error insert stock level", zap.String("product", product.Name), zap.Error(err))
		return err
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Error("error commit transaction", zap.Error(err))
		return err
	}

	s.Log.Info("success insert product", zap.String("product", product.Name))

	return nil
}
