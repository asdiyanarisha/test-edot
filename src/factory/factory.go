package factory

import (
	"go.uber.org/zap"
	"test-edot/database"
	"test-edot/src/repository"
)

type Factory struct {
	Log                  *zap.Logger
	PostRepository       repository.PostRepositoryInterface
	TagRepository        repository.TagRepositoryInterface
	PostTagRepository    repository.PostTagRepositoryInterface
	UserRepository       repository.UserRepositoryInterface
	ShopRepository       repository.ShopRepositoryInterface
	ProductRepository    repository.ProductRepositoryInterface
	WarehouseRepository  repository.WarehouseRepositoryInterface
	StockLevelRepository repository.StockLevelRepositoryInterface
}

func NewFactory() *Factory {
	db := database.GetConn()
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	return &Factory{
		Log:                  logger,
		PostRepository:       repository.NewPostRepository(db),
		TagRepository:        repository.NewTagRepository(db),
		PostTagRepository:    repository.NewPostTagRepository(db),
		UserRepository:       repository.NewUserRepository(db),
		ShopRepository:       repository.NewShopRepository(db),
		ProductRepository:    repository.NewProductRepository(db),
		WarehouseRepository:  repository.NewWarehouseRepository(db),
		StockLevelRepository: repository.NewStockLevelRepository(db),
	}
}
