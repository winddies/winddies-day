package middlewares

import (
	"luck-home/winddies/sso-api/app/global"
	"luck-home/winddies/sso-api/app/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetSession() gin.HandlerFunc {
	return sessions.Sessions(global.Conf.Session.Name, models.RedisSessionStore)
}
