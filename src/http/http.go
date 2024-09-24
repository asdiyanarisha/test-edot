package http

import (
	"github.com/gin-gonic/gin"
	"test-edot/src/app/user"
	"test-edot/src/factory"
	"test-edot/src/middleware"
)

func NewHttp(g *gin.Engine, f *factory.Factory) {
	g.Use(middleware.CORSMiddleware())
	g.Use(gin.Logger(), gin.Recovery())

	// Here we define a router group
	api := g.Group("/api")

	//post.NewHandler(f).PostRouter(api.Group("posts"))

	userGroup := api.Group("users")
	user.NewHandler(f).UserRouter(userGroup)

	// bearer section

	userGroup.Use(middleware.Bearer())
	user.NewHandler(f).UserBearerRouter(userGroup)

}
