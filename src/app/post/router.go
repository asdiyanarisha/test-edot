package post

import (
	"github.com/gin-gonic/gin"
)

func (h *handler) PostRouter(g *gin.RouterGroup) {
	g.POST("", h.AddPostHandler)
}
