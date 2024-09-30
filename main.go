package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"test-edot/database"
	"test-edot/src/factory"
	"test-edot/src/http"
	"test-edot/src/scheduler"
	"test-edot/util"
)

func main() {
	var (
		m string
	)
	database.CreateConn()

	flag.StringVar(&m, "m", "", `This flag is used for mode [scheduler]`)
	flag.Parse()

	f := factory.NewFactory()
	if m == "scheduler" {
		scheduler.RunScheduler(f)
		return
	}

	g := gin.New()
	http.NewHttp(g, f)

	if err := g.Run(":" + util.GetEnv("APP_PORT", "8080")); err != nil {
		log.Fatal("Can't start server.")
	}
}
