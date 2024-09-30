package order

import (
	"errors"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"strconv"
	"test-edot/constants"
	"test-edot/src/dto"
	"test-edot/src/factory"
	"test-edot/src/models"
	"test-edot/src/repository"
	"test-edot/util"
	"time"
)

type Service interface {
	CreateOrder(ctx context.Context, userClaim dto.UserClaimJwt, payload dto.PayloadCreateOrder) (models.Order, error)
	PaymentOrder(ctx context.Context, userClaim dto.UserClaimJwt, orderId int) error
	ReleaseStockOrder()
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

func (s *service) CreateOrder(ctx context.Context, userClaim dto.UserClaimJwt, payload dto.PayloadCreateOrder) (models.Order, error) {
	tx := s.ProductRepository.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		s.Log.Error("error begin transaction", zap.Error(err))
		return models.Order{}, err
	}

	orders, grandTotal, err := s.ProcessOrder(tx, payload)
	if err != nil {
		tx.Rollback()
		return models.Order{}, err
	}

	order, err := s.InsertOrder(tx, userClaim, orders, grandTotal)
	if err != nil {
		tx.Rollback()
		return models.Order{}, err
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Error("error commit transaction", zap.Error(err))
		return models.Order{}, err
	}

	return order, nil
}

func (s *service) ReleaseStockOrder() {
	ctx := context.Background()

	s.Log.Info("running release stock order")

	now := time.Now().In(util.LocationTime)
	orders, err := s.OrderRepository.FindAll(ctx, "id,expired_at", "expired_at < ? and is_release = 0 and is_payment = 0", now.Format("2006-01-02 15:04:05"))
	if err != nil {
		s.Log.Error("error get order", zap.Error(err))
		return
	}
	tx := s.OrderDetailsRepository.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		s.Log.Error("error begin transaction", zap.Error(err))
		return
	}

	for _, order := range orders {
		detailOrders, err := s.OrderDetailsRepository.FindTx(tx, "id,stock_id,product_id,qty", "order_id = ? and expired_at < ?", order.Id, now.Format("2006-01-02 15:04:05"))
		if err != nil {
			tx.Rollback()
			s.Log.Error("error get order", zap.Error(err))
			return
		}

		for _, detailOrder := range detailOrders {
			stock, err := s.StockLevelRepository.FindOneTx(tx, "updated_at asc", "id = ? and product_id = ?", detailOrder.StockId, detailOrder.ProductId)
			if err != nil {
				tx.Rollback()
				s.Log.Error("error get order", zap.Error(err))
				return
			}

			stockLevel := models.StockLevel{ReservedStock: stock.ReservedStock - detailOrder.Qty, Stock: stock.Stock + detailOrder.Qty}
			if err := s.StockLevelRepository.UpdateOneTx(tx, &stockLevel, "reserved_stock,stock", "id = ?", stock.ID); err != nil {
				tx.Rollback()
				s.Log.Error("error update stock", zap.Error(err))
				return
			}
		}

		updatedOrder := models.Order{IsRelease: true, UpdatedAt: time.Now().In(util.LocationTime)}
		if err := s.OrderRepository.UpdateOneTx(tx, &updatedOrder, "is_release,updated_at", "id = ?", order.Id); err != nil {
			tx.Rollback()
			s.Log.Error("error update order", zap.Error(err))
			return
		}

		s.Log.Info("order stock has released", zap.Int("orderId", order.Id))
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Error("error commit transaction", zap.Error(err))
		return
	}
}

func (s *service) ProcessOrder(tx *gorm.DB, payload dto.PayloadCreateOrder) ([]models.OrderDetail, float64, error) {
	var (
		orderDetails []models.OrderDetail
		grandTotal   float64
	)
	mapProductId := make(map[int]bool)

	for _, item := range payload.Items {
		qty := item.Qty

		// filter product id duplicated
		if _, ok := mapProductId[item.ProductId]; ok {
			return nil, 0, constants.DuplicateProduct
		}

		// get stock level
		stocks, err := s.StockLevelRepository.FindTx(tx, "updated_at asc", "product_id = ? and stock > 0", item.ProductId)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.OrderDetail{}, 0, err
		}

		if len(stocks) == 0 {
			return []models.OrderDetail{}, 0, constants.StockProductEmpty
		}

		for _, stock := range stocks {
			if qty == 0 {
				break
			}

			qtyStock := qty
			if qty > stock.Stock {
				qty = qty - stock.Stock
				qtyStock = stock.Stock
			} else {
				qty = 0
			}

			totalPrice := stock.Product.Price * float64(qtyStock)
			// update stock level
			orderDetails = append(orderDetails, models.OrderDetail{
				ProductId: item.ProductId,
				StockId:   stock.ID,
				Qty:       qtyStock,
				Total:     totalPrice,
				CreatedAt: time.Now().In(util.LocationTime),
				UpdatedAt: time.Now().In(util.LocationTime),
			})
			updatedData := models.StockLevel{Stock: stock.Stock - qtyStock, ReservedStock: stock.ReservedStock + qtyStock}
			if err := s.StockLevelRepository.UpdateOneTx(tx, &updatedData, "stock,reserved_stock", "id = ?", stock.ID); err != nil {
				return []models.OrderDetail{}, 0, err
			}

			grandTotal += totalPrice
		}

		if qty > 0 {
			return nil, 0, constants.NotEnoughStockProduct
		}

		mapProductId[item.ProductId] = true
	}

	return orderDetails, grandTotal, nil
}

func (s *service) InsertOrder(tx *gorm.DB, userClaim dto.UserClaimJwt, items []models.OrderDetail, grandTotal float64) (models.Order, error) {
	expireOrderMinutes, err := strconv.Atoi(util.GetEnv("ORDER_EXPIRE_MINUTE", ""))
	if err != nil {
		return models.Order{}, err
	}

	now := time.Now().In(util.LocationTime)
	expiredAt := now.Add(time.Minute * time.Duration(expireOrderMinutes))

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
		return models.Order{}, err
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
			return models.Order{}, err
		}
	}

	return dataOrder, nil
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
