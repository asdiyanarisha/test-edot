package shop

import "github.com/gin-gonic/gin"

func (h *handler) ShopRouter(g *gin.RouterGroup) {
	g.POST("", h.CreateShop)
}
