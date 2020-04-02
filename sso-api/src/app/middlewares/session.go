package middlewares

import (
	"login/src/app/global"
	"login/src/app/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetSession() gin.HandlerFunc {
	return sessions.Sessions(global.Conf.Session.Name, models.RedisSessionStore)
}
