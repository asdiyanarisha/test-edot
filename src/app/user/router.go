package user

import "github.com/gin-gonic/gin"

func (h *handler) UserRouter(g *gin.RouterGroup) {
	g.POST("register", h.RegisterUser)
	g.POST("login", h.LoginUser)
}

func (h *handler) UserBearerRouter(g *gin.RouterGroup) {
	g.GET("me", h.UserMe)
}
