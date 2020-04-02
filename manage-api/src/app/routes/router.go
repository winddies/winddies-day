package routes

import (
	"github.com/gin-gonic/gin"
)

func InitRoutes(app *gin.Engine) {
	routeAuth := app.Group("/auth")
	routeAuth.POST("/article/new")
}
