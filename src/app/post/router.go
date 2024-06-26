package post

import (
	"github.com/gin-gonic/gin"
)

func (h *handler) PostRouter(g *gin.RouterGroup) {
	g.POST("", h.AddPostHandler)
	g.GET("", h.GetPostHandler)
	g.GET("/:id", h.GetPostByIdHandler)
	g.PUT("/:id", h.UpdatePostHandler)
	g.DELETE("/:id", h.DeletePostHandler)
}
