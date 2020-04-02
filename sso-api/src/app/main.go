package main

import (
	"flag"
	"login/src/app/global"
	"login/src/app/middlewares"
	"login/src/app/models"
	"login/src/app/routes"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "get config file")
}

func main() {
	flag.Parse()
	global.Init(configPath)
	models.Init()

	gin.SetMode(getGinMode())
	app := gin.New()
	app.Use(middlewares.Logger(), gin.Recovery())
	routes.InitRoutes(app)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        app,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}

func getGinMode() string {
	switch global.Conf.Mode {
	case global.DevMode:
		return gin.DebugMode
	case global.TestMode:
		return gin.TestMode
	case global.ProdMode:
		return gin.ReleaseMode
	default:
		return gin.DebugMode
	}
}
