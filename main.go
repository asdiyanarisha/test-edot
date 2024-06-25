package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"test-asset-fendr/database"
	"test-asset-fendr/src/factory"
	"test-asset-fendr/src/http"
	"test-asset-fendr/util"
)

func main() {
	database.CreateConn()

	f := factory.NewFactory()

	g := gin.New()
	http.NewHttp(g, f)

	if err := g.Run(":" + util.GetEnv("APP_PORT", "8080")); err != nil {
		log.Fatal("Can't start server.")
	}
}
