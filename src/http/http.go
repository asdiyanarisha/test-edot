package http

import (
	"github.com/gin-gonic/gin"
	"test-edot/src/app/product"
	"test-edot/src/app/shop"
	"test-edot/src/app/user"
	"test-edot/src/factory"
	"test-edot/src/middleware"
)

func NewHttp(g *gin.Engine, f *factory.Factory) {
	g.Use(middleware.CORSMiddleware())
	g.Use(gin.Logger(), gin.Recovery())

	// Here we define a router group
	api := g.Group("/api")

	product.NewHandler(f).ProductRouter(api.Group("/products"))

	// user section
	userGroup := api.Group("users")
	user.NewHandler(f).UserRouter(userGroup)

	userGroup.Use(middleware.Bearer())
	user.NewHandler(f).UserBearerRouter(userGroup)

	// shop section
	shopsGroup := api.Group("shops")
	shopsGroup.Use(middleware.BearerShop())

	shop.NewHandler(f).ShopRouter(shopsGroup)

	// product section
	productGroup := api.Group("product")
	product.NewHandler(f).ProductBearerShopRouter(productGroup)
}
