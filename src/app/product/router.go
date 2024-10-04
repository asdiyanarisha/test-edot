package product

import (
	"github.com/gin-gonic/gin"
	"test-edot/src/middleware"
)

func (h *handler) ProductRouter(g *gin.RouterGroup) {
	g.GET("", h.ProductList)
}

func (h *handler) ProductBearerShopRouter(g *gin.RouterGroup) {
	g.Use(middleware.BearerShop())
	g.POST("", h.AddProduct)
	g.POST("/:product_id/transfer", h.TransferProduct)
	g.GET("/:product_id/detail", h.DetailProduct)
}
