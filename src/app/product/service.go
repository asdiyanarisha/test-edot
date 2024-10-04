package product

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
	AddProduct(ctx context.Context, payload dto.PayloadAddProduct, userClaim dto.UserClaimJwt) error
	ProductList(ctx context.Context, payload dto.ParameterQuery) (any, error)
	TransferProductWarehouse(ctx context.Context, payload dto.TransferProductWarehouse, userClaim dto.UserClaimJwt, productId int) error
	GetProductDetail(ctx context.Context, userClaim dto.UserClaimJwt, productId int) (any, error)
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

func (s *service) TransferProductWarehouse(ctx context.Context, payload dto.TransferProductWarehouse, userClaim dto.UserClaimJwt, productId int) error {
	initialData, err := s.SetupProcessTransferProduct(ctx, userClaim, payload, productId)
	if err != nil {
		return err
	}

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

	if err := s.ProcessTransferProduct(tx, initialData, payload); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Error("error commit transaction", zap.Error(err))
		return err
	}

	return nil
}

func (s *service) GetProductDetail(ctx context.Context, userClaim dto.UserClaimJwt, productId int) (any, error) {
	var (
		stock         int
		reservedStock int
	)

	if userClaim.Role != constants.ROLE_ADMIN_SHOP {
		return nil, constants.RoleUserInvalid
	}

	selectField := "id,name,sku,price,shop_id"
	product, err := s.ProductRepository.GetProductDetail(ctx, selectField, "id = ?", productId)
	if err != nil {
		return nil, err
	}
	for _, level := range product.Stock {
		stock += level.Stock
		reservedStock += level.ReservedStock
	}

	productRes := dto.ProductDetailResponse{
		Id:            product.Id,
		Name:          product.Name,
		Price:         product.Price,
		Sku:           product.Sku,
		Shop:          product.Shop.Name,
		Stock:         stock,
		ReservedStock: reservedStock,
	}

	return productRes, nil
}

func (s *service) ProductList(ctx context.Context, payload dto.ParameterQuery) (any, error) {
	var query string
	limit := 20
	if payload.Limit > 0 {
		limit = payload.Limit
	}

	if payload.Search != "" {
		query += fmt.Sprintf("name LIKE '%s'", "%"+payload.Search+"%")
	}

	selectField := "id,name,sku,price,shop_id"
	products, err := s.ProductRepository.GetProductDetails(ctx, payload.Offset, limit, selectField, query)
	if err != nil {
		return nil, err
	}

	var resProducts []dto.ProductResponse
	for _, product := range products {
		var stock int
		for _, level := range product.Stock {
			stock += level.Stock
		}

		resProducts = append(resProducts, dto.ProductResponse{
			Id:    product.Id,
			Name:  product.Name,
			Price: product.Price,
			Sku:   product.Sku,
			Shop:  product.Shop.Name,
			Stock: stock,
		})
	}

	return resProducts, nil
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

func (s *service) SetupProcessTransferProduct(ctx context.Context, userClaim dto.UserClaimJwt, payload dto.TransferProductWarehouse, productId int) (dto.InitialTransferProduct, error) {
	var (
		wg            sync.WaitGroup
		product       models.Product
		fromWarehouse models.Warehouse
		toWarehouse   models.Warehouse
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
			res, err := s.ProductRepository.FindOne(c, "id,name", "id = ?", productId)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				errChan <- err
				return
			}

			if res == (models.Product{}) {
				errChan <- constants.ProductNotFound
			}

			product = res
		}
	}()

	go func() {
		defer wg.Done()

		select {
		case <-c.Done():
			return
		default:
			res, err := s.WarehouseRepository.FindOne(c, "id,name", "id = ? and user_id = ? and is_active = 1", payload.FromWarehouseId, userClaim.UserId)
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
			res, err := s.WarehouseRepository.FindOne(c, "id,name", "id = ? and user_id = ? and is_active = 1", payload.ToWarehouseId, userClaim.UserId)
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
		return dto.InitialTransferProduct{}, err
	}

	return dto.InitialTransferProduct{
		Product:       product,
		FromWarehouse: fromWarehouse,
		ToWarehouse:   toWarehouse,
	}, nil
}

func (s *service) ProcessTransferProduct(tx *gorm.DB, initialData dto.InitialTransferProduct, payload dto.TransferProductWarehouse) error {
	q := "product_id = ? and warehouse_id = ?"
	stockFrom, err := s.StockLevelRepository.FindOneTx(tx, "updated_at asc", q, initialData.Product.Id, initialData.FromWarehouse.ID)
	if err != nil {
		s.Log.Error("error get stock", zap.String("product", initialData.Product.Name), zap.Error(err))
		return err
	}

	if stockFrom.Stock < payload.Qty {
		return constants.NotEnoughStockToTransfer
	}

	stockDest, err := s.StockLevelRepository.FindOneTx(tx, "updated_at asc", q, initialData.Product.Id, initialData.ToWarehouse.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.Log.Error("error get stock", zap.String("product", initialData.Product.Name), zap.Error(err))
		return err
	}

	if stockDest == (models.StockLevelProduct{}) {
		stockLevel := models.StockLevel{
			ProductId:   initialData.Product.Id,
			WarehouseId: initialData.ToWarehouse.ID,
			Stock:       payload.Qty,
			CreatedAt:   time.Now().In(util.LocationTime),
			UpdatedAt:   time.Now().In(util.LocationTime),
		}
		if err := s.StockLevelRepository.Create(tx, &stockLevel); err != nil {
			s.Log.Error("error creating stock level", zap.Error(err))
			return err
		}
	} else {
		updatedStockLevel := models.StockLevel{
			Stock:     stockDest.Stock + payload.Qty,
			UpdatedAt: time.Now().In(util.LocationTime),
		}

		if err := s.StockLevelRepository.UpdateOneTx(tx, &updatedStockLevel, "stock,updated_at", "warehouse_id = ? and product_id = ?", stockDest.WarehouseId, initialData.Product.Id); err != nil {
			return err
		}
	}

	updatedStockLevel := models.StockLevel{
		Stock:     stockFrom.Stock - payload.Qty,
		UpdatedAt: time.Now().In(util.LocationTime),
	}

	if err := s.StockLevelRepository.UpdateOneTx(tx, &updatedStockLevel, "stock,updated_at", "warehouse_id = ? and product_id = ?", stockFrom.WarehouseId, initialData.Product.Id); err != nil {
		return err
	}

	s.Log.Info("success move product stock to another warehouse")

	return nil

}
