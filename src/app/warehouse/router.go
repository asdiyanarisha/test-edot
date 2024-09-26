package warehouse

import (
	"github.com/gin-gonic/gin"
	"test-edot/src/middleware"
)

func (h *handler) WarehouseBearerShopRouter(g *gin.RouterGroup) {
	g.Use(middleware.BearerShop())
	g.POST("", h.AddWarehouse)
	g.GET("", h.GetWarehouses)
}
