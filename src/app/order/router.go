package order

import (
	"github.com/gin-gonic/gin"
	"test-edot/src/middleware"
)

func (h *handler) OrderBearerRouter(g *gin.RouterGroup) {
	g.Use(middleware.BearerUser())
	g.POST("", h.CreateOrder)
}
