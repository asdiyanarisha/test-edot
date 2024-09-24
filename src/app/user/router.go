package user

import "github.com/gin-gonic/gin"

func (h *handler) UserRouter(g *gin.RouterGroup) {
	g.POST("register", h.RegisterUser)
}
