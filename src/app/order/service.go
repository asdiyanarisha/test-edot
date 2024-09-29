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
	PaymentOrder(ctx context.Context, userClaim dto.UserClaimJwt, orderId int) error
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

func (s *service) PaymentOrder(ctx context.Context, userClaim dto.UserClaimJwt, orderId int) error {
	now := time.Now().In(util.LocationTime)
	tx := s.OrderRepository.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	query := "id = ? and user_id = ? and is_payment = 0 and expired_at > ?"
	order, err := s.OrderRepository.FindOneTx(tx, "id,is_payment,order_no", query, orderId, userClaim.UserId, now.Format("2006-01-02 15:04:05"))
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return constants.OrderNotFound
		}

		s.Log.Error("error get order", zap.Error(err))
		return err
	}

	if err := s.ProcessPaymentOrder(tx, order); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Error; err != nil {
		s.Log.Error("error begin transaction", zap.Error(err))
		return err
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Error("error commit transaction", zap.Error(err))
		return err
	}

	return nil
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
		stock, err := s.StockLevelRepository.FindOneTx(tx, "updated_at asc", "product_id = ? and stock > 0", item.ProductId)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.OrderDetail{}, 0, err
		}

		if stock == (models.StockLevelProduct{}) {
			return []models.OrderDetail{}, 0, constants.StockProductEmpty
		}

		totalPrice := stock.Product.Price * float64(item.Qty)
		// update stock level
		orderDetails = append(orderDetails, models.OrderDetail{
			ProductId: item.ProductId,
			StockId:   stock.ID,
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
		OrderNo:   util.CreateOrderNo(),
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
			StockId:   item.StockId,
			Qty:       item.Qty,
			Total:     item.Total,
			ExpiredAt: expiredAt,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}

		if err := s.OrderDetailsRepository.Create(tx, &orderDetail); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) ProcessPaymentOrder(tx *gorm.DB, order models.Order) error {
	fields := "id,product_id,stock_id,qty"
	orderDetails, err := s.OrderDetailsRepository.FindTx(tx, fields, "order_id = ? ", order.Id)
	if err != nil {
		s.Log.Error("error get order details", zap.Error(err))
		return err
	}

	for _, detail := range orderDetails {
		stock, err := s.StockLevelRepository.FindOneTx(tx, "updated_at asc", "id = ?", detail.StockId)
		if err != nil {
			return err
		}

		updatedData := models.StockLevel{ReservedStock: stock.ReservedStock - detail.Qty, UpdatedAt: time.Now().In(util.LocationTime)}
		if err := s.StockLevelRepository.UpdateOneTx(tx, &updatedData, "reserved_stock,updated_at", "id = ?", detail.StockId); err != nil {
			s.Log.Error("error update stock", zap.Error(err))
			return err
		}

		s.Log.Info("successfully deduct stock", zap.Int("stock", detail.StockId))
	}

	updateOrder := models.Order{IsPayment: true, UpdatedAt: time.Now().In(util.LocationTime)}
	if err := s.OrderRepository.UpdateOneTx(tx, &updateOrder, "is_payment,updated_at", "id = ?", order.Id); err != nil {
		s.Log.Error("error update order", zap.Error(err))
		return err
	}

	return nil
}
