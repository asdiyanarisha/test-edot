package order

import (
	"errors"
	"go.uber.org/zap"
	"golang.org/x/net/context"
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
	CreateOrder(ctx context.Context, userClaim dto.UserClaimJwt, payload dto.PayloadCreateOrder) error
}

type service struct {
	Log                    *zap.Logger
	UserRepository         repository.UserRepositoryInterface
	ShopRepository         repository.ShopRepositoryInterface
	ProductRepository      repository.ProductRepositoryInterface
	OrderRepository        repository.OrderRepositoryInterface
	StockLevelRepository   repository.StockLevelRepositoryInterface
	OrderDetailsRepository repository.OrderDetailRepositoryInterface
}

func NewService(f *factory.Factory) Service {
	return &service{
		Log:                    f.Log,
		UserRepository:         f.UserRepository,
		ShopRepository:         f.ShopRepository,
		ProductRepository:      f.ProductRepository,
		OrderRepository:        f.OrderRepository,
		StockLevelRepository:   f.StockLevelRepository,
		OrderDetailsRepository: f.OrderDetailRepository,
	}
}

func (s *service) CreateOrder(ctx context.Context, userClaim dto.UserClaimJwt, payload dto.PayloadCreateOrder) error {
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

	orders, grandTotal, err := s.ProcessOrder(tx, payload)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := s.InsertOrder(tx, userClaim, orders, grandTotal); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Error("error commit transaction", zap.Error(err))
		return err
	}

	return nil
}

func (s *service) ProcessOrder(tx *gorm.DB, payload dto.PayloadCreateOrder) ([]models.OrderDetail, float64, error) {
	var (
		orderDetails []models.OrderDetail
		grandTotal   float64
	)
	mapProductId := make(map[int]bool)

	for _, item := range payload.Items {
		// filter product id duplicated
		if _, ok := mapProductId[item.ProductId]; ok {
			return nil, 0, constants.DuplicateProduct
		}

		// get stock level
		stock, err := s.StockLevelRepository.FindOneTx(tx, "updated_at asc", "product_id = ?", item.ProductId)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.OrderDetail{}, 0, err
		}

		totalPrice := stock.Product.Price * float64(item.Qty)
		// update stock level
		orderDetails = append(orderDetails, models.OrderDetail{
			ProductId: item.ProductId,
			Qty:       item.Qty,
			Total:     totalPrice,
			CreatedAt: time.Now().In(util.LocationTime),
			UpdatedAt: time.Now().In(util.LocationTime),
		})
		updatedData := models.StockLevel{Stock: stock.Stock - item.Qty, ReservedStock: stock.ReservedStock + item.Qty}
		if err := s.StockLevelRepository.UpdateOneTx(tx, &updatedData, "stock,reserved_stock", "id = ?", stock.ID); err != nil {
			return []models.OrderDetail{}, 0, err
		}

		mapProductId[item.ProductId] = true
		grandTotal += totalPrice
	}

	return orderDetails, grandTotal, nil
}

func (s *service) InsertOrder(tx *gorm.DB, userClaim dto.UserClaimJwt, items []models.OrderDetail, grandTotal float64) error {
	now := time.Now().In(util.LocationTime)
	expiredAt := now.Add(time.Minute * 5)

	dataOrder := models.Order{
		Id:        0,
		OrderNo:   "XWADWA1213451",
		UserId:    userClaim.UserId,
		IsPayment: false,
		Total:     grandTotal,
		ExpiredAt: expiredAt,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.OrderRepository.Create(tx, &dataOrder); err != nil {
		return err
	}

	for _, item := range items {
		orderDetail := models.OrderDetail{
			OrderId:   dataOrder.Id,
			ProductId: item.ProductId,
			Qty:       item.Qty,
			Total:     item.Total,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}

		if err := s.OrderDetailsRepository.Create(tx, &orderDetail); err != nil {
			return err
		}
	}

	return nil
}
