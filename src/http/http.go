package http

import (
	"github.com/gin-gonic/gin"
	"test-asset-fendr/src/app/post"
	"test-asset-fendr/src/factory"
	"test-asset-fendr/src/middleware"
)

func NewHttp(g *gin.Engine, f *factory.Factory) {
	g.Use(middleware.CORSMiddleware())
	g.Use(gin.Logger(), gin.Recovery())

	// Here we define a router group
	api := g.Group("/api")

	post.NewHandler(f).PostRouter(api.Group("posts"))
}
